package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"go-backend/internal/config"
	"go-backend/internal/models"
	"go-backend/pkg/logger"

	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationEmail   NotificationType = "email"
	NotificationSMS     NotificationType = "sms"
	NotificationInApp   NotificationType = "in_app"
	NotificationPush    NotificationType = "push"
	NotificationSlack   NotificationType = "slack"
	NotificationDiscord NotificationType = "discord"
)

// NotificationPriority represents the priority of a notification
type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "low"
	PriorityNormal   NotificationPriority = "normal"
	PriorityHigh     NotificationPriority = "high"
	PriorityCritical NotificationPriority = "critical"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusFailed    NotificationStatus = "failed"
	StatusRead      NotificationStatus = "read"
)

// NotificationTemplate represents a notification template
type NotificationTemplate struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	Name      string           `json:"name" gorm:"unique;not null"`
	Type      NotificationType `json:"type" gorm:"not null"`
	Subject   string           `json:"subject"`
	Body      string           `json:"body" gorm:"type:text"`
	Variables string           `json:"variables" gorm:"type:text"` // JSON array of variable names
	IsActive  bool             `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// Notification represents a notification record
type Notification struct {
	ID          uint                 `json:"id" gorm:"primaryKey"`
	UserID      *uint                `json:"user_id" gorm:"index"`
	Type        NotificationType     `json:"type" gorm:"not null"`
	Priority    NotificationPriority `json:"priority" gorm:"default:normal"`
	Status      NotificationStatus   `json:"status" gorm:"default:pending"`
	Subject     string               `json:"subject"`
	Body        string               `json:"body" gorm:"type:text"`
	Recipient   string               `json:"recipient" gorm:"not null"`
	Metadata    string               `json:"metadata" gorm:"type:text"` // JSON data
	ScheduledAt *time.Time           `json:"scheduled_at"`
	SentAt      *time.Time           `json:"sent_at"`
	DeliveredAt *time.Time           `json:"delivered_at"`
	ReadAt      *time.Time           `json:"read_at"`
	FailedAt    *time.Time           `json:"failed_at"`
	Error       string               `json:"error"`
	Retries     int                  `json:"retries" gorm:"default:0"`
	MaxRetries  int                  `json:"max_retries" gorm:"default:3"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`

	// Relationships
	User *models.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// NotificationChannel interface for different notification channels
type NotificationChannel interface {
	Send(notification *Notification) error
	GetType() NotificationType
}

// EmailChannel implements email notifications
type EmailChannel struct {
	config *config.EmailConfig
	logger *logger.Logger
}

// SMSChannel implements SMS notifications
type SMSChannel struct {
	config *config.Config
	logger *logger.Logger
}

// InAppChannel implements in-app notifications
type InAppChannel struct {
	db     *gorm.DB
	logger *logger.Logger
}

// SlackChannel implements Slack notifications
type SlackChannel struct {
	webhookURL string
	logger     *logger.Logger
}

// NotificationService handles all notification operations
type NotificationService struct {
	db       *gorm.DB
	config   *config.Config
	logger   *logger.Logger
	channels map[NotificationType]NotificationChannel
	audit    *AuditService
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *gorm.DB, config *config.Config, logger *logger.Logger, audit *AuditService) *NotificationService {
	service := &NotificationService{
		db:       db,
		config:   config,
		logger:   logger,
		channels: make(map[NotificationType]NotificationChannel),
		audit:    audit,
	}

	// Initialize channels
	service.channels[NotificationEmail] = &EmailChannel{
		config: &config.Email,
		logger: logger,
	}
	service.channels[NotificationSMS] = &SMSChannel{
		config: config,
		logger: logger,
	}
	service.channels[NotificationInApp] = &InAppChannel{
		db:     db,
		logger: logger,
	}

	return service
}

// SendNotification sends a notification immediately
func (ns *NotificationService) SendNotification(notification *Notification) error {
	// Save notification to database
	if err := ns.db.Create(notification).Error; err != nil {
		ns.logger.Error("Failed to save notification", map[string]interface{}{
			"error": err.Error(),
			"type":  notification.Type,
		})
		return err
	}

	// Get the appropriate channel
	channel, exists := ns.channels[notification.Type]
	if !exists {
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}

	// Send the notification
	err := channel.Send(notification)
	if err != nil {
		notification.Status = StatusFailed
		notification.Error = err.Error()
		notification.FailedAt = &[]time.Time{time.Now()}[0]
		ns.db.Save(notification)

		ns.logger.Error("Failed to send notification", map[string]interface{}{
			"error":           err.Error(),
			"notification_id": notification.ID,
			"type":            notification.Type,
		})
		return err
	}

	// Update notification status
	notification.Status = StatusSent
	notification.SentAt = &[]time.Time{time.Now()}[0]
	ns.db.Save(notification)

	// Log audit event
	if ns.audit != nil {
		ns.audit.LogActivity(notification.UserID, "notification_sent", map[string]interface{}{
			"notification_id": notification.ID,
			"type":            notification.Type,
			"recipient":       notification.Recipient,
		})
	}

	return nil
}

// ScheduleNotification schedules a notification for later delivery
func (ns *NotificationService) ScheduleNotification(notification *Notification, scheduledAt time.Time) error {
	notification.ScheduledAt = &scheduledAt
	notification.Status = StatusPending

	if err := ns.db.Create(notification).Error; err != nil {
		ns.logger.Error("Failed to schedule notification", map[string]interface{}{
			"error": err.Error(),
			"type":  notification.Type,
		})
		return err
	}

	ns.logger.Info("Notification scheduled", map[string]interface{}{
		"notification_id": notification.ID,
		"scheduled_at":    scheduledAt,
		"type":            notification.Type,
	})

	return nil
}

// SendFromTemplate sends a notification using a template
func (ns *NotificationService) SendFromTemplate(templateName string, recipient string, userID *uint, variables map[string]interface{}) error {
	// Get template
	var template NotificationTemplate
	if err := ns.db.Where("name = ? AND is_active = ?", templateName, true).First(&template).Error; err != nil {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Parse template
	subject, body, err := ns.parseTemplate(template, variables)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create notification
	notification := &Notification{
		UserID:    userID,
		Type:      template.Type,
		Priority:  PriorityNormal,
		Subject:   subject,
		Body:      body,
		Recipient: recipient,
	}

	return ns.SendNotification(notification)
}

// ProcessScheduledNotifications processes all pending scheduled notifications
func (ns *NotificationService) ProcessScheduledNotifications() error {
	var notifications []Notification
	now := time.Now()

	// Get notifications that are scheduled and ready to send
	if err := ns.db.Where("status = ? AND scheduled_at <= ?", StatusPending, now).Find(&notifications).Error; err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := ns.SendNotification(&notification); err != nil {
			ns.logger.Error("Failed to send scheduled notification", map[string]interface{}{
				"notification_id": notification.ID,
				"error":           err.Error(),
			})
		}
	}

	return nil
}

// RetryFailedNotifications retries failed notifications
func (ns *NotificationService) RetryFailedNotifications() error {
	var notifications []Notification

	// Get failed notifications that haven't exceeded max retries
	if err := ns.db.Where("status = ? AND retries < max_retries", StatusFailed).Find(&notifications).Error; err != nil {
		return err
	}

	for _, notification := range notifications {
		notification.Retries++

		// Get the appropriate channel
		channel, exists := ns.channels[notification.Type]
		if !exists {
			continue
		}

		// Retry sending
		err := channel.Send(&notification)
		if err != nil {
			notification.Error = err.Error()
			if notification.Retries >= notification.MaxRetries {
				ns.logger.Error("Notification failed after max retries", map[string]interface{}{
					"notification_id": notification.ID,
					"retries":         notification.Retries,
				})
			}
		} else {
			notification.Status = StatusSent
			notification.SentAt = &[]time.Time{time.Now()}[0]
			notification.Error = ""
		}

		ns.db.Save(&notification)
	}

	return nil
}

// GetUserNotifications gets notifications for a user
func (ns *NotificationService) GetUserNotifications(userID uint, limit, offset int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	// Count total notifications
	ns.db.Model(&Notification{}).Where("user_id = ?", userID).Count(&total)

	// Get notifications with pagination
	err := ns.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, total, err
}

// MarkAsRead marks a notification as read
func (ns *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	now := time.Now()
	result := ns.db.Model(&Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"status":  StatusRead,
			"read_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("notification not found or access denied")
	}

	return nil
}

// CreateTemplate creates a new notification template
func (ns *NotificationService) CreateTemplate(template *NotificationTemplate) error {
	return ns.db.Create(template).Error
}

// parseTemplate parses a template with variables
func (ns *NotificationService) parseTemplate(template NotificationTemplate, variables map[string]interface{}) (string, string, error) {
	// Parse subject
	subjectTemplate, err := template2.New("subject").Parse(template.Subject)
	if err != nil {
		return "", "", err
	}

	var subjectBuf bytes.Buffer
	if err := subjectTemplate.Execute(&subjectBuf, variables); err != nil {
		return "", "", err
	}

	// Parse body
	bodyTemplate, err := template2.New("body").Parse(template.Body)
	if err != nil {
		return "", "", err
	}

	var bodyBuf bytes.Buffer
	if err := bodyTemplate.Execute(&bodyBuf, variables); err != nil {
		return "", "", err
	}

	return subjectBuf.String(), bodyBuf.String(), nil
}

// EmailChannel implementation
func (ec *EmailChannel) Send(notification *Notification) error {
	// Create message
	msg := fmt.Sprintf("To: %s\r\n", notification.Recipient)
	msg += fmt.Sprintf("Subject: %s\r\n", notification.Subject)
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += notification.Body

	// Setup authentication
	auth := smtp.PlainAuth("", ec.config.Username, ec.config.Password, ec.config.Host)

	// Setup TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         ec.config.Host,
	}

	// Connect to server
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ec.config.Host, ec.config.Port), tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, ec.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return err
	}

	// Send email
	if err := client.Mail(ec.config.From); err != nil {
		return err
	}

	if err := client.Rcpt(notification.Recipient); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return err
	}

	return writer.Close()
}

func (ec *EmailChannel) GetType() NotificationType {
	return NotificationEmail
}

// SMSChannel implementation (placeholder - would integrate with SMS provider)
func (sc *SMSChannel) Send(notification *Notification) error {
	// This is a placeholder implementation
	// In a real application, you would integrate with an SMS provider like Twilio
	sc.logger.Info("SMS notification sent (placeholder)", map[string]interface{}{
		"recipient": notification.Recipient,
		"message":   notification.Body,
	})
	return nil
}

func (sc *SMSChannel) GetType() NotificationType {
	return NotificationSMS
}

// InAppChannel implementation
func (iac *InAppChannel) Send(notification *Notification) error {
	// For in-app notifications, we just update the database record
	// The frontend would poll or use websockets to get new notifications
	iac.logger.Info("In-app notification created", map[string]interface{}{
		"notification_id": notification.ID,
		"user_id":         notification.UserID,
	})
	return nil
}

func (iac *InAppChannel) GetType() NotificationType {
	return NotificationInApp
}

// SlackChannel implementation
func (sc *SlackChannel) Send(notification *Notification) error {
	payload := map[string]interface{}{
		"text": fmt.Sprintf("*%s*\n%s", notification.Subject, notification.Body),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(sc.webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (sc *SlackChannel) GetType() NotificationType {
	return NotificationSlack
}

// Helper function to get template by name (for backwards compatibility)
var template2 = template

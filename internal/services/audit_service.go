package services

import (
	"encoding/json"
	"go-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// AuditService handles audit logging functionality
type AuditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service instance
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// AuditAction defines types of auditable actions
type AuditAction string

const (
	ActionLogin        AuditAction = "login"
	ActionLogout       AuditAction = "logout"
	ActionCreate       AuditAction = "create"
	ActionUpdate       AuditAction = "update"
	ActionDelete       AuditAction = "delete"
	ActionView         AuditAction = "view"
	ActionPasswordReset AuditAction = "password_reset"
	ActionEmailVerify  AuditAction = "email_verify"
	ActionRoleChange   AuditAction = "role_change"
	ActionPermissionChange AuditAction = "permission_change"
	ActionFileUpload   AuditAction = "file_upload"
	ActionFileDownload AuditAction = "file_download"
	ActionSecurityEvent AuditAction = "security_event"
)

// AuditEventData represents structured data for audit events
type AuditEventData struct {
	EntityType   string      `json:"entity_type,omitempty"`
	EntityID     string      `json:"entity_id,omitempty"`
	Changes      interface{} `json:"changes,omitempty"`
	OldValues    interface{} `json:"old_values,omitempty"`
	NewValues    interface{} `json:"new_values,omitempty"`
	RequestID    string      `json:"request_id,omitempty"`
	SessionID    string      `json:"session_id,omitempty"`
	UserAgent    string      `json:"user_agent,omitempty"`
	RemoteAddr   string      `json:"remote_addr,omitempty"`
	Method       string      `json:"method,omitempty"`
	Path         string      `json:"path,omitempty"`
	StatusCode   int         `json:"status_code,omitempty"`
	Duration     string      `json:"duration,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

// LogEvent creates an audit log entry
func (s *AuditService) LogEvent(userID uint, action AuditAction, data AuditEventData) error {
	oldValuesJSON, _ := json.Marshal(data.OldValues)
	newValuesJSON, _ := json.Marshal(data.NewValues)
	metadataJSON, _ := json.Marshal(map[string]interface{}{
		"request_id":    data.RequestID,
		"session_id":    data.SessionID,
		"method":        data.Method,
		"path":          data.Path,
		"status_code":   data.StatusCode,
		"duration":      data.Duration,
		"error_message": data.ErrorMessage,
	})

	auditLog := &models.AuditLog{
		UserID:     &userID,
		Action:     string(action),
		Resource:   data.EntityType,
		OldValues:  string(oldValuesJSON),
		NewValues:  string(newValuesJSON),
		IPAddress:  data.RemoteAddr,
		UserAgent:  data.UserAgent,
		Metadata:   string(metadataJSON),
		CreatedAt:  time.Now(),
	}

	return s.db.Create(auditLog).Error
}

// LogSystemEvent creates an audit log entry for system events (without user)
func (s *AuditService) LogSystemEvent(action AuditAction, data AuditEventData) error {
	metadataJSON, _ := json.Marshal(map[string]interface{}{
		"request_id":    data.RequestID,
		"session_id":    data.SessionID,
		"method":        data.Method,
		"path":          data.Path,
		"status_code":   data.StatusCode,
		"duration":      data.Duration,
		"error_message": data.ErrorMessage,
	})

	auditLog := &models.AuditLog{
		Action:    string(action),
		Resource:  data.EntityType,
		IPAddress: data.RemoteAddr,
		UserAgent: data.UserAgent,
		Metadata:  string(metadataJSON),
		CreatedAt: time.Now(),
	}

	return s.db.Create(auditLog).Error
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (s *AuditService) GetUserAuditLogs(userID uint, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetSystemAuditLogs retrieves system-wide audit logs
func (s *AuditService) GetSystemAuditLogs(limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByAction retrieves audit logs by action type
func (s *AuditService) GetAuditLogsByAction(action AuditAction, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.Where("action = ?", string(action)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (s *AuditService) GetAuditLogsByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// DeleteOldAuditLogs removes audit logs older than specified days
func (s *AuditService) DeleteOldAuditLogs(daysToKeep int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	return s.db.Where("created_at < ?", cutoffDate).Delete(&models.AuditLog{}).Error
}

// GetAuditLogStats returns statistics about audit logs
func (s *AuditService) GetAuditLogStats() (map[string]interface{}, error) {
	var stats map[string]interface{}
	stats = make(map[string]interface{})

	// Total count
	var totalCount int64
	if err := s.db.Model(&models.AuditLog{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_count"] = totalCount

	// Count by action
	var actionCounts []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	if err := s.db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Find(&actionCounts).Error; err != nil {
		return nil, err
	}
	stats["action_counts"] = actionCounts

	// Recent activity (last 24 hours)
	yesterday := time.Now().Add(-24 * time.Hour)
	var recentCount int64
	if err := s.db.Model(&models.AuditLog{}).
		Where("created_at > ?", yesterday).
		Count(&recentCount).Error; err != nil {
		return nil, err
	}
	stats["recent_activity"] = recentCount

	return stats, nil
}
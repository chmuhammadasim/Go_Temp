package models

import (
	"time"

	"gorm.io/gorm"
)

// EmailVerification represents email verification tokens
type EmailVerification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Email     string         `json:"email" gorm:"not null"`
	Token     string         `json:"token" gorm:"not null;uniqueIndex"`
	Type      string         `json:"type" gorm:"not null"` // verification, reset, otp
	Code      string         `json:"code,omitempty"`       // For OTP
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time     `json:"used_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// IsExpired checks if the verification token is expired
func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.ExpiresAt)
}

// IsUsed checks if the verification token has been used
func (ev *EmailVerification) IsUsed() bool {
	return ev.UsedAt != nil
}

// MarkAsUsed marks the verification as used
func (ev *EmailVerification) MarkAsUsed() {
	now := time.Now()
	ev.UsedAt = &now
}

// AuditLog represents system audit logs
type AuditLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     *uint     `json:"user_id,omitempty" gorm:"index"`
	Action     string    `json:"action" gorm:"not null;index"`   // CREATE, UPDATE, DELETE, LOGIN, etc.
	Resource   string    `json:"resource" gorm:"not null;index"` // user, post, comment, etc.
	ResourceID *uint     `json:"resource_id,omitempty" gorm:"index"`
	OldValues  string    `json:"old_values,omitempty" gorm:"type:jsonb"`
	NewValues  string    `json:"new_values,omitempty" gorm:"type:jsonb"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Metadata   string    `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SecurityEvent represents security-related events
type SecurityEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      *uint     `json:"user_id,omitempty" gorm:"index"`
	EventType   string    `json:"event_type" gorm:"not null;index"` // login_failed, account_locked, suspicious_activity
	Severity    string    `json:"severity" gorm:"not null;index"`   // low, medium, high, critical
	Description string    `json:"description" gorm:"not null"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Metadata    string    `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserSession represents active user sessions
type UserSession struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	LastSeen  time.Time      `json:"last_seen"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Permission represents system permissions
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Resource    string         `json:"resource" gorm:"not null"` // user, post, comment, etc.
	Action      string         `json:"action" gorm:"not null"`   // create, read, update, delete
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Many-to-many relationships
	Roles []*Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// RolePermission represents the junction table for roles and permissions
type RolePermission struct {
	RoleID       uint `json:"role_id" gorm:"primaryKey"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey"`
}

// UserLoginAttempt tracks login attempts for security
type UserLoginAttempt struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null;index"`
	IPAddress string    `json:"ip_address"`
	Success   bool      `json:"success"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// FileUpload represents uploaded files
type FileUpload struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null;index"`
	OriginalName  string         `json:"original_name" gorm:"not null"`
	FileName      string         `json:"file_name" gorm:"not null;uniqueIndex"`
	FilePath      string         `json:"file_path" gorm:"not null"`
	FileSize      int64          `json:"file_size"`
	MimeType      string         `json:"mime_type"`
	FileType      string         `json:"file_type" gorm:"index"` // image, document, video, etc.
	IsPublic      bool           `json:"is_public" gorm:"default:false"`
	DownloadCount int            `json:"download_count" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Notification represents system notifications
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Type      string         `json:"type" gorm:"not null;index"` // email, sms, push, in_app
	Title     string         `json:"title" gorm:"not null"`
	Message   string         `json:"message" gorm:"not null"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	Data      string         `json:"data,omitempty" gorm:"type:jsonb"`
	SentAt    *time.Time     `json:"sent_at,omitempty"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now
}

// SystemSetting represents configurable system settings
type SystemSetting struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Key         string         `json:"key" gorm:"uniqueIndex;not null"`
	Value       string         `json:"value" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"` // string, int, bool, json
	Description string         `json:"description"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// APIKey represents API keys for external access
type APIKey struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Key         string         `json:"key" gorm:"uniqueIndex;not null"`
	Permissions string         `json:"permissions" gorm:"type:jsonb"` // JSON array of permissions
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	LastUsed    *time.Time     `json:"last_used,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TwoFactorAuth represents 2FA settings for users
type TwoFactorAuth struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	Secret      string         `json:"secret" gorm:"not null"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:false"`
	BackupCodes string         `json:"backup_codes,omitempty" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Request/Response DTOs for new features

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest represents resend verification request
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

// FileUploadResponse represents file upload response
type FileUploadResponse struct {
	ID           uint   `json:"id"`
	OriginalName string `json:"original_name"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	FileType     string `json:"file_type"`
	URL          string `json:"url"`
}

// NotificationResponse represents notification response
type NotificationResponse struct {
	ID        uint       `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

// Enable2FARequest represents 2FA setup request
type Enable2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// Verify2FARequest represents 2FA verification request
type Verify2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// AuditLogResponse represents audit log response
type AuditLogResponse struct {
	ID         uint      `json:"id"`
	UserID     *uint     `json:"user_id,omitempty"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID *uint     `json:"resource_id,omitempty"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page     int    `form:"page,default=1" validate:"min=1"`
	Limit    int    `form:"limit,default=10" validate:"min=1,max=100"`
	Sort     string `form:"sort,default=created_at"`
	Order    string `form:"order,default=desc" validate:"oneof=asc desc"`
	Search   string `form:"search"`
	Filter   string `form:"filter"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// GetTotalPages calculates total pages
func (p *PaginationResponse) GetTotalPages() int {
	if p.Limit == 0 {
		return 0
	}
	return int((p.Total + int64(p.Limit) - 1) / int64(p.Limit))
}

// SetPagination sets pagination values
func (p *PaginationResponse) SetPagination(page, limit int, total int64) {
	p.Page = page
	p.Limit = limit
	p.Total = total
	p.TotalPages = p.GetTotalPages()
	p.HasNext = page < p.TotalPages
	p.HasPrev = page > 1
}

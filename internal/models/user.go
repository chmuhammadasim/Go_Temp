package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
)

// User represents a user in the system
type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username  string `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Password  string `json:"-" gorm:"not null" validate:"required,min=6"`
	FirstName string `json:"first_name" gorm:"not null" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" gorm:"not null" validate:"required,min=1,max=50"`
	Role      Role   `json:"role" gorm:"not null;default:'user'" validate:"required,oneof=admin moderator user"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`

	// Enhanced security fields
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	PhoneNumber     string     `json:"phone_number,omitempty" gorm:"uniqueIndex"`
	PhoneVerified   bool       `json:"phone_verified" gorm:"default:false"`
	Avatar          string     `json:"avatar,omitempty"`
	Timezone        string     `json:"timezone" gorm:"default:'UTC'"`
	Language        string     `json:"language" gorm:"default:'en'"`

	// Security fields
	FailedLoginAttempts int        `json:"-" gorm:"default:0"`
	AccountLockedUntil  *time.Time `json:"-"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP         string     `json:"-"`
	PasswordChangedAt   *time.Time `json:"-"`
	MustChangePassword  bool       `json:"-" gorm:"default:false"`
	TwoFactorEnabled    bool       `json:"two_factor_enabled" gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts              []Post              `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Comments           []Comment           `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	EmailVerifications []EmailVerification `json:"-" gorm:"foreignKey:UserID"`
	Sessions           []UserSession       `json:"-" gorm:"foreignKey:UserID"`
	FileUploads        []FileUpload        `json:"-" gorm:"foreignKey:UserID"`
	Notifications      []Notification      `json:"-" gorm:"foreignKey:UserID"`
	APIKeys            []APIKey            `json:"-" gorm:"foreignKey:UserID"`
	TwoFactorAuth      *TwoFactorAuth      `json:"-" gorm:"foreignKey:UserID"`
	AuditLogs          []AuditLog          `json:"-" gorm:"foreignKey:UserID"`
}

// UserCreateRequest represents the request payload for creating a user
type UserCreateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Role      Role   `json:"role,omitempty" validate:"omitempty,oneof=admin moderator user"`
}

// UserUpdateRequest represents the request payload for updating a user
type UserUpdateRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1,max=50"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1,max=50"`
	Role      *Role   `json:"role,omitempty" validate:"omitempty,oneof=admin moderator user"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID               uint       `json:"id"`
	Email            string     `json:"email"`
	Username         string     `json:"username"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Role             Role       `json:"role"`
	IsActive         bool       `json:"is_active"`
	EmailVerified    bool       `json:"email_verified"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty"`
	PhoneNumber      string     `json:"phone_number,omitempty"`
	PhoneVerified    bool       `json:"phone_verified"`
	Avatar           string     `json:"avatar,omitempty"`
	Timezone         string     `json:"timezone"`
	Language         string     `json:"language"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	LastLoginAt      *time.Time `json:"last_login_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response payload for user login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Post represents a blog post or article
type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
	Content   string         `json:"content" gorm:"type:text" validate:"required,min=1"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null"`
	Published bool           `json:"published" gorm:"default:false"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
}

// PostCreateRequest represents the request payload for creating a post
type PostCreateRequest struct {
	Title     string `json:"title" validate:"required,min=1,max=200"`
	Content   string `json:"content" validate:"required,min=1"`
	Published bool   `json:"published,omitempty"`
}

// PostUpdateRequest represents the request payload for updating a post
type PostUpdateRequest struct {
	Title     *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Content   *string `json:"content,omitempty" validate:"omitempty,min=1"`
	Published *bool   `json:"published,omitempty"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null" validate:"required,min=1"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
	Post Post `json:"post" gorm:"foreignKey:PostID"`
}

// CommentCreateRequest represents the request payload for creating a comment
type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1"`
	PostID  uint   `json:"post_id" validate:"required"`
}

// CommentUpdateRequest represents the request payload for updating a comment
type CommentUpdateRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,min=1"`
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash password before saving
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	// Set default role if not provided
	if u.Role == "" {
		u.Role = RoleUser
	}

	return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:               u.ID,
		Email:            u.Email,
		Username:         u.Username,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Role:             u.Role,
		IsActive:         u.IsActive,
		EmailVerified:    u.EmailVerified,
		EmailVerifiedAt:  u.EmailVerifiedAt,
		PhoneNumber:      u.PhoneNumber,
		PhoneVerified:    u.PhoneVerified,
		Avatar:           u.Avatar,
		Timezone:         u.Timezone,
		Language:         u.Language,
		TwoFactorEnabled: u.TwoFactorEnabled,
		LastLoginAt:      u.LastLoginAt,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsModerator checks if the user has moderator role
func (u *User) IsModerator() bool {
	return u.Role == RoleModerator
}

// CanModerate checks if the user can moderate (admin or moderator)
func (u *User) CanModerate() bool {
	return u.Role == RoleAdmin || u.Role == RoleModerator
}

// HasPermission checks if the user has specific permissions
func (u *User) HasPermission(requiredRole Role) bool {
	switch requiredRole {
	case RoleAdmin:
		return u.Role == RoleAdmin
	case RoleModerator:
		return u.Role == RoleAdmin || u.Role == RoleModerator
	case RoleUser:
		return true // All authenticated users have user permission
	default:
		return false
	}
}

// Security-related methods

// IsAccountLocked checks if the account is currently locked
func (u *User) IsAccountLocked() bool {
	if u.AccountLockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.AccountLockedUntil)
}

// LockAccount locks the user account for the specified duration
func (u *User) LockAccount(duration time.Duration) {
	lockUntil := time.Now().Add(duration)
	u.AccountLockedUntil = &lockUntil
}

// UnlockAccount unlocks the user account
func (u *User) UnlockAccount() {
	u.AccountLockedUntil = nil
	u.FailedLoginAttempts = 0
}

// IncrementFailedAttempts increments failed login attempts
func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++
}

// ResetFailedAttempts resets failed login attempts to 0
func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
}

// ShouldLockAccount checks if account should be locked based on failed attempts
func (u *User) ShouldLockAccount(maxAttempts int) bool {
	return u.FailedLoginAttempts >= maxAttempts
}

// UpdateLastLogin updates the last login timestamp and IP
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = ip
	u.ResetFailedAttempts()
}

// MarkEmailAsVerified marks the email as verified
func (u *User) MarkEmailAsVerified() {
	now := time.Now()
	u.EmailVerified = true
	u.EmailVerifiedAt = &now
}

// MarkPhoneAsVerified marks the phone as verified
func (u *User) MarkPhoneAsVerified() {
	u.PhoneVerified = true
}

// UpdatePassword updates the user's password and sets the password changed timestamp
func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	u.Password = string(hashedPassword)
	u.PasswordChangedAt = &now
	u.MustChangePassword = false

	return nil
}

// SetMustChangePassword forces the user to change password on next login
func (u *User) SetMustChangePassword() {
	u.MustChangePassword = true
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// GetDisplayName returns a display-friendly name
func (u *User) GetDisplayName() string {
	if u.FirstName != "" {
		return u.FirstName
	}
	return u.Username
}

// CanAccessResource checks if user can access a specific resource
func (u *User) CanAccessResource(resourceUserID uint) bool {
	// Admin can access everything
	if u.IsAdmin() {
		return true
	}
	// User can access their own resources
	return u.ID == resourceUserID
}

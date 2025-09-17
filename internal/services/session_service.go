package services

import (
	"crypto/rand"
	"encoding/hex"
	"go-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// SessionService handles user session management
type SessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new session service instance
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

// CreateSession creates a new user session
func (s *SessionService) CreateSession(userID uint, ipAddress, userAgent string) (*models.UserSession, error) {
	sessionToken, err := s.generateSessionToken()
	if err != nil {
		return nil, err
	}

	session := &models.UserSession{
		ID:        sessionToken,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		IsActive:  true,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession validates a session token and returns the session
func (s *SessionService) ValidateSession(token string) (*models.UserSession, error) {
	var session models.UserSession
	err := s.db.Where("id = ? AND is_active = ? AND expires_at > ?", 
		token, true, time.Now()).
		Preload("User").
		First(&session).Error
	
	if err != nil {
		return nil, err
	}

	// Update last seen time
	session.LastSeen = time.Now()
	session.UpdatedAt = time.Now()
	s.db.Save(&session)

	return &session, nil
}

// RefreshSession extends the session expiry time
func (s *SessionService) RefreshSession(token string) error {
	return s.db.Model(&models.UserSession{}).
		Where("id = ? AND is_active = ?", token, true).
		Updates(map[string]interface{}{
			"expires_at": time.Now().Add(24 * time.Hour),
			"updated_at": time.Now(),
		}).Error
}

// InvalidateSession deactivates a specific session
func (s *SessionService) InvalidateSession(token string) error {
	return s.db.Model(&models.UserSession{}).
		Where("id = ?", token).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// InvalidateUserSessions deactivates all sessions for a user
func (s *SessionService) InvalidateUserSessions(userID uint) error {
	return s.db.Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// InvalidateUserSessionsExcept deactivates all sessions for a user except the specified one
func (s *SessionService) InvalidateUserSessionsExcept(userID uint, exceptToken string) error {
	return s.db.Model(&models.UserSession{}).
		Where("user_id = ? AND id != ? AND is_active = ?", userID, exceptToken, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// GetUserSessions retrieves all active sessions for a user
func (s *SessionService) GetUserSessions(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := s.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

// GetAllUserSessions retrieves all sessions (active and inactive) for a user
func (s *SessionService) GetAllUserSessions(userID uint, limit, offset int) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions() error {
	return s.db.Where("expires_at < ? OR is_active = ?", time.Now(), false).
		Delete(&models.UserSession{}).Error
}

// GetSessionStats returns session statistics
func (s *SessionService) GetSessionStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total active sessions
	var activeCount int64
	if err := s.db.Model(&models.UserSession{}).
		Where("is_active = ? AND expires_at > ?", true, time.Now()).
		Count(&activeCount).Error; err != nil {
		return nil, err
	}
	stats["active_sessions"] = activeCount

	// Total sessions today
	today := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	if err := s.db.Model(&models.UserSession{}).
		Where("created_at >= ?", today).
		Count(&todayCount).Error; err != nil {
		return nil, err
	}
	stats["sessions_today"] = todayCount

	// Unique users with active sessions
	var uniqueUsers int64
	if err := s.db.Model(&models.UserSession{}).
		Where("is_active = ? AND expires_at > ?", true, time.Now()).
		Distinct("user_id").
		Count(&uniqueUsers).Error; err != nil {
		return nil, err
	}
	stats["unique_active_users"] = uniqueUsers

	return stats, nil
}

// IsUserSessionActive checks if a user has any active sessions
func (s *SessionService) IsUserSessionActive(userID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", 
			userID, true, time.Now()).
		Count(&count).Error
	return count > 0, err
}

// GetSessionByID retrieves a session by its ID
func (s *SessionService) GetSessionByID(sessionID string) (*models.UserSession, error) {
	var session models.UserSession
	err := s.db.Where("id = ?", sessionID).
		Preload("User").
		First(&session).Error
	return &session, err
}

// generateSessionToken generates a secure random session token
func (s *SessionService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// UpdateSessionActivity updates the session's last activity timestamp
func (s *SessionService) UpdateSessionActivity(token string, ipAddress string) error {
	updates := map[string]interface{}{
		"last_seen":  time.Now(),
		"updated_at": time.Now(),
	}
	
	// Update IP address if it has changed
	if ipAddress != "" {
		updates["ip_address"] = ipAddress
	}

	return s.db.Model(&models.UserSession{}).
		Where("id = ? AND is_active = ?", token, true).
		Updates(updates).Error
}
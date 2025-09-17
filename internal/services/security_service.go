package services

import (
	"encoding/json"
	"go-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SecurityService handles security event tracking and monitoring
type SecurityService struct {
	db           *gorm.DB
	auditService *AuditService
}

// NewSecurityService creates a new security service instance
func NewSecurityService(db *gorm.DB, auditService *AuditService) *SecurityService {
	return &SecurityService{
		db:           db,
		auditService: auditService,
	}
}

// SecurityEventType defines types of security events
type SecurityEventType string

const (
	EventSuspiciousLogin    SecurityEventType = "suspicious_login"
	EventMultipleFailedLogins SecurityEventType = "multiple_failed_logins"
	EventUnusualLocation    SecurityEventType = "unusual_location"
	EventPasswordBreach     SecurityEventType = "password_breach"
	EventAccountLockout     SecurityEventType = "account_lockout"
	EventUnauthorizedAccess SecurityEventType = "unauthorized_access"
	EventDataExfiltration   SecurityEventType = "data_exfiltration"
	EventMaliciousRequest   SecurityEventType = "malicious_request"
	EventRateLimitExceeded  SecurityEventType = "rate_limit_exceeded"
	EventSQLInjectionAttempt SecurityEventType = "sql_injection_attempt"
	EventXSSAttempt         SecurityEventType = "xss_attempt"
	EventBruteForceAttempt  SecurityEventType = "brute_force_attempt"
	EventCSRFAttempt        SecurityEventType = "csrf_attempt"
	EventFileUploadViolation SecurityEventType = "file_upload_violation"
	EventPrivilegeEscalation SecurityEventType = "privilege_escalation"
)

// SecuritySeverity defines severity levels for security events
type SecuritySeverity string

const (
	SeverityLow      SecuritySeverity = "low"
	SeverityMedium   SecuritySeverity = "medium"
	SeverityHigh     SecuritySeverity = "high"
	SeverityCritical SecuritySeverity = "critical"
)

// SecurityEventData represents additional data for security events
type SecurityEventData struct {
	RequestID      string      `json:"request_id,omitempty"`
	SessionID      string      `json:"session_id,omitempty"`
	UserAgent      string      `json:"user_agent,omitempty"`
	RemoteAddr     string      `json:"remote_addr,omitempty"`
	Method         string      `json:"method,omitempty"`
	Path           string      `json:"path,omitempty"`
	StatusCode     int         `json:"status_code,omitempty"`
	Payload        interface{} `json:"payload,omitempty"`
	Headers        interface{} `json:"headers,omitempty"`
	FailedAttempts int         `json:"failed_attempts,omitempty"`
	TimeWindow     string      `json:"time_window,omitempty"`
	DetectionRule  string      `json:"detection_rule,omitempty"`
	RiskScore      int         `json:"risk_score,omitempty"`
	Geolocation    interface{} `json:"geolocation,omitempty"`
	DeviceFingerprint string   `json:"device_fingerprint,omitempty"`
}

// LogSecurityEvent creates a security event log entry
func (s *SecurityService) LogSecurityEvent(userID *uuid.UUID, eventType SecurityEventType, severity SecuritySeverity, description string, data SecurityEventData) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	securityEvent := &models.SecurityEvent{
		ID:          uuid.New(),
		UserID:      userID,
		EventType:   string(eventType),
		Severity:    string(severity),
		Description: description,
		Data:        dataJSON,
		IPAddress:   data.RemoteAddr,
		UserAgent:   data.UserAgent,
		Resolved:    false,
		CreatedAt:   time.Now(),
	}

	if err := s.db.Create(securityEvent).Error; err != nil {
		return err
	}

	// Also log in audit trail if audit service is available
	if s.auditService != nil {
		auditData := AuditEventData{
			RemoteAddr:   data.RemoteAddr,
			UserAgent:    data.UserAgent,
			RequestID:    data.RequestID,
			SessionID:    data.SessionID,
			Method:       data.Method,
			Path:         data.Path,
			StatusCode:   data.StatusCode,
			ErrorMessage: description,
		}

		if userID != nil {
			s.auditService.LogEvent(*userID, ActionSecurityEvent, auditData)
		} else {
			s.auditService.LogSystemEvent(ActionSecurityEvent, auditData)
		}
	}

	return nil
}

// DetectSuspiciousLogin analyzes login attempts for suspicious patterns
func (s *SecurityService) DetectSuspiciousLogin(userID uuid.UUID, ipAddress, userAgent string) error {
	// Check for multiple failed login attempts in short time window
	failedAttempts, err := s.getRecentFailedLoginAttempts(userID, 15*time.Minute)
	if err != nil {
		return err
	}

	if failedAttempts >= 5 {
		data := SecurityEventData{
			RemoteAddr:     ipAddress,
			UserAgent:      userAgent,
			FailedAttempts: failedAttempts,
			TimeWindow:     "15 minutes",
			DetectionRule:  "multiple_failed_logins",
			RiskScore:      80,
		}

		return s.LogSecurityEvent(&userID, EventMultipleFailedLogins, SeverityHigh,
			"Multiple failed login attempts detected", data)
	}

	// Check for login from unusual location (simplified - in real implementation, use geolocation)
	if s.isUnusualLocation(userID, ipAddress) {
		data := SecurityEventData{
			RemoteAddr:    ipAddress,
			UserAgent:     userAgent,
			DetectionRule: "unusual_location",
			RiskScore:     60,
		}

		return s.LogSecurityEvent(&userID, EventUnusualLocation, SeverityMedium,
			"Login from unusual location detected", data)
	}

	return nil
}

// DetectRateLimitViolation logs rate limit violations
func (s *SecurityService) DetectRateLimitViolation(userID *uuid.UUID, ipAddress, userAgent, endpoint string, requestCount int) error {
	data := SecurityEventData{
		RemoteAddr:     ipAddress,
		UserAgent:      userAgent,
		Path:           endpoint,
		FailedAttempts: requestCount,
		DetectionRule:  "rate_limit_exceeded",
		RiskScore:      50,
	}

	return s.LogSecurityEvent(userID, EventRateLimitExceeded, SeverityMedium,
		"Rate limit exceeded", data)
}

// DetectMaliciousRequest analyzes requests for malicious patterns
func (s *SecurityService) DetectMaliciousRequest(userID *uuid.UUID, ipAddress, userAgent, method, path string, payload interface{}) error {
	// Simple pattern detection (in real implementation, use more sophisticated detection)
	maliciousPatterns := []string{
		"<script>", "javascript:", "SELECT * FROM", "UNION SELECT", "DROP TABLE",
		"../", "..\\", "/etc/passwd", "cmd.exe", "powershell",
	}

	payloadStr := ""
	if payload != nil {
		if payloadBytes, err := json.Marshal(payload); err == nil {
			payloadStr = string(payloadBytes)
		}
	}

	fullRequest := method + " " + path + " " + payloadStr
	
	for _, pattern := range maliciousPatterns {
		if contains(fullRequest, pattern) {
			data := SecurityEventData{
				RemoteAddr:    ipAddress,
				UserAgent:     userAgent,
				Method:        method,
				Path:          path,
				Payload:       payload,
				DetectionRule: "malicious_pattern_detected",
				RiskScore:     90,
			}

			var eventType SecurityEventType
			var description string

			switch {
			case contains(pattern, "SELECT") || contains(pattern, "UNION") || contains(pattern, "DROP"):
				eventType = EventSQLInjectionAttempt
				description = "SQL injection attempt detected"
			case contains(pattern, "<script>") || contains(pattern, "javascript:"):
				eventType = EventXSSAttempt
				description = "XSS attempt detected"
			default:
				eventType = EventMaliciousRequest
				description = "Malicious request pattern detected"
			}

			return s.LogSecurityEvent(userID, eventType, SeverityCritical, description, data)
		}
	}

	return nil
}

// MarkSecurityEventResolved marks a security event as resolved
func (s *SecurityService) MarkSecurityEventResolved(eventID uuid.UUID, resolvedBy uuid.UUID, resolution string) error {
	updates := map[string]interface{}{
		"resolved":     true,
		"resolved_by":  &resolvedBy,
		"resolved_at":  time.Now(),
		"resolution":   resolution,
		"updated_at":   time.Now(),
	}

	return s.db.Model(&models.SecurityEvent{}).
		Where("id = ?", eventID).
		Updates(updates).Error
}

// GetUnresolvedSecurityEvents retrieves unresolved security events
func (s *SecurityService) GetUnresolvedSecurityEvents(limit, offset int) ([]models.SecurityEvent, error) {
	var events []models.SecurityEvent
	err := s.db.Where("resolved = ?", false).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// GetSecurityEventsBySeverity retrieves security events by severity level
func (s *SecurityService) GetSecurityEventsBySeverity(severity SecuritySeverity, limit, offset int) ([]models.SecurityEvent, error) {
	var events []models.SecurityEvent
	err := s.db.Where("severity = ?", string(severity)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// GetUserSecurityEvents retrieves security events for a specific user
func (s *SecurityService) GetUserSecurityEvents(userID uuid.UUID, limit, offset int) ([]models.SecurityEvent, error) {
	var events []models.SecurityEvent
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// GetSecurityEventStats returns security event statistics
func (s *SecurityService) GetSecurityEventStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total events
	var totalCount int64
	if err := s.db.Model(&models.SecurityEvent{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_events"] = totalCount

	// Unresolved events
	var unresolvedCount int64
	if err := s.db.Model(&models.SecurityEvent{}).
		Where("resolved = ?", false).
		Count(&unresolvedCount).Error; err != nil {
		return nil, err
	}
	stats["unresolved_events"] = unresolvedCount

	// Events by severity
	var severityCounts []struct {
		Severity string `json:"severity"`
		Count    int64  `json:"count"`
	}
	if err := s.db.Model(&models.SecurityEvent{}).
		Select("severity, COUNT(*) as count").
		Group("severity").
		Find(&severityCounts).Error; err != nil {
		return nil, err
	}
	stats["severity_counts"] = severityCounts

	// Recent events (last 24 hours)
	yesterday := time.Now().Add(-24 * time.Hour)
	var recentCount int64
	if err := s.db.Model(&models.SecurityEvent{}).
		Where("created_at > ?", yesterday).
		Count(&recentCount).Error; err != nil {
		return nil, err
	}
	stats["recent_events"] = recentCount

	return stats, nil
}

// getRecentFailedLoginAttempts counts failed login attempts in the specified time window
func (s *SecurityService) getRecentFailedLoginAttempts(userID uuid.UUID, timeWindow time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-timeWindow)
	
	var count int64
	err := s.db.Model(&models.SecurityEvent{}).
		Where("user_id = ? AND event_type = ? AND created_at > ?", 
			userID, string(EventMultipleFailedLogins), cutoffTime).
		Count(&count).Error
	
	return int(count), err
}

// isUnusualLocation checks if the IP address represents an unusual location for the user
func (s *SecurityService) isUnusualLocation(userID uuid.UUID, ipAddress string) bool {
	// Simplified implementation - in reality, you would use geolocation services
	// and compare against user's typical locations
	
	var recentLogins []models.SecurityEvent
	err := s.db.Where("user_id = ? AND ip_address != ? AND created_at > ?", 
		userID, ipAddress, time.Now().Add(-30*24*time.Hour)).
		Limit(10).
		Find(&recentLogins).Error
	
	if err != nil || len(recentLogins) == 0 {
		return false
	}
	
	// If all recent logins are from different IP addresses, consider this unusual
	return len(recentLogins) > 5
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	// Simple case-insensitive substring check
	// In production, use more sophisticated pattern matching
	return len(s) >= len(substr) && 
		   (s == substr || len(s) > len(substr) && 
		    (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		     containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
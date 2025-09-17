package services

import (
	"crypto/rand"
	"fmt"
	"go-backend/internal/models"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TwoFactorService handles two-factor authentication functionality
type TwoFactorService struct {
	db           *gorm.DB
	emailService *EmailService
}

// NewTwoFactorService creates a new two-factor authentication service instance
func NewTwoFactorService(db *gorm.DB, emailService *EmailService) *TwoFactorService {
	return &TwoFactorService{
		db:           db,
		emailService: emailService,
	}
}

// TwoFactorMethod represents different 2FA methods
type TwoFactorMethod string

const (
	TwoFactorMethodEmail TwoFactorMethod = "email"
	TwoFactorMethodSMS   TwoFactorMethod = "sms"
	TwoFactorMethodTOTP  TwoFactorMethod = "totp"
)

// GenerateEmailOTP generates and sends an OTP via email
func (s *TwoFactorService) GenerateEmailOTP(userID uuid.UUID, email string) (string, error) {
	// Generate 6-digit OTP
	otp, err := s.generateOTP(6)
	if err != nil {
		return "", err
	}

	// Store OTP in database
	if err := s.storeOTP(userID, otp, TwoFactorMethodEmail); err != nil {
		return "", err
	}

	// Send OTP via email
	if err := s.emailService.SendOTP(email, otp); err != nil {
		return "", err
	}

	return otp, nil
}

// GenerateSMSOTP generates and sends an OTP via SMS (placeholder - requires SMS service)
func (s *TwoFactorService) GenerateSMSOTP(userID uuid.UUID, phoneNumber string) (string, error) {
	// Generate 6-digit OTP
	otp, err := s.generateOTP(6)
	if err != nil {
		return "", err
	}

	// Store OTP in database
	if err := s.storeOTP(userID, otp, TwoFactorMethodSMS); err != nil {
		return "", err
	}

	// TODO: Implement SMS sending service
	// For now, just return the OTP (in production, this should send SMS)
	fmt.Printf("SMS OTP for %s: %s\n", phoneNumber, otp)

	return otp, nil
}

// VerifyOTP verifies the provided OTP against stored OTP
func (s *TwoFactorService) VerifyOTP(userID uuid.UUID, providedOTP string, method TwoFactorMethod) (bool, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return false, err
	}

	// Check if OTP matches and is not expired
	if user.TwoFactorCode == providedOTP && 
	   user.TwoFactorMethod == string(method) && 
	   user.TwoFactorExpiresAt != nil && 
	   user.TwoFactorExpiresAt.After(time.Now()) {
		
		// Clear the OTP after successful verification
		if err := s.clearOTP(userID); err != nil {
			return false, err
		}
		
		return true, nil
	}

	// Increment failed attempts
	if err := s.incrementFailedAttempts(userID); err != nil {
		return false, err
	}

	return false, nil
}

// EnableTwoFactor enables two-factor authentication for a user
func (s *TwoFactorService) EnableTwoFactor(userID uuid.UUID, method TwoFactorMethod) error {
	updates := map[string]interface{}{
		"two_factor_enabled": true,
		"two_factor_method":  string(method),
		"updated_at":         time.Now(),
	}

	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

// DisableTwoFactor disables two-factor authentication for a user
func (s *TwoFactorService) DisableTwoFactor(userID uuid.UUID) error {
	updates := map[string]interface{}{
		"two_factor_enabled":    false,
		"two_factor_method":     "",
		"two_factor_code":       "",
		"two_factor_expires_at": nil,
		"updated_at":            time.Now(),
	}

	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

// IsTwoFactorEnabled checks if two-factor authentication is enabled for a user
func (s *TwoFactorService) IsTwoFactorEnabled(userID uuid.UUID) (bool, error) {
	var user models.User
	if err := s.db.Select("two_factor_enabled").First(&user, userID).Error; err != nil {
		return false, err
	}
	return user.TwoFactorEnabled, nil
}

// GetTwoFactorMethod returns the two-factor method for a user
func (s *TwoFactorService) GetTwoFactorMethod(userID uuid.UUID) (TwoFactorMethod, error) {
	var user models.User
	if err := s.db.Select("two_factor_method").First(&user, userID).Error; err != nil {
		return "", err
	}
	return TwoFactorMethod(user.TwoFactorMethod), nil
}

// ResendOTP resends the OTP using the user's preferred method
func (s *TwoFactorService) ResendOTP(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if !user.TwoFactorEnabled {
		return fmt.Errorf("two-factor authentication is not enabled for this user")
	}

	method := TwoFactorMethod(user.TwoFactorMethod)
	switch method {
	case TwoFactorMethodEmail:
		_, err := s.GenerateEmailOTP(userID, user.Email)
		return err
	case TwoFactorMethodSMS:
		if user.PhoneNumber == nil {
			return fmt.Errorf("phone number not configured for SMS 2FA")
		}
		_, err := s.GenerateSMSOTP(userID, *user.PhoneNumber)
		return err
	default:
		return fmt.Errorf("unsupported two-factor method: %s", method)
	}
}

// IsOTPValid checks if there's a valid OTP for the user without verification
func (s *TwoFactorService) IsOTPValid(userID uuid.UUID) (bool, error) {
	var user models.User
	if err := s.db.Select("two_factor_code, two_factor_expires_at").First(&user, userID).Error; err != nil {
		return false, err
	}

	return user.TwoFactorCode != "" && 
		   user.TwoFactorExpiresAt != nil && 
		   user.TwoFactorExpiresAt.After(time.Now()), nil
}

// GetOTPExpiryTime returns when the current OTP expires
func (s *TwoFactorService) GetOTPExpiryTime(userID uuid.UUID) (*time.Time, error) {
	var user models.User
	if err := s.db.Select("two_factor_expires_at").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.TwoFactorExpiresAt, nil
}

// generateOTP generates a random numeric OTP of specified length
func (s *TwoFactorService) generateOTP(length int) (string, error) {
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%0*d", length, n), nil
}

// storeOTP stores the OTP in the database with expiry time
func (s *TwoFactorService) storeOTP(userID uuid.UUID, otp string, method TwoFactorMethod) error {
	expiryTime := time.Now().Add(10 * time.Minute) // OTP expires in 10 minutes
	
	updates := map[string]interface{}{
		"two_factor_code":       otp,
		"two_factor_method":     string(method),
		"two_factor_expires_at": &expiryTime,
		"updated_at":            time.Now(),
	}

	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

// clearOTP removes the OTP from the database
func (s *TwoFactorService) clearOTP(userID uuid.UUID) error {
	updates := map[string]interface{}{
		"two_factor_code":       "",
		"two_factor_expires_at": nil,
		"updated_at":            time.Now(),
	}

	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

// incrementFailedAttempts increments the failed login attempts counter
func (s *TwoFactorService) incrementFailedAttempts(userID uuid.UUID) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("failed_login_attempts", gorm.Expr("failed_login_attempts + 1")).
		Update("updated_at", time.Now()).Error
}

// ClearExpiredOTPs removes expired OTPs from all users
func (s *TwoFactorService) ClearExpiredOTPs() error {
	updates := map[string]interface{}{
		"two_factor_code":       "",
		"two_factor_expires_at": nil,
		"updated_at":            time.Now(),
	}

	return s.db.Model(&models.User{}).
		Where("two_factor_expires_at < ?", time.Now()).
		Updates(updates).Error
}
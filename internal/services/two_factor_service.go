package services

import (
	"crypto/rand"
	"fmt"
	"go-backend/internal/models"
	"math/big"
	"time"

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
func (s *TwoFactorService) GenerateEmailOTP(userID uint, email, username string) (string, error) {
	// Generate 6-digit OTP
	otp, err := s.generateOTP(6)
	if err != nil {
		return "", err
	}

	// Store OTP in EmailVerification table
	verification := &models.EmailVerification{
		UserID:    userID,
		Email:     email,
		Type:      "otp",
		Code:      otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(verification).Error; err != nil {
		return "", err
	}

	// Send OTP via email
	if err := s.emailService.SendOTPEmail(email, username, otp); err != nil {
		return "", err
	}

	return otp, nil
}

// GenerateSMSOTP generates and sends an OTP via SMS (placeholder - requires SMS service)
func (s *TwoFactorService) GenerateSMSOTP(userID uint, phoneNumber, username string) (string, error) {
	// Generate 6-digit OTP
	otp, err := s.generateOTP(6)
	if err != nil {
		return "", err
	}

	// Store OTP in EmailVerification table (repurposed for SMS OTP)
	verification := &models.EmailVerification{
		UserID:    userID,
		Email:     phoneNumber, // Using email field for phone number temporarily
		Type:      "sms_otp",
		Code:      otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(verification).Error; err != nil {
		return "", err
	}

	// TODO: Implement SMS sending service
	// For now, just return the OTP (in production, this should send SMS)
	fmt.Printf("SMS OTP for %s: %s\n", phoneNumber, otp)

	return otp, nil
}

// VerifyOTP verifies the provided OTP against stored OTP
func (s *TwoFactorService) VerifyOTP(userID uint, providedOTP string, method TwoFactorMethod) (bool, error) {
	var verification models.EmailVerification
	
	// Determine the verification type based on method
	verificationTypes := []string{"otp"}
	if method == TwoFactorMethodSMS {
		verificationTypes = []string{"sms_otp"}
	}

	// Find the most recent valid OTP for this user
	err := s.db.Where("user_id = ? AND type IN (?) AND code = ? AND expires_at > ? AND used_at IS NULL", 
		userID, verificationTypes, providedOTP, time.Now()).
		Order("created_at DESC").
		First(&verification).Error

	if err != nil {
		// Increment failed attempts
		s.incrementFailedAttempts(userID)
		return false, err
	}

	// Mark the OTP as used
	verification.MarkAsUsed()
	if err := s.db.Save(&verification).Error; err != nil {
		return false, err
	}

	return true, nil
}

// EnableTwoFactor enables two-factor authentication for a user
func (s *TwoFactorService) EnableTwoFactor(userID uint, method TwoFactorMethod) error {
	// Check if TwoFactorAuth record exists
	var twoFA models.TwoFactorAuth
	err := s.db.Where("user_id = ?", userID).First(&twoFA).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new TwoFactorAuth record
		secret, err := s.generateSecret(32)
		if err != nil {
			return err
		}

		twoFA = models.TwoFactorAuth{
			UserID:    userID,
			Secret:    secret,
			IsEnabled: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.db.Create(&twoFA).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// Update existing record
		twoFA.IsEnabled = true
		twoFA.UpdatedAt = time.Now()
		if err := s.db.Save(&twoFA).Error; err != nil {
			return err
		}
	}

	// Update user's two factor enabled flag
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("two_factor_enabled", true).Error
}

// DisableTwoFactor disables two-factor authentication for a user
func (s *TwoFactorService) DisableTwoFactor(userID uint) error {
	// Update TwoFactorAuth record
	err := s.db.Model(&models.TwoFactorAuth{}).
		Where("user_id = ?", userID).
		Update("is_enabled", false).Error

	if err != nil {
		return err
	}

	// Update user's two factor enabled flag
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("two_factor_enabled", false).Error
}

// IsTwoFactorEnabled checks if two-factor authentication is enabled for a user
func (s *TwoFactorService) IsTwoFactorEnabled(userID uint) (bool, error) {
	var user models.User
	if err := s.db.Select("two_factor_enabled").First(&user, userID).Error; err != nil {
		return false, err
	}
	return user.TwoFactorEnabled, nil
}

// ResendOTP resends the OTP using email method
func (s *TwoFactorService) ResendOTP(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if !user.TwoFactorEnabled {
		return fmt.Errorf("two-factor authentication is not enabled for this user")
	}

	// For now, always resend via email
	_, err := s.GenerateEmailOTP(userID, user.Email, user.Username)
	return err
}

// IsOTPValid checks if there's a valid OTP for the user without verification
func (s *TwoFactorService) IsOTPValid(userID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.EmailVerification{}).
		Where("user_id = ? AND type IN (?, ?) AND expires_at > ? AND used_at IS NULL", 
			userID, "otp", "sms_otp", time.Now()).
		Count(&count).Error
	
	return count > 0, err
}

// GetOTPExpiryTime returns when the current OTP expires
func (s *TwoFactorService) GetOTPExpiryTime(userID uint) (*time.Time, error) {
	var verification models.EmailVerification
	err := s.db.Where("user_id = ? AND type IN (?, ?) AND expires_at > ? AND used_at IS NULL", 
		userID, "otp", "sms_otp", time.Now()).
		Order("created_at DESC").
		First(&verification).Error
	
	if err != nil {
		return nil, err
	}
	
	return &verification.ExpiresAt, nil
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

// generateSecret generates a random secret for TOTP
func (s *TwoFactorService) generateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	
	return string(bytes), nil
}

// incrementFailedAttempts increments the failed login attempts counter
func (s *TwoFactorService) incrementFailedAttempts(userID uint) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("failed_login_attempts", gorm.Expr("failed_login_attempts + 1")).Error
}

// ClearExpiredOTPs removes expired OTPs from all users
func (s *TwoFactorService) ClearExpiredOTPs() error {
	return s.db.Where("type IN (?, ?) AND expires_at < ?", "otp", "sms_otp", time.Now()).
		Delete(&models.EmailVerification{}).Error
}
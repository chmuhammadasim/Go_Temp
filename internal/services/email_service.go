package services

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"strconv"
	"strings"
	"time"

	"go-backend/internal/config"
	"go-backend/pkg/logger"

	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

// EmailService handles email operations
type EmailService struct {
	config *config.Config
	logger *logger.Logger
	dialer *mail.Dialer
}

// EmailTemplate represents email template data
type EmailTemplate struct {
	Subject string
	Body    string
	IsHTML  bool
}

// VerificationEmail contains verification email data
type VerificationEmail struct {
	Username string
	Email    string
	Token    string
	Code     string
	Link     string
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config, logger *logger.Logger) *EmailService {
	dialer := mail.NewDialer(
		cfg.Email.Host,
		cfg.Email.Port,
		cfg.Email.Username,
		cfg.Email.Password,
	)
	
	if cfg.Email.TLS {
		dialer.StartTLSPolicy = mail.MandatoryStartTLS
	}

	return &EmailService{
		config: cfg,
		logger: logger,
		dialer: dialer,
	}
}

// SendEmail sends an email
func (e *EmailService) SendEmail(to, subject, body string, isHTML bool) error {
	m := mail.NewMessage()
	m.SetHeader("From", e.config.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	
	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	if err := e.dialer.DialAndSend(m); err != nil {
		e.logger.WithFields(logrus.Fields{
			"to":      to,
			"subject": subject,
			"error":   err.Error(),
		}).Error("Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully")

	return nil
}

// SendVerificationEmail sends email verification
func (e *EmailService) SendVerificationEmail(email, username, token string) error {
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", e.config.App.FrontendURL, token)
	
	tmplData := VerificationEmail{
		Username: username,
		Email:    email,
		Token:    token,
		Link:     verificationLink,
	}

	subject := "Verify Your Email Address"
	body := e.generateVerificationEmailHTML(tmplData)

	return e.SendEmail(email, subject, body, true)
}

// SendPasswordResetEmail sends password reset email
func (e *EmailService) SendPasswordResetEmail(email, username, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", e.config.App.FrontendURL, token)
	
	tmplData := VerificationEmail{
		Username: username,
		Email:    email,
		Token:    token,
		Link:     resetLink,
	}

	subject := "Reset Your Password"
	body := e.generatePasswordResetEmailHTML(tmplData)

	return e.SendEmail(email, subject, body, true)
}

// SendOTPEmail sends OTP code via email
func (e *EmailService) SendOTPEmail(email, username, otp string) error {
	tmplData := VerificationEmail{
		Username: username,
		Email:    email,
		Code:     otp,
	}

	subject := "Your Verification Code"
	body := e.generateOTPEmailHTML(tmplData)

	return e.SendEmail(email, subject, body, true)
}

// SendWelcomeEmail sends welcome email to new users
func (e *EmailService) SendWelcomeEmail(email, username string) error {
	tmplData := VerificationEmail{
		Username: username,
		Email:    email,
	}

	subject := "Welcome to Our Platform!"
	body := e.generateWelcomeEmailHTML(tmplData)

	return e.SendEmail(email, subject, body, true)
}

// GenerateOTP generates a 6-digit OTP
func (e *EmailService) GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	
	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}
	
	return strconv.Itoa(int(n.Add(n, min).Int64())), nil
}

// GenerateToken generates a secure random token
func (e *EmailService) GenerateToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	
	return string(b), nil
}

// Email template generators

func (e *EmailService) generateVerificationEmailHTML(data VerificationEmail) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Email Verification</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Username}}!</h2>
            <p>Thank you for registering with us. To complete your registration, please verify your email address by clicking the button below:</p>
            <a href="{{.Link}}" class="button">Verify Email Address</a>
            <p>If you can't click the button, copy and paste this link into your browser:</p>
            <p><a href="{{.Link}}">{{.Link}}</a></p>
            <p>This verification link will expire in 24 hours.</p>
            <p>If you didn't create an account with us, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>`
	
	t, _ := template.New("verification").Parse(tmpl)
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

func (e *EmailService) generatePasswordResetEmailHTML(data VerificationEmail) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #dc3545; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background: #dc3545; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Username}}!</h2>
            <p>We received a request to reset your password. Click the button below to create a new password:</p>
            <a href="{{.Link}}" class="button">Reset Password</a>
            <p>If you can't click the button, copy and paste this link into your browser:</p>
            <p><a href="{{.Link}}">{{.Link}}</a></p>
            <p>This reset link will expire in 1 hour.</p>
            <p>If you didn't request a password reset, please ignore this email and your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>`
	
	t, _ := template.New("reset").Parse(tmpl)
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

func (e *EmailService) generateOTPEmailHTML(data VerificationEmail) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your Verification Code</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; text-align: center; }
        .otp { font-size: 36px; font-weight: bold; color: #007bff; letter-spacing: 8px; margin: 20px 0; padding: 15px; background: white; border: 2px dashed #007bff; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verification Code</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Username}}!</h2>
            <p>Your verification code is:</p>
            <div class="otp">{{.Code}}</div>
            <p>This code will expire in 10 minutes.</p>
            <p>If you didn't request this code, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>`
	
	t, _ := template.New("otp").Parse(tmpl)
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

func (e *EmailService) generateWelcomeEmailHTML(data VerificationEmail) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome!</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome Aboard!</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Username}}!</h2>
            <p>Welcome to our platform! We're excited to have you on board.</p>
            <p>You can now enjoy all the features and benefits of your account.</p>
            <p>If you have any questions or need assistance, don't hesitate to contact our support team.</p>
            <p>Thank you for choosing us!</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>`
	
	t, _ := template.New("welcome").Parse(tmpl)
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}
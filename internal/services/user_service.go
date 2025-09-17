package services

import (
	"errors"
	"fmt"

	"go-backend/internal/models"
	"go-backend/internal/utils"

	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService struct {
	db         *gorm.DB
	jwtService *utils.JWTService
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB, jwtService *utils.JWTService) *UserService {
	return &UserService{
		db:         db,
		jwtService: jwtService,
	}
}

// Register creates a new user account
func (s *UserService) Register(req *models.UserCreateRequest) (*models.LoginResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		if existingUser.Email == req.Email {
			return nil, errors.New("user with this email already exists")
		}
		return nil, errors.New("user with this username already exists")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = models.RoleUser
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password, // Will be hashed by BeforeCreate hook
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      role,
		IsActive:  true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total users
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch users with pagination
	if err := s.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, total, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Update fields if provided
	updates := make(map[string]interface{})

	if req.Email != nil {
		// Check if email is already taken by another user
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", *req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already taken")
		}
		updates["email"] = *req.Email
	}

	if req.Username != nil {
		// Check if username is already taken by another user
		var existingUser models.User
		if err := s.db.Where("username = ? AND id != ?", *req.Username, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("username already taken")
		}
		updates["username"] = *req.Username
	}

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}

	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Perform update
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Fetch updated user
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return &user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if err := s.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return errors.New("invalid current password")
	}

	// Update password (will be hashed by BeforeUpdate hook if implemented)
	user.Password = newPassword
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

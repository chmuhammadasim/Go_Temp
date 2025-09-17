package handlers

import (
	"net/http"
	"strconv"

	"go-backend/internal/models"
	"go-backend/internal/services"
	"go-backend/internal/utils"
	"go-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
	logger      *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest

	// Bind and validate request
	if errors := utils.BindAndValidate(c, &req); len(errors) > 0 {
		h.logger.WithFields(logrus.Fields{
			"errors": errors,
			"ip":     c.ClientIP(),
		}).Warn("User registration validation failed")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Validation failed",
			"errors": errors,
		})
		return
	}

	// Register user
	response, err := h.userService.Register(&req)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
			"ip":    c.ClientIP(),
		}).Error("User registration failed")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": response.User.ID,
		"email":   response.User.Email,
		"ip":      c.ClientIP(),
	}).Info("User registered successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

// Login handles user authentication
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind and validate request
	if errors := utils.BindAndValidate(c, &req); len(errors) > 0 {
		h.logger.WithFields(logrus.Fields{
			"errors": errors,
			"ip":     c.ClientIP(),
		}).Warn("User login validation failed")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Validation failed",
			"errors": errors,
		})
		return
	}

	// Authenticate user
	response, err := h.userService.Login(&req)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
			"ip":    c.ClientIP(),
		}).Warn("User login failed")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": response.User.ID,
		"email":   response.User.Email,
		"ip":      c.ClientIP(),
	}).Info("User logged in successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// GetProfile gets the current user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID type",
		})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user profile")
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.ToResponse(),
	})
}

// GetUser gets a user by ID (admin only)
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user")
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.ToResponse(),
	})
}

// GetUsers gets all users with pagination (admin only)
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	// Convert to response format
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": userResponses,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// UpdateUser updates a user (admin or owner)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req models.UserUpdateRequest

	// Bind and validate request
	if errors := utils.BindAndValidate(c, &req); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Validation failed",
			"errors": errors,
		})
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update user")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":    user.ID,
		"updated_by": c.GetUint("user_id"),
	}).Info("User updated successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user.ToResponse(),
	})
}

// DeleteUser deletes a user (admin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Prevent self-deletion
	currentUserID := c.GetUint("user_id")
	if currentUserID == uint(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You cannot delete your own account",
		})
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		h.logger.WithError(err).Error("Failed to delete user")
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":    id,
		"deleted_by": currentUserID,
	}).Info("User deleted successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ChangePassword changes the current user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	// Bind and validate request
	if errors := utils.BindAndValidate(c, &req); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Validation failed",
			"errors": errors,
		})
		return
	}

	err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		h.logger.WithError(err).Error("Failed to change password")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithField("user_id", userID).Info("Password changed successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// Helper function for owner check middleware
func (h *UserHandler) GetUserIDFromParam(c *gin.Context) uint {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

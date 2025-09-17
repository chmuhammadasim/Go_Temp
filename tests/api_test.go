package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go-backend/internal/config"
	"go-backend/internal/database"
	"go-backend/internal/handlers"
	"go-backend/internal/middleware"
	"go-backend/internal/models"
	"go-backend/internal/services"
	"go-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// TestSuite provides a test suite for the application
type TestSuite struct {
	suite.Suite
	db          *gorm.DB
	router      *gin.Engine
	config      *config.Config
	logger      *logger.Logger
	userService *services.UserService
	authService *services.AuthService
	testUser    *models.User
	adminUser   *models.User
	authToken   string
	adminToken  string
}

// SetupSuite runs before all tests in the suite
func (suite *TestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test configuration
	suite.config = &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite",
			DSN:    ":memory:",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret-key-for-testing-only",
			Expiration: 24,
		},
		App: config.AppConfig{
			Environment: "test",
			Port:        "8080",
		},
	}

	// Setup logger
	suite.logger = logger.New(suite.config.App.Environment)

	// Setup database
	db, err := database.NewConnection(&suite.config.Database, suite.logger)
	suite.Require().NoError(err)
	suite.db = db

	// Run migrations
	err = database.RunMigrations(suite.db, suite.logger)
	suite.Require().NoError(err)

	// Setup services
	auditService := services.NewAuditService(suite.db, suite.logger)
	suite.userService = services.NewUserService(suite.db, suite.logger, auditService)
	suite.authService = services.NewAuthService(suite.db, suite.config, suite.logger, auditService)

	// Setup router
	suite.setupRouter()

	// Create test users
	suite.createTestUsers()
}

// TearDownSuite runs after all tests in the suite
func (suite *TestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

// SetupTest runs before each test
func (suite *TestSuite) SetupTest() {
	// Clean up data before each test if needed
}

// TearDownTest runs after each test
func (suite *TestSuite) TearDownTest() {
	// Clean up data after each test if needed
}

// setupRouter sets up the Gin router for testing
func (suite *TestSuite) setupRouter() {
	suite.router = gin.New()

	// Add middleware
	suite.router.Use(middleware.Logger(suite.logger))
	suite.router.Use(middleware.Recovery())
	suite.router.Use(middleware.CORS())

	// Setup handlers
	authHandler := handlers.NewAuthHandler(suite.authService, suite.logger)
	userHandler := handlers.NewUserHandler(suite.userService, suite.logger)

	// API routes
	v1 := suite.router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(suite.authService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", userHandler.ListUsers)
				admin.PUT("/users/:id", userHandler.UpdateUser)
				admin.DELETE("/users/:id", userHandler.DeleteUser)
			}
		}
	}
}

// createTestUsers creates test users for testing
func (suite *TestSuite) createTestUsers() {
	// Create regular test user
	testUser := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
	}

	hashedPassword, err := suite.authService.HashPassword("password123")
	suite.Require().NoError(err)
	testUser.Password = hashedPassword

	err = suite.db.Create(testUser).Error
	suite.Require().NoError(err)
	suite.testUser = testUser

	// Create admin test user
	adminUser := &models.User{
		Username:  "admin",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Role:      "admin",
		IsActive:  true,
	}

	adminUser.Password = hashedPassword
	err = suite.db.Create(adminUser).Error
	suite.Require().NoError(err)
	suite.adminUser = adminUser

	// Generate tokens
	suite.authToken, err = suite.authService.GenerateToken(suite.testUser)
	suite.Require().NoError(err)

	suite.adminToken, err = suite.authService.GenerateToken(suite.adminUser)
	suite.Require().NoError(err)
}

// Helper methods for making HTTP requests

func (suite *TestSuite) makeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(jsonBody)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *TestSuite) makeJSONRequest(method, path string, body interface{}, token string) (*httptest.ResponseRecorder, map[string]interface{}) {
	w := suite.makeRequest(method, path, body, token)

	var response map[string]interface{}
	if w.Body.Len() > 0 {
		err := json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)
	}

	return w, response
}

// Authentication Tests

func (suite *TestSuite) TestRegister() {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid registration",
			payload: map[string]interface{}{
				"username":  "newuser",
				"email":     "newuser@example.com",
				"password":  "password123",
				"firstName": "New",
				"lastName":  "User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"email": "incomplete@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Duplicate email",
			payload: map[string]interface{}{
				"username":  "duplicate",
				"email":     "test@example.com", // Already exists
				"password":  "password123",
				"firstName": "Duplicate",
				"lastName":  "User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"username":  "invaliduser",
				"email":     "invalid-email",
				"password":  "password123",
				"firstName": "Invalid",
				"lastName":  "User",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w, response := suite.makeJSONRequest("POST", "/api/v1/auth/register", tt.payload, "")
			suite.Equal(tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				suite.Contains(response, "message")
			} else {
				suite.Contains(response, "error")
			}
		})
	}
}

func (suite *TestSuite) TestLogin() {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid login",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid email",
			payload: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Missing fields",
			payload: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w, response := suite.makeJSONRequest("POST", "/api/v1/auth/login", tt.payload, "")
			suite.Equal(tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				suite.Contains(response, "token")
				suite.Contains(response, "user")
				suite.Contains(response, "expiresAt")
			} else {
				suite.Contains(response, "error")
			}
		})
	}
}

// User Tests

func (suite *TestSuite) TestGetProfile() {
	// Test with valid token
	w, response := suite.makeJSONRequest("GET", "/api/v1/users/profile", nil, suite.authToken)
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(response, "id")
	suite.Contains(response, "username")
	suite.Contains(response, "email")
	suite.Equal("testuser", response["username"])

	// Test without token
	w, response = suite.makeJSONRequest("GET", "/api/v1/users/profile", nil, "")
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(response, "error")

	// Test with invalid token
	w, response = suite.makeJSONRequest("GET", "/api/v1/users/profile", nil, "invalid-token")
	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(response, "error")
}

func (suite *TestSuite) TestUpdateProfile() {
	payload := map[string]interface{}{
		"firstName": "Updated",
		"lastName":  "Name",
	}

	w, response := suite.makeJSONRequest("PUT", "/api/v1/users/profile", payload, suite.authToken)
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(response, "user")

	// Verify the update
	w, response = suite.makeJSONRequest("GET", "/api/v1/users/profile", nil, suite.authToken)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("Updated", response["firstName"])
	suite.Equal("Name", response["lastName"])
}

// Admin Tests

func (suite *TestSuite) TestListUsers_AdminAccess() {
	w, response := suite.makeJSONRequest("GET", "/api/v1/admin/users", nil, suite.adminToken)
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(response, "users")
	suite.Contains(response, "pagination")
}

func (suite *TestSuite) TestListUsers_UserAccess() {
	w, response := suite.makeJSONRequest("GET", "/api/v1/admin/users", nil, suite.authToken)
	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(response, "error")
}

func (suite *TestSuite) TestUpdateUser_AdminAccess() {
	payload := map[string]interface{}{
		"firstName": "AdminUpdated",
	}

	path := fmt.Sprintf("/api/v1/admin/users/%d", suite.testUser.ID)
	w, response := suite.makeJSONRequest("PUT", path, payload, suite.adminToken)
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(response, "user")
}

func (suite *TestSuite) TestDeleteUser_AdminAccess() {
	// Create a user to delete
	userToDelete := &models.User{
		Username:  "todelete",
		Email:     "delete@example.com",
		FirstName: "To",
		LastName:  "Delete",
		Role:      "user",
		IsActive:  true,
		Password:  "hashedpassword",
	}
	err := suite.db.Create(userToDelete).Error
	suite.Require().NoError(err)

	path := fmt.Sprintf("/api/v1/admin/users/%d", userToDelete.ID)
	w, response := suite.makeJSONRequest("DELETE", path, nil, suite.adminToken)
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(response, "message")

	// Verify user is deleted
	var count int64
	suite.db.Model(&models.User{}).Where("id = ?", userToDelete.ID).Count(&count)
	suite.Equal(int64(0), count)
}

// Benchmark Tests

func (suite *TestSuite) BenchmarkLogin() {
	payload := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	suite.T().Run("BenchmarkLogin", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			w := suite.makeRequest("POST", "/api/v1/auth/login", payload, "")
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}

func (suite *TestSuite) BenchmarkGetProfile() {
	suite.T().Run("BenchmarkGetProfile", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			w := suite.makeRequest("GET", "/api/v1/users/profile", nil, suite.authToken)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}

// Integration Tests

func (suite *TestSuite) TestCompleteUserFlow() {
	// 1. Register a new user
	registerPayload := map[string]interface{}{
		"username":  "flowtest",
		"email":     "flowtest@example.com",
		"password":  "password123",
		"firstName": "Flow",
		"lastName":  "Test",
	}

	w, response := suite.makeJSONRequest("POST", "/api/v1/auth/register", registerPayload, "")
	suite.Equal(http.StatusCreated, w.Code)

	// 2. Login with the new user
	loginPayload := map[string]interface{}{
		"email":    "flowtest@example.com",
		"password": "password123",
	}

	w, response = suite.makeJSONRequest("POST", "/api/v1/auth/login", loginPayload, "")
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(response, "token")

	token := response["token"].(string)

	// 3. Get profile
	w, response = suite.makeJSONRequest("GET", "/api/v1/users/profile", nil, token)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("flowtest", response["username"])

	// 4. Update profile
	updatePayload := map[string]interface{}{
		"firstName": "Updated Flow",
	}

	w, response = suite.makeJSONRequest("PUT", "/api/v1/users/profile", updatePayload, token)
	suite.Equal(http.StatusOK, w.Code)

	// 5. Verify update
	w, response = suite.makeJSONRequest("GET", "/api/v1/users/profile", nil, token)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("Updated Flow", response["firstName"])
}

// Test Runner

func TestAPISuite(t *testing.T) {
	// Skip if running in CI without database
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(TestSuite))
}

// Additional Test Utilities

// TestHelper provides utility functions for testing
type TestHelper struct {
	suite *TestSuite
}

// NewTestHelper creates a new test helper
func NewTestHelper(suite *TestSuite) *TestHelper {
	return &TestHelper{suite: suite}
}

// CreateTestUser creates a test user and returns it
func (th *TestHelper) CreateTestUser(username, email string) *models.User {
	user := &models.User{
		Username:  username,
		Email:     email,
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
		Password:  "hashedpassword",
	}

	err := th.suite.db.Create(user).Error
	th.suite.Require().NoError(err)
	return user
}

// GenerateTestToken generates a JWT token for a user
func (th *TestHelper) GenerateTestToken(user *models.User) string {
	token, err := th.suite.authService.GenerateToken(user)
	th.suite.Require().NoError(err)
	return token
}

// CleanupTestData removes test data from the database
func (th *TestHelper) CleanupTestData() {
	// Clean up in reverse order of dependencies
	th.suite.db.Exec("DELETE FROM audit_logs")
	th.suite.db.Exec("DELETE FROM posts")
	th.suite.db.Exec("DELETE FROM users WHERE email LIKE '%test%'")
}

package docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &OpenAPIInfo{
	Version:     "1.0.0",
	Title:       "Go Backend API",
	Description: "A robust and scalable Go backend with authentication, role-based access control, and comprehensive features",
	Contact: Contact{
		Name:  "API Support",
		Email: "support@example.com",
		URL:   "https://example.com/support",
	},
	License: License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	},
}

// OpenAPIInfo represents basic API information
type OpenAPIInfo struct {
	Version     string  `json:"version"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Contact     Contact `json:"contact"`
	License     License `json:"license"`
}

// Contact represents contact information
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// License represents license information
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Server represents server information
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// SecurityScheme represents security scheme
type SecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Description  string `json:"description,omitempty"`
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      Schema      `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

// RequestBody represents request body
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType represents media type
type MediaType struct {
	Schema   Schema                 `json:"schema"`
	Example  interface{}            `json:"example,omitempty"`
	Examples map[string]interface{} `json:"examples,omitempty"`
}

// Response represents an API response
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Headers     map[string]Header    `json:"headers,omitempty"`
}

// Header represents response header
type Header struct {
	Description string `json:"description,omitempty"`
	Schema      Schema `json:"schema"`
}

// Schema represents JSON schema
type Schema struct {
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Description string            `json:"description,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Items       *Schema           `json:"items,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Enum        []interface{}     `json:"enum,omitempty"`
	Ref         string            `json:"$ref,omitempty"`
	AllOf       []Schema          `json:"allOf,omitempty"`
	AnyOf       []Schema          `json:"anyOf,omitempty"`
	OneOf       []Schema          `json:"oneOf,omitempty"`
	Minimum     *float64          `json:"minimum,omitempty"`
	Maximum     *float64          `json:"maximum,omitempty"`
	MinLength   *int              `json:"minLength,omitempty"`
	MaxLength   *int              `json:"maxLength,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
}

// PathItem represents path item
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Trace   *Operation `json:"trace,omitempty"`
}

// OpenAPISpec represents the complete OpenAPI specification
type OpenAPISpec struct {
	OpenAPI      string                 `json:"openapi"`
	Info         OpenAPIInfo            `json:"info"`
	Servers      []Server               `json:"servers,omitempty"`
	Paths        map[string]PathItem    `json:"paths"`
	Components   Components             `json:"components,omitempty"`
	Security     []map[string][]string  `json:"security,omitempty"`
	Tags         []Tag                  `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// Components represents reusable components
type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// Tag represents API tag
type Tag struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// ExternalDocumentation represents external documentation
type ExternalDocumentation struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// APIDocumentation manages API documentation
type APIDocumentation struct {
	spec     *OpenAPISpec
	examples map[string]interface{}
}

// NewAPIDocumentation creates a new API documentation instance
func NewAPIDocumentation() *APIDocumentation {
	return &APIDocumentation{
		spec: &OpenAPISpec{
			OpenAPI: "3.0.3",
			Info:    *SwaggerInfo,
			Servers: []Server{
				{
					URL:         "http://localhost:8080",
					Description: "Development server",
				},
				{
					URL:         "https://api.example.com",
					Description: "Production server",
				},
			},
			Paths: make(map[string]PathItem),
			Components: Components{
				Schemas:         make(map[string]Schema),
				SecuritySchemes: make(map[string]SecurityScheme),
			},
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
			Tags: []Tag{
				{Name: "auth", Description: "Authentication operations"},
				{Name: "users", Description: "User management operations"},
				{Name: "posts", Description: "Post management operations"},
				{Name: "admin", Description: "Administrative operations"},
				{Name: "files", Description: "File management operations"},
				{Name: "notifications", Description: "Notification operations"},
			},
		},
		examples: make(map[string]interface{}),
	}
}

// SetupSwagger initializes the Swagger documentation
func (api *APIDocumentation) SetupSwagger() {
	// Setup security schemes
	api.spec.Components.SecuritySchemes["bearerAuth"] = SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT Authorization header using the Bearer scheme",
	}

	// Setup common schemas
	api.setupCommonSchemas()
	api.setupAuthEndpoints()
	api.setupUserEndpoints()
	api.setupPostEndpoints()
	api.setupFileEndpoints()
	api.setupNotificationEndpoints()
	api.setupAdminEndpoints()
}

// setupCommonSchemas sets up common schema definitions
func (api *APIDocumentation) setupCommonSchemas() {
	// Error response schema
	api.spec.Components.Schemas["Error"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"error":   {Type: "string", Description: "Error type"},
			"message": {Type: "string", Description: "Error message"},
			"details": {Type: "object", Description: "Additional error details"},
		},
		Required: []string{"error", "message"},
	}

	// Success response schema
	api.spec.Components.Schemas["Success"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"message": {Type: "string", Description: "Success message"},
			"data":    {Type: "object", Description: "Response data"},
		},
		Required: []string{"message"},
	}

	// Pagination schema
	api.spec.Components.Schemas["Pagination"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"page":       {Type: "integer", Description: "Current page number"},
			"limit":      {Type: "integer", Description: "Items per page"},
			"total":      {Type: "integer", Description: "Total number of items"},
			"totalPages": {Type: "integer", Description: "Total number of pages"},
		},
	}

	// User schema
	api.spec.Components.Schemas["User"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"id":         {Type: "integer", Description: "User ID"},
			"username":   {Type: "string", Description: "Username"},
			"email":      {Type: "string", Format: "email", Description: "Email address"},
			"firstName":  {Type: "string", Description: "First name"},
			"lastName":   {Type: "string", Description: "Last name"},
			"role":       {Type: "string", Enum: []interface{}{"admin", "moderator", "user"}, Description: "User role"},
			"isActive":   {Type: "boolean", Description: "Account status"},
			"isVerified": {Type: "boolean", Description: "Email verification status"},
			"createdAt":  {Type: "string", Format: "date-time", Description: "Creation timestamp"},
			"updatedAt":  {Type: "string", Format: "date-time", Description: "Last update timestamp"},
		},
		Required: []string{"id", "username", "email", "role"},
	}

	// Login request schema
	api.spec.Components.Schemas["LoginRequest"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"email":    {Type: "string", Format: "email", Description: "Email address"},
			"password": {Type: "string", MinLength: &[]int{6}[0], Description: "Password"},
		},
		Required: []string{"email", "password"},
		Example: map[string]interface{}{
			"email":    "user@example.com",
			"password": "password123",
		},
	}

	// Register request schema
	api.spec.Components.Schemas["RegisterRequest"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"username":  {Type: "string", MinLength: &[]int{3}[0], Description: "Username"},
			"email":     {Type: "string", Format: "email", Description: "Email address"},
			"password":  {Type: "string", MinLength: &[]int{6}[0], Description: "Password"},
			"firstName": {Type: "string", Description: "First name"},
			"lastName":  {Type: "string", Description: "Last name"},
		},
		Required: []string{"username", "email", "password"},
		Example: map[string]interface{}{
			"username":  "newuser",
			"email":     "newuser@example.com",
			"password":  "password123",
			"firstName": "John",
			"lastName":  "Doe",
		},
	}

	// Token response schema
	api.spec.Components.Schemas["TokenResponse"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"token":     {Type: "string", Description: "JWT access token"},
			"expiresAt": {Type: "string", Format: "date-time", Description: "Token expiration time"},
			"user":      {Ref: "#/components/schemas/User"},
		},
		Required: []string{"token", "expiresAt", "user"},
	}
}

// setupAuthEndpoints sets up authentication endpoints documentation
func (api *APIDocumentation) setupAuthEndpoints() {
	// Login endpoint
	api.spec.Paths["/api/v1/auth/login"] = PathItem{
		Post: &Operation{
			Tags:        []string{"auth"},
			Summary:     "User login",
			Description: "Authenticate user and return JWT token",
			RequestBody: &RequestBody{
				Required:    true,
				Description: "Login credentials",
				Content: map[string]MediaType{
					"application/json": {
						Schema: Schema{Ref: "#/components/schemas/LoginRequest"},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Login successful",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/TokenResponse"},
						},
					},
				},
				"400": {
					Description: "Invalid request",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/Error"},
						},
					},
				},
				"401": {
					Description: "Invalid credentials",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/Error"},
						},
					},
				},
			},
			Security: []map[string][]string{}, // No auth required for login
		},
	}

	// Register endpoint
	api.spec.Paths["/api/v1/auth/register"] = PathItem{
		Post: &Operation{
			Tags:        []string{"auth"},
			Summary:     "User registration",
			Description: "Register a new user account",
			RequestBody: &RequestBody{
				Required:    true,
				Description: "Registration details",
				Content: map[string]MediaType{
					"application/json": {
						Schema: Schema{Ref: "#/components/schemas/RegisterRequest"},
					},
				},
			},
			Responses: map[string]Response{
				"201": {
					Description: "Registration successful",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/Success"},
						},
					},
				},
				"400": {
					Description: "Invalid request or user already exists",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/Error"},
						},
					},
				},
			},
			Security: []map[string][]string{}, // No auth required for registration
		},
	}
}

// setupUserEndpoints sets up user management endpoints documentation
func (api *APIDocumentation) setupUserEndpoints() {
	// Get current user profile
	api.spec.Paths["/api/v1/users/profile"] = PathItem{
		Get: &Operation{
			Tags:        []string{"users"},
			Summary:     "Get current user profile",
			Description: "Retrieve the profile of the authenticated user",
			Responses: map[string]Response{
				"200": {
					Description: "User profile retrieved successfully",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/User"},
						},
					},
				},
				"401": {
					Description: "Unauthorized",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/Error"},
						},
					},
				},
			},
		},
	}

	// List users (admin only)
	api.spec.Paths["/api/v1/users"] = PathItem{
		Get: &Operation{
			Tags:        []string{"users"},
			Summary:     "List users",
			Description: "Retrieve a paginated list of users (admin only)",
			Parameters: []Parameter{
				{
					Name:        "page",
					In:          "query",
					Description: "Page number",
					Schema:      Schema{Type: "integer", Minimum: &[]float64{1}[0]},
					Example:     1,
				},
				{
					Name:        "limit",
					In:          "query",
					Description: "Items per page",
					Schema:      Schema{Type: "integer", Minimum: &[]float64{1}[0], Maximum: &[]float64{100}[0]},
					Example:     10,
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Users retrieved successfully",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{
								Type: "object",
								Properties: map[string]Schema{
									"users":      {Type: "array", Items: &Schema{Ref: "#/components/schemas/User"}},
									"pagination": {Ref: "#/components/schemas/Pagination"},
								},
							},
						},
					},
				},
				"403": {
					Description: "Forbidden - Admin access required",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{Ref: "#/components/schemas/Error"},
						},
					},
				},
			},
		},
	}
}

// setupPostEndpoints sets up post management endpoints documentation
func (api *APIDocumentation) setupPostEndpoints() {
	// Post schema
	api.spec.Components.Schemas["Post"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"id":        {Type: "integer", Description: "Post ID"},
			"title":     {Type: "string", Description: "Post title"},
			"content":   {Type: "string", Description: "Post content"},
			"authorId":  {Type: "integer", Description: "Author user ID"},
			"author":    {Ref: "#/components/schemas/User"},
			"createdAt": {Type: "string", Format: "date-time", Description: "Creation timestamp"},
			"updatedAt": {Type: "string", Format: "date-time", Description: "Last update timestamp"},
		},
	}

	// List posts
	api.spec.Paths["/api/v1/posts"] = PathItem{
		Get: &Operation{
			Tags:        []string{"posts"},
			Summary:     "List posts",
			Description: "Retrieve a paginated list of posts",
			Parameters: []Parameter{
				{
					Name:        "page",
					In:          "query",
					Description: "Page number",
					Schema:      Schema{Type: "integer", Minimum: &[]float64{1}[0]},
					Example:     1,
				},
				{
					Name:        "limit",
					In:          "query",
					Description: "Items per page",
					Schema:      Schema{Type: "integer", Minimum: &[]float64{1}[0], Maximum: &[]float64{100}[0]},
					Example:     10,
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Posts retrieved successfully",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{
								Type: "object",
								Properties: map[string]Schema{
									"posts":      {Type: "array", Items: &Schema{Ref: "#/components/schemas/Post"}},
									"pagination": {Ref: "#/components/schemas/Pagination"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// setupFileEndpoints sets up file management endpoints documentation
func (api *APIDocumentation) setupFileEndpoints() {
	// File upload
	api.spec.Paths["/api/v1/files/upload"] = PathItem{
		Post: &Operation{
			Tags:        []string{"files"},
			Summary:     "Upload file",
			Description: "Upload a file to the server",
			RequestBody: &RequestBody{
				Required:    true,
				Description: "File to upload",
				Content: map[string]MediaType{
					"multipart/form-data": {
						Schema: Schema{
							Type: "object",
							Properties: map[string]Schema{
								"file": {
									Type:        "string",
									Format:      "binary",
									Description: "File to upload",
								},
								"category": {
									Type:        "string",
									Description: "File category",
									Example:     "document",
								},
							},
							Required: []string{"file"},
						},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "File uploaded successfully",
					Content: map[string]MediaType{
						"application/json": {
							Schema: Schema{
								Type: "object",
								Properties: map[string]Schema{
									"fileId":   {Type: "string", Description: "File ID"},
									"filename": {Type: "string", Description: "Original filename"},
									"url":      {Type: "string", Description: "File access URL"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// setupNotificationEndpoints sets up notification endpoints documentation
func (api *APIDocumentation) setupNotificationEndpoints() {
	// Get user notifications
	api.spec.Paths["/api/v1/notifications"] = PathItem{
		Get: &Operation{
			Tags:        []string{"notifications"},
			Summary:     "Get user notifications",
			Description: "Retrieve notifications for the authenticated user",
			Parameters: []Parameter{
				{
					Name:        "limit",
					In:          "query",
					Description: "Number of notifications to retrieve",
					Schema:      Schema{Type: "integer", Minimum: &[]float64{1}[0], Maximum: &[]float64{50}[0]},
					Example:     10,
				},
				{
					Name:        "offset",
					In:          "query",
					Description: "Number of notifications to skip",
					Schema:      Schema{Type: "integer", Minimum: &[]float64{0}[0]},
					Example:     0,
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Notifications retrieved successfully",
				},
			},
		},
	}
}

// setupAdminEndpoints sets up admin endpoints documentation
func (api *APIDocumentation) setupAdminEndpoints() {
	// Admin dashboard
	api.spec.Paths["/api/v1/admin/dashboard"] = PathItem{
		Get: &Operation{
			Tags:        []string{"admin"},
			Summary:     "Get admin dashboard data",
			Description: "Retrieve dashboard statistics and metrics (admin only)",
			Responses: map[string]Response{
				"200": {
					Description: "Dashboard data retrieved successfully",
				},
				"403": {
					Description: "Forbidden - Admin access required",
				},
			},
		},
	}
}

// GetOpenAPISpec returns the complete OpenAPI specification
func (api *APIDocumentation) GetOpenAPISpec() *OpenAPISpec {
	return api.spec
}

// ServeSwaggerUI serves the Swagger UI
func (api *APIDocumentation) ServeSwaggerUI() gin.HandlerFunc {
	return func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>API Documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/api/docs/openapi.json',
      dom_id: '#swagger-ui',
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.presets.standalone,
      ],
    });
  };
</script>
</body>
</html>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	}
}

// ServeOpenAPIJSON serves the OpenAPI specification as JSON
func (api *APIDocumentation) ServeOpenAPIJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, api.spec)
	}
}

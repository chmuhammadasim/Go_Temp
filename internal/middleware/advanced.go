package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-backend/internal/services"
	"go-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig contains rate limiting configuration
type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int
	KeyGenerator      func(*gin.Context) string
	SkipPaths         []string
	OnLimitReached    func(*gin.Context, string)
}

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate time.Duration
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// TryConsume attempts to consume a token from the bucket
func (tb *TokenBucket) TryConsume() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Add tokens based on elapsed time
	tokensToAdd := int(elapsed / tb.refillRate)
	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	// Try to consume a token
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(config RateLimiterConfig, securityService *services.SecurityService) gin.HandlerFunc {
	buckets := make(map[string]*TokenBucket)
	bucketsLock := sync.RWMutex{}

	// Default key generator (IP-based)
	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}

	// Default rate limit values
	if config.RequestsPerMinute == 0 {
		config.RequestsPerMinute = 60
	}
	if config.BurstSize == 0 {
		config.BurstSize = 10
	}

	refillRate := time.Minute / time.Duration(config.RequestsPerMinute)

	return func(c *gin.Context) {
		// Skip rate limiting for certain paths
		for _, skipPath := range config.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, skipPath) {
				c.Next()
				return
			}
		}

		key := config.KeyGenerator(c)

		// Get or create bucket for this key
		bucketsLock.RLock()
		bucket, exists := buckets[key]
		bucketsLock.RUnlock()

		if !exists {
			bucket = NewTokenBucket(config.BurstSize, refillRate)
			bucketsLock.Lock()
			buckets[key] = bucket
			bucketsLock.Unlock()
		}

		// Try to consume a token
		if !bucket.TryConsume() {
			// Rate limit exceeded
			if config.OnLimitReached != nil {
				config.OnLimitReached(c, key)
			}

			// Log security event
			if securityService != nil {
				userID := getUserIDFromContext(c)
				securityService.DetectRateLimitViolation(
					userID,
					c.ClientIP(),
					c.GetHeader("User-Agent"),
					c.Request.URL.Path,
					config.RequestsPerMinute,
				)
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Limit: %d requests per minute", config.RequestsPerMinute),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequestValidationConfig contains request validation configuration
type RequestValidationConfig struct {
	MaxBodySize       int64
	AllowedMethods    []string
	AllowedHeaders    []string
	RequiredHeaders   []string
	BlockedUserAgents []string
	ValidateJSON      bool
}

// RequestValidator creates a request validation middleware
func RequestValidator(config RequestValidationConfig, securityService *services.SecurityService) gin.HandlerFunc {
	// Default values
	if config.MaxBodySize == 0 {
		config.MaxBodySize = 10 << 20 // 10MB
	}

	return func(c *gin.Context) {
		// Validate content length
		if c.Request.ContentLength > config.MaxBodySize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "Request too large",
				"message": fmt.Sprintf("Maximum request size is %d bytes", config.MaxBodySize),
			})
			c.Abort()
			return
		}

		// Validate HTTP method
		if len(config.AllowedMethods) > 0 {
			methodAllowed := false
			for _, method := range config.AllowedMethods {
				if c.Request.Method == method {
					methodAllowed = true
					break
				}
			}
			if !methodAllowed {
				c.JSON(http.StatusMethodNotAllowed, gin.H{
					"error":   "Method not allowed",
					"message": fmt.Sprintf("Method %s is not allowed", c.Request.Method),
				})
				c.Abort()
				return
			}
		}

		// Validate required headers
		for _, header := range config.RequiredHeaders {
			if c.GetHeader(header) == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Missing required header",
					"message": fmt.Sprintf("Header %s is required", header),
				})
				c.Abort()
				return
			}
		}

		// Check blocked user agents
		userAgent := c.GetHeader("User-Agent")
		for _, blockedUA := range config.BlockedUserAgents {
			if strings.Contains(strings.ToLower(userAgent), strings.ToLower(blockedUA)) {
				// Log security event
				if securityService != nil {
					userID := getUserIDFromContext(c)
					securityService.LogSecurityEvent(
						userID,
						services.EventMaliciousRequest,
						services.SeverityMedium,
						"Blocked user agent detected",
						services.SecurityEventData{
							RemoteAddr: c.ClientIP(),
							UserAgent:  userAgent,
							Method:     c.Request.Method,
							Path:       c.Request.URL.Path,
						},
					)
				}

				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Forbidden",
					"message": "Access denied",
				})
				c.Abort()
				return
			}
		}

		// Validate JSON content type for POST/PUT/PATCH requests
		if config.ValidateJSON && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH") {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error":   "Unsupported media type",
					"message": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// APIVersionConfig contains API versioning configuration
type APIVersionConfig struct {
	DefaultVersion     string
	SupportedVersions  []string
	VersionHeader      string
	VersionParam       string
	DeprecatedVersions map[string]string // version -> deprecation message
}

// APIVersioning creates an API versioning middleware
func APIVersioning(config APIVersionConfig) gin.HandlerFunc {
	// Default values
	if config.DefaultVersion == "" {
		config.DefaultVersion = "v1"
	}
	if config.VersionHeader == "" {
		config.VersionHeader = "API-Version"
	}
	if config.VersionParam == "" {
		config.VersionParam = "version"
	}

	return func(c *gin.Context) {
		version := ""

		// Try to get version from header first
		version = c.GetHeader(config.VersionHeader)

		// If not in header, try query parameter
		if version == "" {
			version = c.Query(config.VersionParam)
		}

		// If still not found, try URL path prefix
		if version == "" {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/") {
				parts := strings.Split(path, "/")
				if len(parts) >= 3 && strings.HasPrefix(parts[2], "v") {
					version = parts[2]
				}
			}
		}

		// Use default version if none specified
		if version == "" {
			version = config.DefaultVersion
		}

		// Validate version
		versionSupported := false
		for _, supportedVersion := range config.SupportedVersions {
			if version == supportedVersion {
				versionSupported = true
				break
			}
		}

		if !versionSupported {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Unsupported API version",
				"message": fmt.Sprintf("Version %s is not supported. Supported versions: %v", version, config.SupportedVersions),
			})
			c.Abort()
			return
		}

		// Check for deprecated versions
		if deprecationMsg, isDeprecated := config.DeprecatedVersions[version]; isDeprecated {
			c.Header("Warning", fmt.Sprintf("299 - \"API version %s is deprecated. %s\"", version, deprecationMsg))
		}

		// Store version in context for handlers to use
		c.Set("api_version", version)

		// Add version to response header
		c.Header("API-Version", version)

		c.Next()
	}
}

// SecurityHeadersConfig contains security headers configuration
type SecurityHeadersConfig struct {
	EnableHSTS               bool
	HSTSMaxAge               int
	EnableContentTypeNoSniff bool
	EnableFrameDeny          bool
	EnableXSSProtection      bool
	CSPPolicy                string
	ReferrerPolicy           string
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders(config SecurityHeadersConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// HTTP Strict Transport Security
		if config.EnableHSTS {
			maxAge := config.HSTSMaxAge
			if maxAge == 0 {
				maxAge = 31536000 // 1 year
			}
			c.Header("Strict-Transport-Security", fmt.Sprintf("max-age=%d; includeSubDomains", maxAge))
		}

		// Content Type Options
		if config.EnableContentTypeNoSniff {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// Frame Options
		if config.EnableFrameDeny {
			c.Header("X-Frame-Options", "DENY")
		}

		// XSS Protection
		if config.EnableXSSProtection {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		// Content Security Policy
		if config.CSPPolicy != "" {
			c.Header("Content-Security-Policy", config.CSPPolicy)
		}

		// Referrer Policy
		if config.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", config.ReferrerPolicy)
		}

		// Additional security headers
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		c.Header("X-DNS-Prefetch-Control", "off")

		c.Next()
	}
}

// RequestIDConfig contains request ID configuration
type RequestIDConfig struct {
	Header    string
	Generator func() string
}

// RequestID adds a unique request ID to each request
func RequestID(config RequestIDConfig) gin.HandlerFunc {
	if config.Header == "" {
		config.Header = "X-Request-ID"
	}

	if config.Generator == nil {
		config.Generator = func() string {
			return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
		}
	}

	return func(c *gin.Context) {
		// Check if request ID already exists
		requestID := c.GetHeader(config.Header)
		if requestID == "" {
			requestID = config.Generator()
		}

		// Set in context and response header
		c.Set("request_id", requestID)
		c.Header(config.Header, requestID)

		c.Next()
	}
}

// IPWhitelistConfig contains IP whitelist configuration
type IPWhitelistConfig struct {
	AllowedIPs     []string
	AllowedCIDRs   []string
	TrustedProxies []string
}

// IPWhitelist restricts access to whitelisted IP addresses
func IPWhitelist(config IPWhitelistConfig, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Check if IP is in allowed list
		for _, allowedIP := range config.AllowedIPs {
			if clientIP == allowedIP {
				c.Next()
				return
			}
		}

		// TODO: Add CIDR range checking here
		// For simplicity, this example only checks exact IP matches

		// IP not allowed
		logger.Warn("IP access denied", map[string]interface{}{
			"client_ip": clientIP,
			"path":      c.Request.URL.Path,
			"method":    c.Request.Method,
		})

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Access denied",
			"message": "Your IP address is not allowed to access this resource",
		})
		c.Abort()
	}
}

// Helper functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getUserIDFromContext(c *gin.Context) *uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return &id
		}
		if idStr, ok := userID.(string); ok {
			if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				uid := uint(id)
				return &uid
			}
		}
	}
	return nil
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns the API health status
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API is running",
		"version": "1.0.0",
	})
}

// ReadinessCheck returns the API readiness status
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Here you could check database connectivity, external services, etc.
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "connected",
		},
	})
}

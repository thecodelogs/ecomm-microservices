package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// GET /health
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "api-gateway",
		"timestamp": time.Now().UTC(),
	})
}

// GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	// TODO: check downstream services
	c.JSON(http.StatusOK, gin.H{"ready": true})
}

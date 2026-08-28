package routes

import (
	"ratelimiter/internal/health"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoute(r *gin.Engine, h *health.Handler) {
	r.GET("/health", h.Check)
}

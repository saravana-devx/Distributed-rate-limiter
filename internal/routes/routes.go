package routes

import (
	"ratelimiter/internal/check"
	"ratelimiter/internal/client"
	"ratelimiter/internal/health"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, healthHandler *health.Handler, clientHandler *client.Handler, clientService *client.Service, checkHandler *check.Handler) {
	RegisterHealthRoute(r, healthHandler)

	v1 := r.Group("/api/v1")
	RegisterClientRoute(v1, clientHandler)
	RegisterCheckRoute(v1, checkHandler, clientService)
}

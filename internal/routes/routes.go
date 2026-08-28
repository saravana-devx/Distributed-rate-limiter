package routes

import (
	"ratelimiter/internal/client"
	"ratelimiter/internal/health"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, healthHandler *health.Handler, clientHandler *client.Handler, clientService *client.Service) {
	RegisterHealthRoute(r, healthHandler)
	RegisterClientRoute(r, clientHandler, clientService)
}

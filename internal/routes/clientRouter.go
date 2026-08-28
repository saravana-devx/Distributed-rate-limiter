package routes

import (
	"ratelimiter/internal/client"
	"ratelimiter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoute(r *gin.Engine, h *client.Handler, svc *client.Service) {
	clients := r.Group("/clients")
	{
		clients.POST("", h.CreateClient)
		clients.GET("", middleware.APIKeyAuth(svc), h.GetAllClients)
		clients.GET("/:id", middleware.APIKeyAuth(svc), h.GetClientByID)
		clients.PUT("/:id", middleware.APIKeyAuth(svc), h.UpdateClient)
		clients.DELETE("/:id", middleware.APIKeyAuth(svc), h.DeleteClient)
	}
}

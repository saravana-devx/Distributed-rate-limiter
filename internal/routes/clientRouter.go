package routes

import (
	"ratelimiter/internal/client"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoute(r gin.IRouter, h *client.Handler) {
	clients := r.Group("/clients")
	{
		clients.POST("", h.CreateClient)
		clients.GET("", h.GetAllClients)
		clients.GET("/:id", h.GetClientByID)
		clients.PUT("/:id", h.UpdateClient)
		clients.DELETE("/:id", h.DeleteClient)
	}
}

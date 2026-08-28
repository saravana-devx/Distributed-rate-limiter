package routes

import (
	"ratelimiter/internal/check"
	"ratelimiter/internal/client"
	"ratelimiter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCheckRoute(r gin.IRouter, h *check.Handler, svc *client.Service) {
	check := r.Group("/check")
	{
		check.POST("", middleware.APIKeyAuth(svc), h.Check)
	}
}

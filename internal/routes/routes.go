package routes

import (
	"ratelimiter/internal/check"
	"ratelimiter/internal/client"
	"ratelimiter/internal/health"
	"ratelimiter/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Register(r *gin.Engine, healthHandler *health.Handler, clientHandler *client.Handler, clientService *client.Service, checkHandler *check.Handler) {
	r.Use(middleware.PrometheusMiddleware())

	RegisterHealthRoute(r, healthHandler)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("/api/v1")
	RegisterClientRoute(v1, clientHandler)
	RegisterCheckRoute(v1, checkHandler, clientService)
}

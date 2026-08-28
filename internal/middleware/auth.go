package middleware

import (
	"context"
	"net/http"
	"time"

	"ratelimiter/internal/client"
	"ratelimiter/internal/httpx"

	"github.com/gin-gonic/gin"
)

const ClientContextKey = "client"

func APIKeyAuth(svc *client.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			httpx.AbortError(c, http.StatusUnauthorized, httpx.MsgUnauthorized)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		cl, err := svc.GetClientByAPIKeyService(ctx, apiKey)
		if err != nil {
			httpx.AbortError(c, http.StatusUnauthorized, httpx.MsgUnauthorized)
			return
		}
		c.Set(ClientContextKey, cl)
		c.Next()
	}
}

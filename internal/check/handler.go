package check

import (
	"context"
	"fmt"
	"net/http"
	"ratelimiter/internal/client"
	"ratelimiter/internal/httpx"
	"ratelimiter/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Check(c *gin.Context) {
	// Get the client data from the middleware
	cl := c.MustGet(middleware.ClientContextKey).(*client.Client)
	var body CheckRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.BindError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	data, err := h.svc.Check(ctx, cl, &body)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, httpx.MsgInternalError)
		return
	}
	// handle if allowed is false/true
	if !data.Allowed {
		c.Header("Retry-After", fmt.Sprintf("%d", data.ResetAt-time.Now().Unix()))
		httpx.Success(c, http.StatusTooManyRequests, "rate limit exceeded", data)
		return
	}
	httpx.Success(c, http.StatusOK, "request allowed", data)
}

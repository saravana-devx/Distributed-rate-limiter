package health

import (
	"context"
	"net/http"
	"time"

	"ratelimiter/internal/httpx"
	"ratelimiter/internal/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Redis
}

func NewHealthHandler(db *gorm.DB, rdb *redis.Redis) *Handler {
	return &Handler{db: db, rdb: rdb}
}

func (h *Handler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		httpx.Error(c, http.StatusServiceUnavailable, MsgUnhealthy)
		return
	}

	if err := h.rdb.Ping(ctx); err != nil {
		httpx.Error(c, http.StatusServiceUnavailable, MsgUnhealthy)
		return
	}

	httpx.Success(c, http.StatusOK, MsgHealthy, nil)
}

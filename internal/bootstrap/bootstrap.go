package bootstrap

import (
	"context"
	"fmt"
	"log"

	"ratelimiter/internal/config"
	"ratelimiter/internal/database"
	"ratelimiter/internal/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	router *gin.Engine
	DB     *gorm.DB
	Redis  *redis.Redis
}

func New(cfg *config.Config) (*App, error) {
	db, err := database.ConnectDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	rdb := redis.NewRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err := rdb.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %w", err)
	}

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	return &App{router: router, DB: db, Redis: rdb}, nil
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func (a *App) Stop() {
	log.Println("bootstrap stop...")
}

package bootstrap

import (
	"context"
	"fmt"
	"log"

	"ratelimiter/internal/client"
	"ratelimiter/internal/config"
	"ratelimiter/internal/database"
	"ratelimiter/internal/health"
	"ratelimiter/internal/redis"
	"ratelimiter/internal/routes"

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

	clientRepo := client.NewRepository(db)
	clientService := client.NewService(clientRepo, rdb)
	clientHandler := client.NewClientHandler(clientService)
	healthHandler := health.NewHealthHandler(db, rdb)

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	routes.Register(router, healthHandler, clientHandler, clientService)

	return &App{router: router, DB: db, Redis: rdb}, nil
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func (a *App) Stop() {
	log.Println("bootstrap stop...")
}

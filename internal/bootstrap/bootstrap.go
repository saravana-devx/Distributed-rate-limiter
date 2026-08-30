package bootstrap

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"ratelimiter/internal/check"
	"ratelimiter/internal/client"
	"ratelimiter/internal/config"
	"ratelimiter/internal/database"
	"ratelimiter/internal/health"
	"ratelimiter/internal/metrics"
	"ratelimiter/internal/redis"
	"ratelimiter/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	metrics.Init()

	healthHandler := health.NewHealthHandler(db, rdb)

	clientRepo := client.NewRepository(db)
	clientService := client.NewService(clientRepo, rdb)
	clientHandler := client.NewClientHandler(clientService)

	checkService := check.NewService(rdb)
	checkHandler := check.NewHandler(checkService)

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
			if name == "-" {
				return ""
			}
			return name
		})
	}
	routes.Register(router, healthHandler, clientHandler, clientService, checkHandler)

	return &App{router: router, DB: db, Redis: rdb}, nil
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func (a *App) Stop() {
	log.Println("bootstrap stop...")
}

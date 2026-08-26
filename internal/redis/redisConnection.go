package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

type Redis struct {
	Client *redis.Client
}

var (
	redisInstance *Redis
	once          sync.Once
)

// singleton pattern guarantees that this redis connection will be created only once
// every goroutine / function will get the same pointer reference
// without singleton pattern we might get multiple redis connections
func NewRedis(redisAddr string, redisPassword string) *Redis {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		})
		redisInstance = &Redis{Client: rdb}
	})
	return redisInstance
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// script.Run()  →  sends Lua to Redis, returns *redis.Cmd
// .Result()     →  unpacks it into (interface{}, error)
func (r *Redis) RunScript(ctx context.Context, script *redis.Script, keys []string, args ...interface{}) (interface{}, error) {
	return script.Run(ctx, r.Client, keys, args...).Result()
}

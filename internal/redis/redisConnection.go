package redis

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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

// Set JSON-encodes value and stores it under key. ttl of 0 means no expiry.
func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Del removes key. Not an error if key doesn't exist.
func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// script.Run()  →  sends Lua to Redis, returns *redis.Cmd
// .Result()     →  unpacks it into (interface{}, error)
func (r *Redis) RunScript(ctx context.Context, script *redis.Script, keys []string, args ...interface{}) (interface{}, error) {
	return script.Run(ctx, r.Client, keys, args...).Result()
}

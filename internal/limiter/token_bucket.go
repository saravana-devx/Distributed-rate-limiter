package limiter

import (
	"context"
	_ "embed"

	"github.com/redis/go-redis/v9"
	redisClient "ratelimiter/internal/redis"
)

//go:embed lua/token_bucket.lua
var tokenBucketScriptSrc string

var tokenBucketScript = redis.NewScript(tokenBucketScriptSrc)

type TokenBucket struct {
	rdb *redisClient.Redis
}

func NewTokenBucket(rdb *redisClient.Redis) *TokenBucket {
	return &TokenBucket{rdb: rdb}
}

func (t *TokenBucket) Allow(ctx context.Context, key string, limit int, window int) (*Result, error) {
	res, err := t.rdb.RunScript(ctx, tokenBucketScript, []string{key}, limit, window)
	if err != nil {
		return nil, err
	}
	vals := res.([]interface{})
	return &Result{
		Allowed:   vals[0].(int64) == 1,
		Remaining: vals[1].(int64),
		ResetAt:   vals[2].(int64),
	}, nil
}

// compile-time check that TokenBucket implements RateLimiter
var _ RateLimiter = (*TokenBucket)(nil)

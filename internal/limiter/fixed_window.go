package limiter

import (
	"context"
	_ "embed"

	"github.com/redis/go-redis/v9"
	redisClient "ratelimiter/internal/redis"
)

// the below line is not a comment it load the lua-script at runtime
// into the below mention variable
//go:embed lua/fixed_window.lua
var fixedWindowScriptSrc string

var fixedWindowScript = redis.NewScript(fixedWindowScriptSrc)

// FixedWindow implements RateLimiter using Redis INCR/EXPIRE over a fixed
// window that resets at the start of every window rather than sliding.
type FixedWindow struct {
	rdb *redisClient.Redis
}

func NewFixedWindow(rdb *redisClient.Redis) *FixedWindow {
	return &FixedWindow{rdb: rdb}
}

func (f *FixedWindow) Allow(ctx context.Context, key string, limit int, window int) (*Result, error) {
	res, err := f.rdb.RunScript(ctx, fixedWindowScript, []string{key}, limit, window)
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

var _ RateLimiter = (*FixedWindow)(nil)

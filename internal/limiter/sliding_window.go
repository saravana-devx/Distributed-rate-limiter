package limiter

import (
	"context"
	_ "embed"

	"github.com/redis/go-redis/v9"
	redisClient "ratelimiter/internal/redis"
)

//go:embed lua/sliding_window.lua
var slidingWindowScriptSrc string

var slidingWindowScript = redis.NewScript(slidingWindowScriptSrc)

// SlidingWindow implements RateLimiter using a Redis sorted set to track
// individual request timestamps within a continuously sliding window.
type SlidingWindow struct {
	rdb *redisClient.Redis
}

func NewSlidingWindow(rdb *redisClient.Redis) *SlidingWindow {
	return &SlidingWindow{rdb: rdb}
}

func (s *SlidingWindow) Allow(ctx context.Context, key string, limit int, window int) (*Result, error) {
	res, err := s.rdb.RunScript(ctx, slidingWindowScript, []string{key}, limit, window)
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

var _ RateLimiter = (*SlidingWindow)(nil)

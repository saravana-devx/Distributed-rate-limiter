package check

import (
	"context"
	"fmt"
	"ratelimiter/internal/client"
	"ratelimiter/internal/limiter"
	"ratelimiter/internal/redis"
)

type Service struct {
	rdb        *redis.Redis
	strategies map[string]limiter.RateLimiter
}

func NewService(rdb *redis.Redis) *Service {
	return &Service{
		rdb: rdb,
		strategies: map[string]limiter.RateLimiter{
			"fixed_window":   limiter.NewFixedWindow(rdb),
			"sliding_window": limiter.NewSlidingWindow(rdb),
			"token_bucket":   limiter.NewTokenBucket(rdb),
		},
	}
}

func (s *Service) Check(ctx context.Context, cl *client.Client, req *CheckRequest) (*CheckResult, error) {
	// Build redis key
	key := fmt.Sprintf("limiter:%s:%s:%s", cl.Algorithm, cl.ClientID, req.Identifier)

	strategy, ok := s.strategies[string(cl.Algorithm)]
	if !ok {
		return nil, fmt.Errorf("unknown algorithm: %s", cl.Algorithm)
	}

	res, err := strategy.Allow(ctx, key, cl.Limit, cl.WindowSeconds)
	if err != nil {
		return nil, err
	}
	return &CheckResult{Allowed: res.Allowed, Remaining: int(res.Remaining), ResetAt: res.ResetAt}, nil
}

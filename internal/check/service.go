package check

import (
	"context"
	"fmt"
	"ratelimiter/internal/client"
	"ratelimiter/internal/limiter"
	"ratelimiter/internal/redis"
)

type CheckRepository interface {
}

type Service struct {
	rdb *redis.Redis
}

func NewService(rdb *redis.Redis) *Service {
	return &Service{rdb: rdb}
}

func (s *Service) Check(ctx context.Context, cl *client.Client, req *CheckRequest) (*CheckResult, error) {

	// Build redis key
	key := fmt.Sprintf("limiter:%s:%s:%s", cl.Algorithm, cl.ClientID, req.Identifier)

	// switch on Client Algorithm -> call right algorithm
	switch cl.Algorithm {
	case "token_bucket":
		res, err := limiter.TokenBucket(ctx, s.rdb, key, cl.Limit, cl.WindowSeconds)
		if err != nil {
			return nil, err
		}
		return &CheckResult{Allowed: res.Allowed, Remaining: int(res.Remaining)}, nil
	case "fixed_window":
		res, err := limiter.FixedWindow(ctx, s.rdb, key, cl.Limit, cl.WindowSeconds)
		if err != nil {
			return nil, err
		}
		return &CheckResult{Allowed: res.Allowed, Remaining: int(res.Remaining)}, nil
	case "sliding_window":
		res, err := limiter.SlidingWindow(ctx, s.rdb, key, cl.Limit, cl.WindowSeconds)
		if err != nil {
			return nil, err
		}
		return &CheckResult{Allowed: res.Allowed, Remaining: int(res.Remaining)}, nil

	}
	return nil, fmt.Errorf("unknown algorithm: %s", cl.Algorithm)
}

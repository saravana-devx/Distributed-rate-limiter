package limiter

import "context"

type Result struct {
	Allowed   bool
	Remaining int64
	ResetAt   int64
}

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window int) (*Result, error)
}

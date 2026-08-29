package test

import (
	"context"
	"testing"

	limiter "ratelimiter/internal/limiter"
	redisClient "ratelimiter/internal/redis"
)

// t.Errorf → logs failure but continues test
// t.Fatalf → logs failure and stops test immediately
func TestFixedWindow(t *testing.T) {
	rdb := redisClient.NewRedis("localhost:6379", "ratelimiter")
	fw := limiter.NewFixedWindow(rdb)
	ctx := context.Background()
	key := "test:fixed_window"

	// cleanup before and after
	rdb.Client.Del(ctx, key)
	defer rdb.Client.Del(ctx, key)

	// first request should be allowed
	res, err := fw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("FixedWindow returned an error: %v", err)
	}

	if !res.Allowed {
		t.Errorf("Expected Allowed, but got rejected")
	}

	if res.Remaining != 2 {
		t.Errorf("expected remaining 2, got %d", res.Remaining)
	}

	if res.ResetAt <= 0 {
		t.Errorf("expected positive reset_at, got %d", res.ResetAt)
	}

	// second request should be allowed
	res2, err := fw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res2.Allowed {
		t.Errorf("expected allowed, got rejected")
	}
	if res2.Remaining != 1 {
		t.Errorf("expected remaining 1, got %d", res2.Remaining)
	}

	// third request should be allowed
	res3, err := fw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res3.Allowed {
		t.Errorf("expected allowed, got rejected")
	}
	if res3.Remaining != 0 {
		t.Errorf("expected remaining 0, got %d", res3.Remaining)
	}

	// fourth request should be rejected — limit exceeded
	res4, err := fw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res4.Allowed {
		t.Errorf("expected rejected, got allowed")
	}
	if res4.Remaining != 0 {
		t.Errorf("expected remaining 0, got %d", res4.Remaining)
	}
}

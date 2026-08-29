package test

import (
	"context"
	"testing"

	limiter "ratelimiter/internal/limiter"
	redisClient "ratelimiter/internal/redis"
)

func TestSlidingWindow(t *testing.T) {
	rdb := redisClient.NewRedis("localhost:6379", "ratelimiter")
	sw := limiter.NewSlidingWindow(rdb)
	ctx := context.Background()
	key := "test:sliding_window"

	rdb.Client.Del(ctx, key)
	defer rdb.Client.Del(ctx, key)

	res, err := sw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("Sliding Window returned an error : %v", err)
	}

	if !res.Allowed {
		t.Errorf("Expected Allowed, but got rejected")
	}

	if res.Remaining != 2 {
		t.Errorf("excepted remaining 2, got %d", res.Remaining)
	}

	if res.ResetAt <= 0 {
		t.Errorf("expected positive reset_at, got %d", res.ResetAt)
	}

	res2, err := sw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("Unexcepted error : %v", err)
	}

	if !res2.Allowed {
		t.Errorf("excepted allowed, got rejected")
	}

	if res2.Remaining != 1 {
		t.Errorf("excepted remaining 1, got %d", res2.Remaining)
	}

	res3, err := sw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("Unexcepted error : %v", err)
	}

	if !res3.Allowed {
		t.Errorf("excepted allowed, got rejected")
	}

	if res3.Remaining != 0 {
		t.Errorf("excepted remaining 0, got %d", res3.Remaining)
	}

	res4, err := sw.Allow(ctx, key, 3, 60)
	if err != nil {
		t.Fatalf("Unexcepted error : %v", err)
	}

	if res4.Allowed {
		t.Errorf("excepted rejected, got allowed")
	}

	if res4.Remaining != 0 {
		t.Errorf("excepted remaining 0, got %d", res4.Remaining)
	}
}

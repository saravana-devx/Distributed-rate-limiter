package test

import (
	"context"
	"ratelimiter/internal/limiter"
	redisClient "ratelimiter/internal/redis"
	"testing"
)

func TestSlidingWindow(t *testing.T) {
	rdb := redisClient.NewRedis("localhost:6379", "ratelimiter")
	ctx := context.Background()
	key := "test:sliding_window"

	rdb.Client.Del(ctx, key)
	defer rdb.Client.Del(ctx, key)

	res, err := limiter.SlidingWindow(ctx, rdb, key, 3, 60)
	if err != nil {
		t.Fatalf("Sliding Window returned an error : %v", err)
	}

	if !res.Allowed {
		t.Errorf("Expected Allowed, but got rejected")
	}

	if res.Remaining != 2 {
		t.Errorf("excepted remaining 2, got %d", res.Remaining)
	}

	res2, err := limiter.SlidingWindow(ctx, rdb, key, 3, 60)
	if err != nil {
		t.Fatalf("Unexcepted error : %v", err)
	}

	if !res2.Allowed {
		t.Errorf("excepted allowed, got rejected")
	}

	if res2.Remaining != 1 {
		t.Errorf("excepted remaining 1, got %d", res2.Remaining)
	}

	res3, err := limiter.SlidingWindow(ctx, rdb, key, 3, 60)
	if err != nil {
		t.Fatalf("Unexcepted error : %v", err)
	}

	if !res3.Allowed {
		t.Errorf("excepted allowed, got rejected")
	}

	if res3.Remaining != 0 {
		t.Errorf("excepted remaining 0, got %d", res3.Remaining)
	}

	res4, err := limiter.SlidingWindow(ctx, rdb, key, 3, 60)
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

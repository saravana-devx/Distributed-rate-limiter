package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	redisClient "ratelimiter/internal/redis"
)

var fixedWindowScript = redis.NewScript(`
	local key = KEYS[1]
	local limit = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])

	-- increment the count for the key
	local count = redis.call("INCR", key)

	-- if first request, set the expiry window
	if count == 1 then
		redis.call("EXPIRE", key, window)
	end

	local ttl = redis.call("TTL", key)
	local reset_at = now + ttl

	-- {allowed, remaining, reset_at} | 1 = allowed, 0 = rejected
	if count > limit then
		return {0, 0, reset_at}
	end
	return {1, limit - count, reset_at}
`)

type FixedWindowResult struct {
	Allowed   bool
	Remaining int64
	ResetAt   int64
}

func FixedWindow(ctx context.Context, rdb *redisClient.Redis, key string, limit int, windowSecs int) (*FixedWindowResult, error) {
	now := time.Now().Unix()
	res, err := rdb.RunScript(ctx, fixedWindowScript, []string{key}, limit, windowSecs, now)
	if err != nil {
		return nil, err
	}
	vals := res.([]interface{})
	return &FixedWindowResult{
		Allowed:   vals[0].(int64) == 1,
		Remaining: vals[1].(int64),
		ResetAt:   vals[2].(int64),
	}, nil
}

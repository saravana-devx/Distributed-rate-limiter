package limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	redisClient "ratelimiter/internal/redis"
)

var slidingWindowScript = redis.NewScript(`
	local key = KEYS[1]
	local limit = tonumber(ARGV[1])
	local windowSize = tonumber(ARGV[2])

	local now = tonumber(redis.call("TIME")[1])
	local windowStart = now - windowSize

	-- remove all entries outside the window
	redis.call("ZREMRANGEBYSCORE", key, "-inf", windowStart)

	-- add current request timestamps
	redis.call("ZADD", key, now, now)
	-- count how many requests are in the window

	-- auto cleanup so key doesn't resides in redis forever
	local count = redis.call("ZCARD", key)
	redis.call("EXPIRE", key, windowSize)

	if count > limit then
		return {0, 0}
	else
		return {1, limit - count}
	end
`)

type SlidingWindowResult struct {
	Allowed   bool
	Remaining int64
}

func SlidingWindow(ctx context.Context, rdb *redisClient.Redis, key string, limit int, windowSize int) (*SlidingWindowResult, error) {
	res, err := rdb.RunScript(ctx, slidingWindowScript, []string{key}, limit, windowSize)
	if err != nil {
		return nil, err
	}

	vals := res.([]interface{})
	return &SlidingWindowResult{
		Allowed:   vals[0].(int64) == 1,
		Remaining: vals[1].(int64),
	}, nil
}

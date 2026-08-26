package limiter

import (
	"context"

	"github.com/redis/go-redis/v9"
	redisClient "ratelimiter/internal/redis"
)

var tokenBucketScript = redis.NewScript(`
	local key = KEYS[1]
	local capacity = tonumber(ARGV[1])
	local refill_rate = tonumber(ARGV[2])
	local now = tonumber(redis.call("TIME")[1])

	local tokens = tonumber(redis.call("HGET",key,"tokens"))

	if tokens == nil then
		tokens = capacity
	end

	local last_refill = tonumber(redis.call("HGET",key,"last_refill"))
	if last_refill == nil then
		last_refill = now
	end

	local elapsed = now - last_refill

	if elapsed > 0 then
		tokens = math.min(tokens + elapsed * refill_rate, capacity)
	end

	
	if tokens < 1 then 
		return {0, 0}
	end

	redis.call("HSET",key, "tokens", tokens - 1, "last_refill", now)

	--if tokens < 1 then 
	--	return {0, 0, now + math.ceil((1 - tokens) / refill_rate)}
	--end

	return {1, tokens - 1}
`)

type TokenBucketResult struct {
	Allowed   bool
	Remaining int64
}

func TokenBucket(ctx context.Context, rdb *redisClient.Redis, key string, capacity int, refillRate int) (*TokenBucketResult, error) {
	res, err := rdb.RunScript(ctx, tokenBucketScript, []string{key}, capacity, refillRate)
	if err != nil {
		return nil, err
	}

	vals := res.([]interface{})
	return &TokenBucketResult{
		Allowed:   vals[0].(int64) == 1,
		Remaining: vals[1].(int64),
	}, nil
}

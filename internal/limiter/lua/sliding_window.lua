local key = KEYS[1]
local limit = tonumber(ARGV[1])
local windowSize = tonumber(ARGV[2])

local time = redis.call("TIME")
local now = tonumber(time[1])
local windowStart = now - windowSize

-- remove all entries outside the window
redis.call("ZREMRANGEBYSCORE", key, "-inf", windowStart)

-- add current request timestamp; member must be unique per request so that
-- multiple requests within the same second don't collide and overwrite
-- each other's sorted-set entry (which would keep ZCARD from increasing)
local member = now .. "-" .. time[2] .. "-" .. math.random()
redis.call("ZADD", key, now, member)
-- count how many requests are in the window
local count = redis.call("ZCARD", key)

-- auto cleanup so key doesn't resides in redis forever
redis.call("EXPIRE", key, windowSize)

-- reset_at is when the oldest entry currently in the window falls out of it
local oldest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")
local reset_at = now + windowSize
if oldest[2] ~= nil then
	reset_at = tonumber(oldest[2]) + windowSize
end

if count > limit then
	return {0, 0, reset_at}
else
	return {1, limit - count, reset_at}
end

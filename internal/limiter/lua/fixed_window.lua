local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

-- increment the count for the key
local count = redis.call("INCR", key)

-- if first request, set the expiry window
if count == 1 then
	redis.call("EXPIRE", key, window)
end

local ttl = redis.call("TTL", key)
local now = tonumber(redis.call("TIME")[1])
local reset_at = now + ttl

-- {allowed, remaining, reset_at} | 1 = allowed, 0 = rejected
if count > limit then
	return {0, 0, reset_at}
end
return {1, limit - count, reset_at}

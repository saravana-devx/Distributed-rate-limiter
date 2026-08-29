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
	-- reset_at is when enough tokens will have refilled to allow the next request
	local reset_at = now + math.ceil((1 - tokens) / refill_rate)
	return {0, 0, reset_at}
end

redis.call("HSET",key, "tokens", tokens - 1, "last_refill", now)

return {1, tokens - 1, now}

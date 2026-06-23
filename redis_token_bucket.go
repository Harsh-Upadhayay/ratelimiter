package ratelimiter

import "strconv"

type RedisTokenBucket struct {
	config tokenBucketConfig
}

const tokenScale = 1000

func NewRedisTokenBucket(capacity int, refillRate float64) (*RedisTokenBucket, error) {
	config, err := newTokenBucketConfig(capacity, refillRate)
	if err != nil {
		return nil, err
	}

	if int(refillRate*float64(tokenScale)) <= 0 {
		return nil, ErrInvalidRefillRate
	}

	return &RedisTokenBucket{
		config: config,
	}, nil
}

func (tb *RedisTokenBucket) script() string {
	return `
		local key = KEYS[1]
		local capacity = tonumber(ARGV[1])
		local refillRate = tonumber(ARGV[2])
		local tokenScale = tonumber(ARGV[3])
		local oneReqCost = tokenScale
		local ttlSeconds = math.ceil(capacity / refillRate)
		
		local now = redis.call("TIME")
		local nowMs = tonumber(now[1]) * 1000 + math.floor(tonumber(now[2]) / 1000)

		local state = redis.call("HMGET", key, "tokens", "last_refill_ms")
		local tokens = state[1]
		local lastRefillMs = state[2]

		if not tokens or not lastRefillMs then 
			tokens = capacity
			lastRefillMs = nowMs
		else
			tokens = tonumber(tokens)
			lastRefillMs = tonumber(lastRefillMs)

		local elapsedMs = nowMs - lastRefillMs
		if elapsedMs > 0 then 
			local refilled = math.floor((elapsedMs * refillRate) / 1000)
			if refilled > 0 then
				tokens = math.min(capacity, tokens + refilled)

				local consumedMs = math.ceil((refilled * 1000) / refillRate)
				lastRefillMs = math.min(nowMs, lastRefillMs + consumedMs)
		
				end
			end
		end
		
		if tokens >= oneReqCost then
			tokens = tokens - oneReqCost

			redis.call("HSET", key, "tokens", tokens, "last_refill_ms", lastRefillMs)
			redis.call("EXPIRE", key, ttlSeconds)

			return {1, math.floor(tokens / tokenScale), 0}
		end

		redis.call("HSET", key, "tokens", tokens, "last_refill_ms", lastRefillMs)
		redis.call("EXPIRE", key, ttlSeconds)

	
		local deficit = oneReqCost - tokens
		local retryAfterSeconds = math.ceil(deficit / refillRate)

		return {0, math.floor(tokens/tokenScale), retryAfterSeconds}
	`
}

func (tb *RedisTokenBucket) args() []string {
	capacityScaled := tb.config.capacity * tokenScale
	refillUnitsPerSecond := int(tb.config.refillRate * float64(tokenScale))

	return []string{
		strconv.Itoa(capacityScaled),
		strconv.Itoa(refillUnitsPerSecond),
		strconv.Itoa(tokenScale),
	}
}

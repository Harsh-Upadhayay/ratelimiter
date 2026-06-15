package ratelimiter

import (
	"strconv"
	"time"
)

type RedisFixedWindow struct {
	limit          int
	windowDuration time.Duration
}

func NewRedisFixedWindow(limit int, windowDuration time.Duration) (*RedisFixedWindow, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	if windowDuration <= 0 {
		return nil, ErrInvalidWindowDuration
	}
	return &RedisFixedWindow{
		limit:          limit,
		windowDuration: windowDuration,
	}, nil
}

func (fw *RedisFixedWindow) script() string {
	script := `
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local windowDuration = tonumber(ARGV[2])

		local remaining = redis.call("GET", key)

		if not remaining then
			redis.call("SET", key, limit - 1, "EX", windowDuration)
			return {1, limit - 1, 0}
		end

		remaining = tonumber(remaining)

		if remaining > 0 then 
			redis.call("DECR", key)
			return {1, remaining - 1, 0}
		end

		local ttl = redis.call("TTL", key)
		if ttl == -1 then
			return redis.error_reply("rate limiter key has no ttl")
		end
		
		return {0, 0, ttl}
	`
	return script
}

func (fw *RedisFixedWindow) args() []string {
	args := []string{}
	args = append(args, strconv.Itoa(fw.limit))
	args = append(args, strconv.Itoa(int(fw.windowDuration.Seconds())))
	return args
}

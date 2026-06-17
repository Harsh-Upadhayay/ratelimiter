package ratelimiter

import (
	"strconv"
	"time"
)

type RedisFixedWindow struct {
	config fixedWindowConfig
}

func NewRedisFixedWindow(limit int, windowDuration time.Duration) (*RedisFixedWindow, error) {
	config, err := newFixedWindowConfig(limit, windowDuration)
	if err != nil {
		return nil, err
	}
	return &RedisFixedWindow{
		config: config,
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
	args = append(args, strconv.Itoa(fw.config.requestLimit))
	args = append(args, strconv.Itoa(int(fw.config.windowDuration.Seconds())))
	return args
}

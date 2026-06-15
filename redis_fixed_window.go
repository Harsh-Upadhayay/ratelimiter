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
	// TODO: Write lua code, after learning syntax.
	return ""
}

func (fw *RedisFixedWindow) args() []string {
	args := []string{}
	args = append(args, strconv.Itoa(fw.limit))
	args = append(args, strconv.Itoa(int(fw.windowDuration.Seconds())))
	return args
}

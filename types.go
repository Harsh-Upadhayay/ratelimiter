package ratelimiter

import (
	"context"
	"time"
)

// memoryAlgorithmState is an interface that all in-memory algorithm states must implement.
// It is used to store the state of each key in the MemoryLimiter.
type memoryAlgorithmState interface {
	isMemoryAlgorithmState()
}

// memoryAlgorithm is an interface that all in-memory rate limiting algorithms must implement.
type memoryAlgorithm interface {
	Decide(now time.Time, state memoryAlgorithmState, exists bool) (Result, memoryAlgorithmState, error)
}

// Result represents the outcome of a rate limit check,
// including whether the request is allowed,
// how many requests are remaining in the current window,
// and how long to wait before retrying if the limit has been exceeded.
type Result struct {
	Allowed    bool
	Remaining  int
	RetryAfter time.Duration
}

type fixedWindowConfig struct {
	requestLimit   int
	windowDuration time.Duration
}

func newFixedWindowConfig(requestLimit int, windowDuration time.Duration) (fixedWindowConfig, error) {
	if requestLimit <= 0 {
		return fixedWindowConfig{}, ErrInvalidLimit
	}
	if windowDuration <= 0 {
		return fixedWindowConfig{}, ErrInvalidWindowDuration
	}

	return fixedWindowConfig{
		requestLimit:   requestLimit,
		windowDuration: windowDuration,
	}, nil
}

type tokenBucketConfig struct {
	capacity   int
	refillRate float64
}

func newTokenBucketConfig(capacity int, refillRate float64) (tokenBucketConfig, error) {
	if capacity <= 0 {
		return tokenBucketConfig{}, ErrInvalidCapacity
	}

	if refillRate <= 0 {
		return tokenBucketConfig{}, ErrInvalidRefillRate
	}

	return tokenBucketConfig{
		capacity:   capacity,
		refillRate: refillRate,
	}, nil
}

type redisAlgorithm interface {
	script() string
	args() []string
}

type redisAdapter interface {
	eval(ctx context.Context, script string, keys []string, args []string) (any, error)
}

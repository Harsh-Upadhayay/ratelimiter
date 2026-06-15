package ratelimiter

import (
	"context"
	"time"
)

// algorithmState is an interface that all algorithm states must implement.
// It is used to store the state of each key in the Limiter.
type algorithmState interface {
	isAlgorithmState()
}

// algorithm is an interface that all rate limiting algorithms must implement.
type algorithm interface {
	Decide(now time.Time, state algorithmState, exists bool) (Result, algorithmState, error)
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

type redisAlgorithm interface {
	script() string
	args() []string
}

type redisAdapter interface {
	eval(ctx context.Context, script string, keys []string, args []string) (any, error)
}

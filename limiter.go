package ratelimiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, key string) (Result, error)
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

// Compile-time assertions to ensure that MemoryLimiter and RedisLimiter implement the Limiter interface.
var _ Limiter = (*MemoryLimiter)(nil) // Convert the value nill to the type *MemoryLimiter: (T)(v) converts v to T.
var _ Limiter = (*RedisLimiter)(nil)

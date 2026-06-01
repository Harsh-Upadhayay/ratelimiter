package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

type userState struct {
	windowStartTime  time.Time
	consumedRequests int
}

// Limiter implements a fixed window rate limiter that tracks the
// number of requests made by each key within a
// specified time window. It allows a certain number of requests
// per window and resets the count when the window expires.
type Limiter struct {
	limit          int
	windowDuration time.Duration
	states         map[string]userState
	mu             sync.Mutex
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

// ErrEmptyKey is returned when the rate limit key is empty
var ErrEmptyKey = errors.New("rate limit key is required")

// ErrInvalidLimit is returned when the limit is not greater than 0
var ErrInvalidLimit = errors.New("limit must be greater than 0")

// ErrInvalidWindowDuration is returned when the window duration is not greater than 0
var ErrInvalidWindowDuration = errors.New("window duration must be greater than 0")

// NewLimiter creates a new Limiter with the specified limit and window duration.
func NewLimiter(limit int, windowDuration time.Duration) (*Limiter, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	if windowDuration <= 0 {
		return nil, ErrInvalidWindowDuration
	}

	limiter := &Limiter{
		limit:          limit,
		windowDuration: windowDuration,
		states:         make(map[string]userState),
	}

	return limiter, nil
}

// Allow checks if a request with the given key is allowed at the current time.
// It returns a Result indicating whether the request is allowed,
// how many requests are remaining in the current window,
// and how long to wait before retrying if the limit has been exceeded.
func (l *Limiter) Allow(key string, now time.Time) (Result, error) {
	if key == "" {
		return Result{}, ErrEmptyKey
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	curState, exists := l.states[key]

	result, updatedState := decideFixedWindow(
		now,
		l.limit,
		l.windowDuration,
		curState,
		exists,
	)
	l.states[key] = updatedState
	return result, nil
}

// Fixed window algorithm implementation.
// Doesn't interact with the storage layer
// Returns the result of the rate limit check and the updated user state.
// Parameter ordering is config first, then state
func decideFixedWindow(
	now time.Time,
	requestLimit int,
	windowDuration time.Duration,
	state userState,
	exists bool,
) (Result, userState) {
	if !exists {
		state = userState{windowStartTime: now, consumedRequests: 1}
		return Result{
			Allowed:    true,
			Remaining:  requestLimit - 1,
			RetryAfter: 0,
		}, state
	}

	if !state.windowStartTime.Add(windowDuration).After(now) {
		state.windowStartTime = now
		state.consumedRequests = 1
		return Result{
			Allowed:    true,
			Remaining:  requestLimit - 1,
			RetryAfter: 0,
		}, state
	}
	if state.consumedRequests < requestLimit {
		state.consumedRequests++
		return Result{
			Allowed:    true,
			Remaining:  requestLimit - state.consumedRequests,
			RetryAfter: 0,
		}, state
	}
	return Result{
		Allowed:    false,
		Remaining:  0,
		RetryAfter: state.windowStartTime.Add(windowDuration).Sub(now),
	}, state
}

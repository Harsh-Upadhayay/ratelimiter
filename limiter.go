package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

type algorithmState interface {
	isAlgorithmState()
}

type fixedWindowState struct {
	windowStartTime  time.Time
	consumedRequests int
}

func (f fixedWindowState) isAlgorithmState() {}

type algorithm interface {
	Decide(now time.Time, state algorithmState, exists bool) (Result, algorithmState, error)
}

// FixedWindow implements a fixed window rate limiter algorithm that tracks the
// number of requests made by each key within a
// specified time window. It allows a certain number of requests
// per window and resets the count when the window expires.
type FixedWindow struct {
	requestLimit   int
	windowDuration time.Duration
}

func NewFixedWindow(requestLimit int, windowDuration time.Duration) (*FixedWindow, error) {
	if requestLimit <= 0 {
		return nil, ErrInvalidLimit
	}
	if windowDuration <= 0 {
		return nil, ErrInvalidWindowDuration
	}

	return &FixedWindow{
		requestLimit:   requestLimit,
		windowDuration: windowDuration,
	}, nil
}

// Fixed window algorithm implementation.
// Doesn't interact with the storage layer
// Returns the result of the rate limit check and the updated user state.
// Parameter ordering is config first, then state
func (f FixedWindow) Decide(now time.Time, state algorithmState, exists bool) (Result, algorithmState, error) {
	if !exists {
		fwState := fixedWindowState{windowStartTime: now, consumedRequests: 1}
		return Result{
			Allowed:    true,
			Remaining:  f.requestLimit - 1,
			RetryAfter: 0,
		}, fwState, nil
	}

	fwState, ok := state.(fixedWindowState)

	if !ok {
		return Result{}, state, ErrUnsupportedAlgorithmState
	}

	if !fwState.windowStartTime.Add(f.windowDuration).After(now) {
		fwState.windowStartTime = now
		fwState.consumedRequests = 1
		return Result{
			Allowed:    true,
			Remaining:  f.requestLimit - 1,
			RetryAfter: 0,
		}, fwState, nil
	}
	if fwState.consumedRequests < f.requestLimit {
		fwState.consumedRequests++
		return Result{
			Allowed:    true,
			Remaining:  f.requestLimit - fwState.consumedRequests,
			RetryAfter: 0,
		}, fwState, nil
	}
	return Result{
		Allowed:    false,
		Remaining:  0,
		RetryAfter: fwState.windowStartTime.Add(f.windowDuration).Sub(now),
	}, fwState, nil
}

type TokenBucket struct {
	capacity   int
	refillRate float64
}

type tokenBucketState struct {
	lastRefillTime  time.Time
	availableTokens float64
}

func NewTokenBucket(capacity int, refilRate float64) (*TokenBucket, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	if refilRate <= 0 {
		return nil, ErrInvalidRefilRate
	}

	return &TokenBucket{
		capacity:   capacity,
		refillRate: refilRate,
	}, nil
}

func (t tokenBucketState) isAlgorithmState() {}

func (t TokenBucket) Decide(now time.Time, state algorithmState, exists bool) (Result, algorithmState, error) {
	if !exists {
		tbState := tokenBucketState{lastRefillTime: now, availableTokens: float64(t.capacity - 1)}
		return Result{
			Allowed:    true,
			Remaining:  t.capacity - 1,
			RetryAfter: 0,
		}, tbState, nil
	}

	tbState, ok := state.(tokenBucketState)

	if !ok {
		return Result{}, state, ErrUnsupportedAlgorithmState
	}

	elapsedTime := now.Sub(tbState.lastRefillTime)
	if elapsedTime > 0 {
		tbState.lastRefillTime = now
	} else {
		elapsedTime = 0
	}

	refilled := elapsedTime.Seconds() * t.refillRate
	tbState.availableTokens = min(float64(t.capacity), tbState.availableTokens+refilled)

	if tbState.availableTokens >= 1 {
		tbState.availableTokens -= 1
		return Result{
			Allowed:    true,
			Remaining:  int(tbState.availableTokens),
			RetryAfter: 0,
		}, tbState, nil
	}

	return Result{
		Allowed:    false,
		Remaining:  0,
		RetryAfter: time.Duration((float64(1) - tbState.availableTokens) / t.refillRate * float64(time.Second)),
	}, tbState, nil
}

type Limiter struct {
	algo   algorithm
	states map[string]algorithmState
	mu     sync.Mutex
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

// ErrEmptyKey is returned when the rate limit key is empty;
var ErrEmptyKey = errors.New("rate limit key is required")

// ErrInvalidLimit is returned when the limit is not greater than 0;
var ErrInvalidLimit = errors.New("limit must be greater than 0")

// ErrInvalidWindowDuration is returned when the window duration is not greater than 0;
var ErrInvalidWindowDuration = errors.New("window duration must be greater than 0")

// ErrUnsupportedAlgorithmState is returned when the algorithm receives a state that it doesn't recognize;
var ErrUnsupportedAlgorithmState = errors.New("unsupported algorithm state passed to the algorithm")

// ErrNilAlgorithm is returned when trying to create a Limiter with a nil algorithm;
var ErrNilAlgorithm = errors.New("algorithm cannot be nil")

var ErrInvalidCapacity = errors.New("capacity must be greater than 0")

var ErrInvalidRefilRate = errors.New("refil rate must be greater than 0")

// NewLimiter creates a new Limiter with the specified algorithm
func NewLimiter(algo algorithm) (*Limiter, error) {

	if algo == nil {
		return nil, ErrNilAlgorithm
	}

	limiter := &Limiter{
		algo:   algo,
		states: make(map[string]algorithmState),
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

	result, updatedState, err := l.algo.Decide(
		now,
		curState,
		exists,
	)

	if err != nil {
		return Result{}, err
	}

	l.states[key] = updatedState
	return result, nil
}

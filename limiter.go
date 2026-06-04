package ratelimiter

import (
	"sync"
	"time"
)

// Limiter is the main struct that users interact with to perform rate limiting checks.
// It holds a reference to the rate limiting algorithm and a map of states for each key.
type Limiter struct {
	algo   algorithm
	states map[string]algorithmState
	mu     sync.Mutex
}

// NewLimiter creates a new Limiter with the specified algorithm.
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

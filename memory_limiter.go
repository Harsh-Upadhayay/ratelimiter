package ratelimiter

import (
	"time"
)

// MemoryLimiter is the in-process limiter that users interact with to perform rate limiting checks.
// It holds a reference to the rate limiting algorithm and a state store for each key.
type MemoryLimiter struct {
	algo  memoryAlgorithm
	store StateStore
}

// NewMemoryLimiter creates a new MemoryLimiter with the specified algorithm.
func NewMemoryLimiter(algo memoryAlgorithm, store StateStore) (*MemoryLimiter, error) {
	if algo == nil {
		return nil, ErrNilAlgorithm
	}

	limiter := &MemoryLimiter{
		algo:  algo,
		store: store,
	}

	return limiter, nil
}

// Allow checks if a request with the given key is allowed at the current time.
// It returns a Result indicating whether the request is allowed,
// how many requests are remaining in the current window,
// and how long to wait before retrying if the limit has been exceeded.
func (l *MemoryLimiter) Allow(key string, now time.Time) (Result, error) {
	const maxCASRetries = 10

	if key == "" {
		return Result{}, ErrEmptyKey
	}

	for attempt := 0; attempt < maxCASRetries; attempt++ {
		curState, version, exists, err := l.store.Get(key)
		if err != nil {
			return Result{}, err
		}

		result, updatedState, err := l.algo.Decide(
			now,
			curState,
			exists,
		)
		if err != nil {
			return Result{}, err
		}

		ok, err := l.store.CompareAndSwap(key, version, updatedState)

		if err != nil {
			return Result{}, err
		}

		if ok {
			return result, nil
		}
	}
	return Result{}, ErrCASConflict
}

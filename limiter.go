package ratelimiter

import (
	"time"
)

// Limiter is the main struct that users interact with to perform rate limiting checks.
// It holds a reference to the rate limiting algorithm and a map of states for each key.
type Limiter struct {
	algo  algorithm
	store StateStore
}

// NewLimiter creates a new Limiter with the specified algorithm.
func NewLimiter(algo algorithm, store StateStore) (*Limiter, error) {
	if algo == nil {
		return nil, ErrNilAlgorithm
	}

	limiter := &Limiter{
		algo:  algo,
		store: store,
	}

	return limiter, nil
}

// Allow checks if a request with the given key is allowed at the current time.
// It returns a Result indicating whether the request is allowed,
// how many requests are remaining in the current window,
// and how long to wait before retrying if the limit has been exceeded.
func (l *Limiter) Allow(key string, now time.Time) (Result, error) {
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

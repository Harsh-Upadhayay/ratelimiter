package ratelimiter

import (
	"context"
)

// MemoryLimiter is the in-process limiter that users interact with to perform rate limiting checks.
// It holds a reference to the rate limiting algorithm and a state store for each key.
type MemoryLimiter struct {
	algo  memoryAlgorithm
	store StateStore
	clock clock
}

// NewMemoryLimiter creates a new MemoryLimiter with the specified algorithm.
func NewMemoryLimiter(algo memoryAlgorithm, store StateStore) (*MemoryLimiter, error) {
	if algo == nil {
		return nil, ErrNilAlgorithm
	}

	limiter := &MemoryLimiter{
		algo:  algo,
		store: store,
		clock: realClock{},
	}

	return limiter, nil
}

// Allow checks if a request with the given key is allowed at the current time.
// It returns a Result indicating whether the request is allowed,
// how many requests are remaining in the current window,
// and how long to wait before retrying if the limit has been exceeded.
func (ml *MemoryLimiter) Allow(ctx context.Context, key string) (Result, error) {

	if ctx.Err() != nil {
		return Result{}, ctx.Err()
	}

	const maxCASRetries = 10

	if key == "" {
		return Result{}, ErrEmptyKey
	}

	for attempt := 0; attempt < maxCASRetries; attempt++ {
		curState, version, exists, err := ml.store.Get(key)
		if err != nil {
			return Result{}, err
		}

		result, updatedState, err := ml.algo.Decide(
			ml.clock.now(),
			curState,
			exists,
		)
		if err != nil {
			return Result{}, err
		}

		ok, err := ml.store.CompareAndSwap(key, version, updatedState)

		if err != nil {
			return Result{}, err
		}

		if ok {
			return result, nil
		}
	}
	return Result{}, ErrCASConflict
}

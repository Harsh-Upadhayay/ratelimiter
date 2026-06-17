package ratelimiter

import "time"

// tokenBucketState is the state for the MemoryTokenBucket algorithm,
// which keeps track of the last time the bucket was refilled
// and the number of available tokens.
type tokenBucketState struct {
	lastRefillTime  time.Time
	availableTokens float64
}

// isMemoryAlgorithmState is a marker method to ensure that
// tokenBucketState implements the memoryAlgorithmState interface.
func (t tokenBucketState) isMemoryAlgorithmState() {}

// MemoryTokenBucket implements a token bucket rate limiter algorithm.
// It sets the token bucket's capacity and refill rate. Each request consumes a token.
type MemoryTokenBucket struct {
	config tokenBucketConfig
}

// NewMemoryTokenBucket creates a new MemoryTokenBucket with the specified capacity and refill rate.
// Returns an error if the capacity is not positive or if the refill rate is not positive.
func NewMemoryTokenBucket(capacity int, refillRate float64) (*MemoryTokenBucket, error) {
	config, err := newTokenBucketConfig(capacity, refillRate)
	if err != nil {
		return nil, err
	}

	return &MemoryTokenBucket{
		config: config,
	}, nil
}

// Decide applies the token-bucket algorithm for the current state and returns
// the rate-limit result plus the state that should be stored for the key.
func (t MemoryTokenBucket) Decide(now time.Time, state memoryAlgorithmState, exists bool) (Result, memoryAlgorithmState, error) {
	if !exists {
		tbState := tokenBucketState{lastRefillTime: now, availableTokens: float64(t.config.capacity - 1)}
		return Result{
			Allowed:    true,
			Remaining:  t.config.capacity - 1,
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

	refilled := elapsedTime.Seconds() * t.config.refillRate
	tbState.availableTokens = min(float64(t.config.capacity), tbState.availableTokens+refilled)

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
		RetryAfter: time.Duration((float64(1) - tbState.availableTokens) / t.config.refillRate * float64(time.Second)),
	}, tbState, nil
}

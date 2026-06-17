package ratelimiter

import "time"

// fixedWindowState used by the MemoryFixedWindow algorithm to
// track the start time of the current window and the
// number of consumed requests within that window.
type fixedWindowState struct {
	windowStartTime  time.Time
	consumedRequests int
}

// isMemoryAlgorithmState is a marker method to ensure that
// fixedWindowState implements the memoryAlgorithmState interface.
func (f fixedWindowState) isMemoryAlgorithmState() {}

// MemoryFixedWindow implements a fixed window rate limiter algorithm that tracks the
// number of requests made by each key within a
// specified time window. It allows a certain number of requests
// per window and resets the count when the window expires.
type MemoryFixedWindow struct {
	config fixedWindowConfig
}

// NewMemoryFixedWindow creates a new MemoryFixedWindow with the specified request limit and window duration.
// Returns an error if the request limit is not positive or if the window duration is not positive.
func NewMemoryFixedWindow(requestLimit int, windowDuration time.Duration) (*MemoryFixedWindow, error) {
	config, err := newFixedWindowConfig(requestLimit, windowDuration)
	if err != nil {
		return nil, err
	}

	return &MemoryFixedWindow{
		config: config,
	}, nil
}

// Decide applies the fixed-window algorithm for the current state and returns
// the rate-limit result plus the state that should be stored for the key.
func (f MemoryFixedWindow) Decide(now time.Time, state memoryAlgorithmState, exists bool) (Result, memoryAlgorithmState, error) {
	if !exists {
		fwState := fixedWindowState{windowStartTime: now, consumedRequests: 1}
		return Result{
			Allowed:    true,
			Remaining:  f.config.requestLimit - 1,
			RetryAfter: 0,
		}, fwState, nil
	}

	fwState, ok := state.(fixedWindowState)

	if !ok {
		return Result{}, state, ErrUnsupportedAlgorithmState
	}

	if !fwState.windowStartTime.Add(f.config.windowDuration).After(now) {
		fwState.windowStartTime = now
		fwState.consumedRequests = 1
		return Result{
			Allowed:    true,
			Remaining:  f.config.requestLimit - 1,
			RetryAfter: 0,
		}, fwState, nil
	}
	if fwState.consumedRequests < f.config.requestLimit {
		fwState.consumedRequests++
		return Result{
			Allowed:    true,
			Remaining:  f.config.requestLimit - fwState.consumedRequests,
			RetryAfter: 0,
		}, fwState, nil
	}
	return Result{
		Allowed:    false,
		Remaining:  0,
		RetryAfter: fwState.windowStartTime.Add(f.config.windowDuration).Sub(now),
	}, fwState, nil
}

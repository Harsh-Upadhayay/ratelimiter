package ratelimiter

import "time"

// fixedWindowState used by the FixedWindow algorithm to
// track the start time of the current window and the
// number of consumed requests within that window.
type fixedWindowState struct {
	windowStartTime  time.Time
	consumedRequests int
}

// isAlgorithmState is a marker method to ensure that
// fixedWindowState implements the algorithmState interface.
func (f fixedWindowState) isAlgorithmState() {}

// FixedWindow implements a fixed window rate limiter algorithm that tracks the
// number of requests made by each key within a
// specified time window. It allows a certain number of requests
// per window and resets the count when the window expires.
type FixedWindow struct {
	requestLimit   int
	windowDuration time.Duration
}

// NewFixedWindow creates a new FixedWindow with the specified request limit and window duration.
// Returns an error if the request limit is not positive or if the window duration is not positive.
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

// Decide applies the fixed-window algorithm for the current state and returns
// the rate-limit result plus the state that should be stored for the key.
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

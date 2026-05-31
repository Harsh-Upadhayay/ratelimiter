package ratelimiter

import (
	"sync"
	"time"
)

type userState struct {
	windowStartTime  time.Time
	consumedRequests int
}

type Limiter struct {
	limit          int
	windowDuration time.Duration
	states         map[string]userState
	mu             sync.Mutex
}

func NewLimiter(limit int, windowDuration time.Duration) *Limiter {

	if limit <= 0 {
		panic("limit must be must be greater than 0")
	}
	if windowDuration <= 0 {
		panic("window duration must be greater than 0")
	}

	limiter := &Limiter{
		limit:          limit,
		windowDuration: windowDuration,
		states:         make(map[string]userState)}

	return limiter
}

func (l *Limiter) Allow(userID string, now time.Time) bool {

	if userID == "" {
		panic("invalid user ID")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	curUserState, exists := l.states[userID]

	// New user
	if !exists {
		curUserState = userState{now, 1}
		l.states[userID] = curUserState
		return true
	}

	// Expired window: reset
	if !curUserState.windowStartTime.Add(l.windowDuration).After(now) {
		curUserState.windowStartTime = now
		curUserState.consumedRequests = 1
		l.states[userID] = curUserState
		return true
	}
	// Active window with available limit
	if curUserState.consumedRequests < l.limit {
		curUserState.consumedRequests++
		l.states[userID] = curUserState
		return true
	}

	// Active window with limit exhausted
	return false
}

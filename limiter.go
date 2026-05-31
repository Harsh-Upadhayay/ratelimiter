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
	states         map[string]*userState
	mu             sync.Mutex
}

func NewLimiter(limit int, windowDuration time.Duration) *Limiter {

	if limit < 0 {
		panic("limit must be must be greater than 0")
	}
	if windowDuration <= 0 {
		panic("window duration must be greater than 0")
	}

	limiter := new(Limiter)
	limiter.limit = limit
	limiter.windowDuration = windowDuration

	return limiter
}

func (l *Limiter) Allow(userID string, now time.Time) bool {

	if userID == "" {
		panic("invalid userid")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	curUserState, exists := l.states[userID]

	if !exists {
		curUserState = &userState{now, 1}
		l.states[userID] = curUserState
	}

	// Expired window: reset
	if curUserState.windowStartTime.Add(l.windowDuration).Before(now) {
		curUserState.windowStartTime = now
		curUserState.consumedRequests = 1
		return true
	} else {
		// Active window

		if curUserState.consumedRequests < l.limit {
			curUserState.consumedRequests += 1
			return true
		} else {
			return false
		}
	}

}

package ratelimiter

import (
	"sync"
	"time"
)

type clock interface {
	now() time.Time
}

type realClock struct{}

func (realClock) now() time.Time {
	return time.Now()
}

type testClock struct {
	mu sync.Mutex
	t  time.Time
}

func (tc *testClock) now() time.Time {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return tc.t
}

func (tc *testClock) set(nt time.Time) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.t = nt
}

func (tc *testClock) advance(d time.Duration) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.t = tc.t.Add(d)
}

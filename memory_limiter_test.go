package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"
)

func newTestMemoryLimiter(algo memoryAlgorithm) (*MemoryLimiter, *testClock, error) {
	ml, err := NewMemoryLimiter(algo, NewMemoryStore())
	if err != nil {
		return nil, nil, err
	}
	fakeClock := &testClock{t: time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)}
	ml.clock = fakeClock
	return ml, fakeClock, err
}
func TestEmptyKeyError(t *testing.T) {
	algo, err := NewMemoryFixedWindow(1, time.Minute)
	limiter, _, err := newTestMemoryLimiter(algo)

	if err != nil {
		t.Fatalf("new limiter returned error : %v", err)
	}

	_, err = limiter.Allow(context.Background(), "")

	if !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("empty key allowed in limiter.Allow")
	}
}

func TestNilAlgorithmError(t *testing.T) {
	_, _, err := newTestMemoryLimiter(nil)

	if !errors.Is(err, ErrNilAlgorithm) {
		t.Fatalf("limiter created with nil memoryAlgorithm")
	}
}

func TestAllowConcurrentSameKeyMemoryFixedWindow(t *testing.T) {
	type allowOutcome struct {
		allowed bool
		err     error
	}

	fw, err := NewMemoryFixedWindow(1, time.Minute)

	if err != nil {
		t.Fatalf("fixedwindow creation failed with error: %v", err)
	}

	limiter, _, err := newTestMemoryLimiter(fw)

	if err != nil {
		t.Fatalf("limiter creation failed with error: %v", err)
	}

	outcomes := make(chan allowOutcome, 100)
	goroutines := 100

	for range goroutines {
		go func() {
			result, err := limiter.Allow(context.Background(), "key-1")
			outcomes <- allowOutcome{allowed: result.Allowed, err: err}
		}()
	}

	allowed, rejected, errors := 0, 0, 0
	for i := 0; i < goroutines; i++ {
		outcome := <-outcomes
		if outcome.allowed {
			allowed += 1
		} else {
			rejected += 1
		}
		if outcome.err != nil {
			errors += 1
		}
	}

	if allowed != 1 {
		t.Fatalf("expected 1 allowed request, got %d", allowed)
	}
	if rejected != goroutines-1 {
		t.Fatalf("expected %d rejected requests, got %d", goroutines-1, rejected)
	}
	if errors != 0 {
		t.Fatalf("expected 0 errors, got %d", errors)
	}
}

func TestAllowConcurrentSameKeyMemoryTokenBucket(t *testing.T) {
	type allowOutcome struct {
		allowed bool
		err     error
	}

	tb, err := NewMemoryTokenBucket(1, 1)

	if err != nil {
		t.Fatalf("fixedwindow creation failed with error: %v", err)
	}

	limiter, _, err := newTestMemoryLimiter(tb)

	if err != nil {
		t.Fatalf("limiter creation failed with error: %v", err)
	}

	outcomes := make(chan allowOutcome, 100)
	goroutines := 100

	for range goroutines {
		go func() {
			result, err := limiter.Allow(context.Background(), "key-1")
			outcomes <- allowOutcome{allowed: result.Allowed, err: err}
		}()
	}

	allowed, rejected, errors := 0, 0, 0
	for i := 0; i < goroutines; i++ {
		outcome := <-outcomes
		if outcome.allowed {
			allowed += 1
		} else {
			rejected += 1
		}
		if outcome.err != nil {
			errors += 1
		}
	}

	if allowed != 1 {
		t.Fatalf("expected 1 allowed request, got %d", allowed)
	}
	if rejected != goroutines-1 {
		t.Fatalf("expected %d rejected requests, got %d", goroutines-1, rejected)
	}
	if errors != 0 {
		t.Fatalf("expected 0 errors, got %d", errors)
	}
}

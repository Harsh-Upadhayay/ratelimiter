package ratelimiter

import (
	"errors"
	"testing"
	"time"
)

func TestEmptyKeyError(t *testing.T) {
	algo, err := NewFixedWindow(1, time.Minute)
	limiter, err := NewLimiter(algo)

	if err != nil {
		t.Fatalf("new limiter returned error : %v", err)
	}

	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)
	_, err = limiter.Allow("", now)

	if !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("empty key allowed in limiter.Allow")
	}
}

func TestNilAlgorithmError(t *testing.T) {
	_, err := NewLimiter(nil)

	if !errors.Is(err, ErrNilAlgorithm) {
		t.Fatalf("limiter created with nil algorithm")
	}
}

func TestAllowConcurrentSameKeyFixedWindow(t *testing.T) {
	type allowOutcome struct {
		allowed bool
		err     error
	}

	fw, err := NewFixedWindow(1, time.Minute)

	if err != nil {
		t.Fatalf("fixedwindow creation failed with error: %v", err)
	}

	limiter, err := NewLimiter(fw)

	if err != nil {
		t.Fatalf("limiter creation failed with error: %v", err)
	}

	outcomes := make(chan allowOutcome, 100)
	goroutines := 100
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	for range goroutines {
		go func() {
			result, err := limiter.Allow("key-1", now)
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

func TestAllowConcurrentSameKeyTokenBucket(t *testing.T) {
	type allowOutcome struct {
		allowed bool
		err     error
	}

	tb, err := NewTokenBucket(1, 1)

	if err != nil {
		t.Fatalf("fixedwindow creation failed with error: %v", err)
	}

	limiter, err := NewLimiter(tb)

	if err != nil {
		t.Fatalf("limiter creation failed with error: %v", err)
	}

	outcomes := make(chan allowOutcome, 100)
	goroutines := 100
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	for range goroutines {
		go func() {
			result, err := limiter.Allow("key-1", now)
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

package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestInvalidLimitError(t *testing.T) {
	_, err := NewMemoryFixedWindow(0, time.Minute)

	if !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("limiter created with invalid limit")
	}
}

func TestInvalidWindowDurationError(t *testing.T) {
	_, err := NewMemoryFixedWindow(1, 0)

	if !errors.Is(err, ErrInvalidWindowDuration) {
		t.Fatalf("limiter created with invalid window duration")
	}
}

func TestAllowRejectsAfterLimit(t *testing.T) {
	algo, err := NewMemoryFixedWindow(2, time.Minute)
	if err != nil {
		t.Fatalf("new fixed window returned error: %v", err)
	}
	limiter, _, err := newTestMemoryLimiter(algo)

	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}

	result, err := limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}

	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}

	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Allowed {
		t.Fatalf("request over limit was allowed")
	}
}

func TestAllowTracksUsersIndependently(t *testing.T) {
	algo, err := NewMemoryFixedWindow(1, time.Minute)
	if err != nil {
		t.Fatalf("new fixed window returned error: %v", err)
	}
	limiter, _, err := newTestMemoryLimiter(algo)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}

	result, err := limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}
	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Allowed {
		t.Fatalf("request over limit was allowed")
	}

	result, err = limiter.Allow(context.Background(), "key-2")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}
}

func TestAllowResetsAtWindowBoundary(t *testing.T) {
	algo, err := NewMemoryFixedWindow(1, time.Minute)
	if err != nil {
		t.Fatalf("new fixed window returned error: %v", err)
	}
	limiter, clk, err := newTestMemoryLimiter(algo)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}

	result, err := limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}

	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if result.Allowed {
		t.Fatalf("request over limit was allowed")
	}
	clk.advance(time.Minute)
	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request at window boundary was rejected")
	}
}

func TestAllowRejectsBeforeWindowBoundary(t *testing.T) {
	algo, err := NewMemoryFixedWindow(1, time.Minute)
	if err != nil {
		t.Fatalf("new fixed window returned error: %v", err)
	}
	limiter, clk, err := newTestMemoryLimiter(algo)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}
	result, err := limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}
	clk.advance(time.Minute - time.Nanosecond)
	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if result.Allowed {
		t.Fatalf("request before window boundary was allowed")
	}
}

func TestRemainingQuotaCalculation(t *testing.T) {
	algo, err := NewMemoryFixedWindow(2, time.Minute)
	if err != nil {
		t.Fatalf("new fixed window returned error: %v", err)
	}
	limiter, _, err := newTestMemoryLimiter(algo)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}

	result, err := limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Remaining != 1 {
		t.Fatalf("expected remaining limit of 1 != %v ", result.Remaining)
	}

	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Remaining != 0 {
		t.Fatalf("expected remaining limit of 0 != %v ", result.Remaining)
	}

	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Remaining != 0 {
		t.Fatalf("expected remaining limit of 0 != %v ", result.Remaining)
	}
}

func TestRetryAfterCalculation(t *testing.T) {
	algo, err := NewMemoryFixedWindow(1, time.Minute)
	if err != nil {
		t.Fatalf("new fixed window returned error: %v", err)
	}
	limiter, _, err := newTestMemoryLimiter(algo)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}

	result, err := limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.RetryAfter != 0 {
		t.Fatalf("expected retry-after of 0 != %v ", result.RetryAfter)
	}

	result, err = limiter.Allow(context.Background(), "key-1")
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.RetryAfter != time.Minute {
		t.Fatalf("expected retry-after of 1 minute != %v ", result.RetryAfter)
	}
}

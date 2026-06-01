package ratelimiter

import (
	"errors"
	"testing"
	"time"
)

func TestInvalidLimitError(t *testing.T) {
	_, err := NewLimiter(0, time.Minute)

	if !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("limiter created with invalid limit")
	}
}

func TestInvalidWindowDurationError(t *testing.T) {
	_, err := NewLimiter(1, 0)

	if !errors.Is(err, ErrInvalidWindowDuration) {
		t.Fatalf("limiter created with invalid window duration")
	}
}

func TestEmptyKeyError(t *testing.T) {
	limiter, err := NewLimiter(1, time.Minute)

	if err != nil {
		t.Fatalf("new limiter returned error : %v", err)
	}

	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)
	_, err = limiter.Allow("", now)

	if !errors.Is(err, ErrEmptyKey) {
		t.Fatalf("empty key allowed in limiter.Allow")
	}
}

func TestAllowRejectsAfterLimit(t *testing.T) {
	limiter, err := NewLimiter(2, time.Minute)

	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}

	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	result, err := limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}

	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}

	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Allowed {
		t.Fatalf("request over limit was allowed")
	}
}

func TestAllowTracksUsersIndependently(t *testing.T) {
	limiter, err := NewLimiter(1, time.Minute)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	result, err := limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}
	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Allowed {
		t.Fatalf("request over limit was allowed")
	}

	result, err = limiter.Allow("key-2", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}
}

func TestAllowResetsAtWindowBoundary(t *testing.T) {
	limiter, err := NewLimiter(1, time.Minute)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	result, err := limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}

	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if result.Allowed {
		t.Fatalf("request over limit was allowed")
	}

	result, err = limiter.Allow("key-1", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request at window boundary was rejected")
	}
}

func TestAllowRejectsBeforeWindowBoundary(t *testing.T) {
	limiter, err := NewLimiter(1, time.Minute)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	result, err := limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if !result.Allowed {
		t.Fatalf("request within limit was rejected")
	}
	result, err = limiter.Allow("key-1", now.Add(time.Minute-time.Nanosecond))
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}
	if result.Allowed {
		t.Fatalf("request before window boundary was allowed")
	}
}

func TestRemainingQuotaCalculation(t *testing.T) {
	limiter, err := NewLimiter(2, time.Minute)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	result, err := limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Remaining != 1 {
		t.Fatalf("expected remaining limit of 1 != %v ", result.Remaining)
	}

	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Remaining != 0 {
		t.Fatalf("expected remaining limit of 0 != %v ", result.Remaining)
	}

	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.Remaining != 0 {
		t.Fatalf("expected remaining limit of 0 != %v ", result.Remaining)
	}
}

func TestRetryAfterCalculation(t *testing.T) {
	limiter, err := NewLimiter(1, time.Minute)
	if err != nil {
		t.Fatalf("new limiter returned error: %v", err)
	}
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	result, err := limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.RetryAfter != 0 {
		t.Fatalf("expected retry-after of 0 != %v ", result.RetryAfter)
	}

	result, err = limiter.Allow("key-1", now)
	if err != nil {
		t.Fatalf("allow returned error: %v", err)
	}

	if result.RetryAfter != time.Minute {
		t.Fatalf("expected retry-after of 1 minute != %v ", result.RetryAfter)
	}
}

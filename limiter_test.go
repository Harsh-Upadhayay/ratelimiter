package ratelimiter

import (
	"testing"
	"time"
)

func TestAllowRejectsAfterLimit(t *testing.T) {
	limiter := NewLimiter(2, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	if !limiter.Allow("user-1", now) {
		t.Fatalf("request within limit was rejected")
	}

	if !limiter.Allow("user-1", now) {
		t.Fatalf("request within limit was rejected")
	}

	if limiter.Allow("user-1", now) {
		t.Fatalf("request over limit was allowed")
	}
}

func TestAllowTracksUsersIndependently(t *testing.T) {
	limiter := NewLimiter(1, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	if !limiter.Allow("user-1", now) {
		t.Fatalf("request within limit was rejected")
	}
	if limiter.Allow("user-1", now) {
		t.Fatalf("request over limit was allowed")
	}
	if !limiter.Allow("user-2", now) {
		t.Fatalf("request within limit was rejected")
	}
}

func TestAllowResetsAtWindowBoundary(t *testing.T) {
	limiter := NewLimiter(1, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	if !limiter.Allow("user-1", now) {
		t.Fatalf("request within limit was rejected")
	}
	if limiter.Allow("user-1", now) {
		t.Fatalf("request over limit was allowed")
	}
	if !limiter.Allow("user-1", now.Add(time.Minute)) {
		t.Fatalf("request at window boundary was rejected")
	}
}

func TestAllowRejectsBeforeWindowBoundary(t *testing.T) {
	limiter := NewLimiter(1, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	if !limiter.Allow("user-1", now) {
		t.Fatalf("request within limit was rejected")
	}
	if limiter.Allow("user-1", now.Add(time.Minute-time.Nanosecond)) {
		t.Fatalf("request before window boundary was allowed")
	}
}

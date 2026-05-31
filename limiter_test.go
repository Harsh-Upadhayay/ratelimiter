package ratelimiter

import (
	"testing"
	"time"
)

func TestAllowRejectsAfterLimit(t *testing.T) {
	limiter := NewLimiter(2, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	if !limiter.Allow("user-1", now) {
		t.Fatalf("Valid request rejected")
	}

	if !limiter.Allow("user-1", now) {
		t.Fatalf("Valid request rejected")
	}

	if limiter.Allow("user-1", now) {
		t.Fatalf("Invalid request accepted")
	}
}

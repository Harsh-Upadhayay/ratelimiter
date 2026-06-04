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

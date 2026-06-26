package ratelimiter

import (
	"testing"
)

func TestInvalidCapacityError(t *testing.T) {
	_, err := NewMemoryTokenBucket(0, 5)

	if err != ErrInvalidCapacity {
		t.Fatalf("Bucket created with 0 capacity")
	}
	_, err = NewMemoryTokenBucket(-1, 5)

	if err != ErrInvalidCapacity {
		t.Fatalf("Bucket created with negative capacity")
	}
}

func TestInvalidRefillRateError(t *testing.T) {
	_, err := NewMemoryTokenBucket(5, 0)

	if err != ErrInvalidRefillRate {
		t.Fatalf("bucket created with 0 rate")
	}
	_, err = NewMemoryTokenBucket(4, -1)

	if err != ErrInvalidRefillRate {
		t.Fatalf("bucket created with negative refill rate")
	}
}

// TODO: Add tests
// new bucket starts full and consumes one token
// rejects after capacity is exhausted
// refills over elapsed time
// does not allow with less than one token
// caps refill at capacity
// retry-after is based on missing tokens
// backward time does not grant tokens

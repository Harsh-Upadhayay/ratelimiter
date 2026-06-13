package ratelimiter

import (
	"errors"
	"fmt"
)

// ErrEmptyKey is returned when the rate limit key is empty.
var ErrEmptyKey = errors.New("rate limit key is required")

// ErrInvalidLimit is returned when the limit is not greater than 0.
var ErrInvalidLimit = errors.New("limit must be greater than 0")

// ErrInvalidWindowDuration is returned when the window duration is not greater than 0.
var ErrInvalidWindowDuration = errors.New("window duration must be greater than 0")

// ErrUnsupportedAlgorithmState is returned when the algorithm receives a state that it doesn't recognize.
var ErrUnsupportedAlgorithmState = errors.New("unsupported algorithm state passed to the algorithm")

// ErrNilAlgorithm is returned when trying to create a Limiter with a nil algorithm.
var ErrNilAlgorithm = errors.New("algorithm cannot be nil")

// ErrInvalidCapacity is returned when the capacity is not greater than 0.
var ErrInvalidCapacity = errors.New("capacity must be greater than 0")

// ErrInvalidRefillRate is returned when the refill rate is not greater than 0.
var ErrInvalidRefillRate = errors.New("refill rate must be greater than 0")

// ErrCASConflict is returned when the limiter cannot commit state after retrying CAS conflicts.
var ErrCASConflict = errors.New("max number of CAS conflict attempts exhausted")

var ErrInvalidShardCount = fmt.Errorf("shard count must lie between 1 and %d", MAXSHARDSIZE-1)

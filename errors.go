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

// ErrNilAlgorithm is returned when trying to create a MemoryLimiter with a nil algorithm.
var ErrNilAlgorithm = errors.New("algorithm cannot be nil")

// ErrInvalidCapacity is returned when the capacity is not greater than 0.
var ErrInvalidCapacity = errors.New("capacity must be greater than 0")

// ErrInvalidRefillRate is returned when the refill rate is not greater than 0.
var ErrInvalidRefillRate = errors.New("refill rate must be greater than 0")

// ErrCASConflict is returned when the limiter cannot commit state after retrying CAS conflicts.
var ErrCASConflict = errors.New("max number of CAS conflict attempts exhausted")

// ErrInvalidShardCount is returned when the shard count is not between 1 and MAXSHARDSIZE-1.
var ErrInvalidShardCount = fmt.Errorf("shard count must lie between 1 and %d", MAXSHARDSIZE-1)

// ErrUnexpectedRedisResult is returned when Redis returns a value that does not match the limiter result contract.
var ErrUnexpectedRedisResult = errors.New("unexpected redis result")

// ErrNilRedisClient is returned when creating a RedisLimiter with a nil Redis client.
var ErrNilRedisClient = errors.New("redis client cannot be nil")

var ErrNilLimiter = errors.New("limiter cannot be nil")

var ErrNilKeyFunc = errors.New("key func cannot be nil")

var ErrInvalidFailurePolicy = errors.New("failure policy is invalid")

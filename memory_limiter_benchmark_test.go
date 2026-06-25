package ratelimiter

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"
)

func newMemoryFixedWindowLimiter(b *testing.B, limit int, dur time.Duration) *MemoryLimiter {
	b.Helper()
	fw, err := NewMemoryFixedWindow(limit, dur)
	if err != nil {
		b.Fatalf("fixedwindow initialization failed %v", err)
	}

	// Benchmarks use the production realClock (via NewMemoryLimiter's default)
	// rather than the test clock: a shared test clock serializes every Allow on
	// one mutex, which would mask the per-shard scaling these benchmarks measure.
	lim, err := NewMemoryLimiter(fw, NewMemoryStore())
	if err != nil {
		b.Fatalf("limiter initalization failed %v", err)
	}

	return lim
}

func newMemoryTokenBucketLimiter(b *testing.B, capacity int, rate float64) *MemoryLimiter {
	b.Helper()
	fw, err := NewMemoryTokenBucket(capacity, rate)
	if err != nil {
		b.Fatalf("tokenbucket initialization failed %v", err)
	}

	lim, err := NewMemoryLimiter(fw, NewMemoryStore())
	if err != nil {
		b.Fatalf("limiter initalization failed %v", err)
	}

	return lim
}

func BenchmarkAllowSameKeyMemoryFixedWindow(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 1, time.Minute)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow(context.Background(), "key_1")
		if err != nil && !errors.Is(err, ErrCASConflict) {
			b.Fatal(err)
		}
	}

}

func BenchmarkAllowSameKeyMemoryFixedWindowParallel(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 1, time.Minute)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := limiter.Allow(context.Background(), "key_1")
			// ErrCASConflict is an expected outcome under high single-key
			// contention (bounded CAS retry exhausted), not a harness failure.
			if err != nil && !errors.Is(err, ErrCASConflict) {
				b.Fatal(err)
			}
		}
	})

}

func BenchmarkAllowManyKeyMemoryFixedWindow(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 10, time.Nanosecond)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow(context.Background(), keys[i%(len(keys))])
		if err != nil && !errors.Is(err, ErrCASConflict) {
			b.Fatal(err)
		}
	}
}

func BenchmarkAllowManyKeyMemoryFixedWindowParallel(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 1, time.Nanosecond)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, err := limiter.Allow(context.Background(), keys[i%len(keys)])
			// ErrCASConflict is an expected outcome under high single-key
			// contention (bounded CAS retry exhausted), not a harness failure.
			if err != nil && !errors.Is(err, ErrCASConflict) {
				b.Fatal(err)
			}
			i++
		}
	})

}

func BenchmarkAllowSameKeyMemoryTokenBucket(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 1, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow(context.Background(), "key_1")
		if err != nil && !errors.Is(err, ErrCASConflict) {
			b.Fatal(err)
		}
	}

}

func BenchmarkAllowSameKeyMemoryTokenBucketParallel(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 1, 1)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := limiter.Allow(context.Background(), "key_1")
			// ErrCASConflict is an expected outcome under high single-key
			// contention (bounded CAS retry exhausted), not a harness failure.
			if err != nil && !errors.Is(err, ErrCASConflict) {
				b.Fatal(err)
			}
		}
	})

}

func BenchmarkAllowManyKeyMemoryTokenBucket(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 10, 1)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow(context.Background(), keys[i%(len(keys))])
		if err != nil && !errors.Is(err, ErrCASConflict) {
			b.Fatal(err)
		}
	}
}

func BenchmarkAllowManyKeyMemoryTokenBucketParallel(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 1, 1)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, err := limiter.Allow(context.Background(), keys[i%len(keys)])
			// ErrCASConflict is an expected outcome under high single-key
			// contention (bounded CAS retry exhausted), not a harness failure.
			if err != nil && !errors.Is(err, ErrCASConflict) {
				b.Fatal(err)
			}
			i++
		}
	})

}

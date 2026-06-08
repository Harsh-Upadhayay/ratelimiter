package ratelimiter

import (
	"strconv"
	"testing"
	"time"
)

func newFixedWindowLimiter(b *testing.B, limit int, dur time.Duration) *Limiter {
	b.Helper()
	fw, err := NewFixedWindow(limit, dur)
	if err != nil {
		b.Fatalf("fixedwindow initialization failed %v", err)
	}

	lim, err := NewLimiter(fw)
	if err != nil {
		b.Fatalf("limiter initalization failed %v", err)
	}

	return lim
}

func newTokenBucketLimiter(b *testing.B, capacity int, rate float64) *Limiter {
	b.Helper()
	fw, err := NewTokenBucket(capacity, rate)
	if err != nil {
		b.Fatalf("tokenbucket initialization failed %v", err)
	}

	lim, err := NewLimiter(fw)
	if err != nil {
		b.Fatalf("limiter initalization failed %v", err)
	}

	return lim
}

func BenchmarkAllowSameKeyFixedWindow(b *testing.B) {
	limiter := newFixedWindowLimiter(b, 1, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow("key_1", now)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkAllowSameKeyFixedWindowParallel(b *testing.B) {
	limiter := newFixedWindowLimiter(b, 1, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := limiter.Allow("key_1", now)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

}

func BenchmarkAllowManyKeyFixedWindow(b *testing.B) {
	limiter := newFixedWindowLimiter(b, 10, time.Nanosecond)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow(keys[i%(len(keys))], now)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAllowManyKeyFixedWindowParallel(b *testing.B) {
	limiter := newFixedWindowLimiter(b, 1, time.Nanosecond)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, err := limiter.Allow(keys[i%len(keys)], now)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})

}

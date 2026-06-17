package ratelimiter

import (
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

	lim, err := newTestMemoryLimiter(fw)
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

	lim, err := newTestMemoryLimiter(fw)
	if err != nil {
		b.Fatalf("limiter initalization failed %v", err)
	}

	return lim
}

func BenchmarkAllowSameKeyMemoryFixedWindow(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 1, time.Minute)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow("key_1", now)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkAllowSameKeyMemoryFixedWindowParallel(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 1, time.Minute)
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

func BenchmarkAllowManyKeyMemoryFixedWindow(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 10, time.Nanosecond)
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

func BenchmarkAllowManyKeyMemoryFixedWindowParallel(b *testing.B) {
	limiter := newMemoryFixedWindowLimiter(b, 1, time.Nanosecond)
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

func BenchmarkAllowSameKeyMemoryTokenBucket(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 1, 1)
	now := time.Date(1, time.January, 1, 1, 1, 1, 1, time.UTC)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := limiter.Allow("key_1", now)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkAllowSameKeyMemoryTokenBucketParallel(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 1, 1)
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

func BenchmarkAllowManyKeyMemoryTokenBucket(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 10, 1)
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

func BenchmarkAllowManyKeyMemoryTokenBucketParallel(b *testing.B) {
	limiter := newMemoryTokenBucketLimiter(b, 1, 1)
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

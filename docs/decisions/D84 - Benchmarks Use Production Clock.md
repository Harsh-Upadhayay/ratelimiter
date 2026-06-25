# D84 - Benchmarks Use Production Clock

Back to [[V9 HTTP Middleware Index]].

## Context

After [[D83 - MemoryLimiter Implements Limiter via Private Clock]], `Allow` reads time from
`ml.clock.now()` on every call. The test clock (`*testClock`) is a single shared struct guarded
by a mutex (it must be — `TestAllowConcurrentSameKey...` and the parallel benchmarks read it from
many goroutines; see [[G29 - Race Detector]]).

When benchmarks were routed through the test-clock helper, **every `Allow` took that one mutex**.

## Problem

The many-key parallel benchmarks exist to prove the [[D48 ...]] / [[V7 Sharded MemoryStore Index]]
result: different keys hit different shards, so contention scales out. That story depends on there
being **no global serialization point** in the hot path.

A shared test-clock mutex *is* a global serialization point. It flattened `ManyKeyParallel` toward
`SameKeyParallel` and made the sharding win disappear in the numbers — a measurement artifact, not
a production effect (`realClock` is stateless and lock-free).

## Decision

Benchmarks build limiters through the production constructor, which defaults to `realClock{}`:

```go
lim, err := NewMemoryLimiter(fw, NewMemoryStore())  // realClock by default
```

They do **not** route through `newTestMemoryLimiter`. Benchmarks never advance time, so they have
no reason to hold the controllable test clock — and every reason to measure the real, lock-free
production path.

## Why

A benchmark must measure the production code path. Injecting a test-only synchronization primitive
into the thing under measurement corrupts the result. The fix restored the real sharding gap
(many-key serial ~190 ns/op vs the flattened parallel numbers now reflecting true store contention
only).

## Consequence

Removing the accidental global lock exposed real single-key CAS contention — see
[[D85 - Hot Key CAS Exhaustion Is Expected]].

## Links

- [[D83 - MemoryLimiter Implements Limiter via Private Clock]]
- [[D85 - Hot Key CAS Exhaustion Is Expected]]
- [[G28 - Go Benchmarks]]
- [[G29 - Race Detector]]

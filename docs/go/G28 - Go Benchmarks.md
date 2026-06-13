# G28 - Go Benchmarks

Back to [[Rate Limiter Learning Map]].

## Concept

Go benchmarks are functions named `Benchmark...` in `_test.go` files.

The benchmark loop runs `b.N` iterations. The Go tool chooses `b.N` large enough to get a stable measurement.

## Practices used

- Put setup before `b.ResetTimer`.
- Avoid allocations in the measured loop when they are not the subject of the benchmark.
- Precompute key strings before timing.
- Use `b.RunParallel` to benchmark concurrent callers.
- Use a local counter inside each `RunParallel` worker when a shared atomic counter would add measurement noise.
- Use `b.Helper` only for helpers that report benchmark failures.

## Rate limiter use

The benchmarks compared same-key and many-key paths to show the effect of a global mutex.

## Links

- [[D48 - Benchmark Before Storage Refactor]]
- [[G29 - Race Detector]]

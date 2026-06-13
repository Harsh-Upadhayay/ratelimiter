# D48 - Benchmark Before Storage Refactor

Back to [[Rate Limiter Learning Map]].

## Context

Before changing the storage design, we added concurrency tests and benchmarks around the existing limiter.

## Decision

Use benchmarks to expose the current concurrency bottleneck before replacing the limiter-owned map and mutex.

## What it showed

Same-key and many-key parallel benchmarks both serialize behind one limiter-level mutex.

That means unrelated keys cannot progress independently in the current design.

## Why

This keeps the storage refactor motivated by observed pressure instead of jumping directly to the target architecture.

## Tradeoff

Benchmarks add some test code and interpretation overhead, but they give a concrete baseline before changing concurrency control.

## Links

- [[D04 - One Global Mutex]]
- [[D45 - Split Storage from Limiter]]
- [[G28 - Go Benchmarks]]
- [[G29 - Race Detector]]

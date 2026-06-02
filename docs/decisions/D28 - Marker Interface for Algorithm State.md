# D28 - Marker Interface for Algorithm State

Back to [[Rate Limiter Learning Map]].

## Context

The limiter must store per-key state without knowing whether that state belongs to fixed window or token bucket.

Using `any` would allow arbitrary values and weaken the boundary. A shared state struct would accumulate fields that only apply to some algorithms.

## Decision

Use a marker interface for algorithm state.

## Conceptual shape

```go
type algorithmState interface {
	algorithmState()
}
```

Each algorithm-specific state type implements the no-op marker method.

## Why

The limiter can store opaque state while the package controls which state types are valid.

## Tradeoff

This adds syntax and still requires each algorithm to assert its expected concrete state type internally.

## Links

- [[D27 - Algorithm Owned State]]
- [[G19 - Marker Interfaces and Opaque State]]
- [[G20 - Type Assertions]]

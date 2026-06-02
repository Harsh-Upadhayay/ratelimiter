# D31 - Algorithm Owns Config Validation

Back to [[Rate Limiter Learning Map]].

## Context

Before V4, `Limiter` owned fixed-window configuration and validated request limit and window duration.

After introducing multiple algorithms, fixed window and token bucket have different configuration invariants.

## Decision

Each algorithm constructor validates its own configuration.

## Why

The component that owns a configuration invariant should validate it at construction time.

## Current ownership

- `NewFixedWindow` validates request limit and window duration.
- Future `NewTokenBucket` validates capacity and refill rate.
- `NewLimiter` validates only limiter-level concerns, such as a missing algorithm.

## Tradeoff

Callers perform two construction steps, but generic limiter code stays independent of algorithm-specific settings.

## Links

- [[D26 - Introduce Algorithm Interface]]
- [[D32 - Manual Assembly with Built In Algorithms]]
- [[G21 - Constructor Validation Ownership]]

# D33 - Decide Returns State Errors

Back to [[Rate Limiter Learning Map]].

## Context

Algorithms receive opaque state through `algorithmState` and must assert their expected concrete state type.

## Decision

Algorithm decision methods return an error in addition to result and updated state.

## Why

If an existing key contains state for the wrong algorithm, the mismatch should be reported deliberately instead of panicking.

## Tradeoff

The method signature becomes larger, and `Limiter.Allow` must propagate algorithm errors before storing updated state.

## Links

- [[D28 - Marker Interface for Algorithm State]]
- [[D34 - Initialize Missing State Before Assertion]]
- [[G20 - Type Assertions]]

# D17 - Constructor Returns Pointer

Back to [[Rate Limiter Learning Map]].

## Context

The limiter is stateful and contains a `sync.Mutex`.

## Decision

`NewLimiter` should return `*Limiter`.

## Why

Callers should share one limiter instance so all goroutines coordinate through the same mutex and mutate the same state map intentionally.

## Tradeoff

Pointer sharing makes mutation explicit. The implementation must maintain clear locking rules.

## Revisit when

If the design later separates immutable configuration from mutable runtime state.

## Links

- [[D16 - Limiter Owns Runtime State]]
- [[G07 - Do Not Copy Mutexes]]
- [[G04 - Mutexes and Critical Sections]]

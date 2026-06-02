# D29 - Exists Flag for Missing State

Back to [[Rate Limiter Learning Map]].

## Context

Algorithms must distinguish a new key from an existing key with stored state.

## Decision

For now, keep an explicit `exists` boolean in the algorithm decision method.

## Why

This matches the current fixed-window helper and avoids subtle nil-interface behavior.

## Tradeoff

The method signature is slightly awkward. Revisit if interface usage creates pressure for an `InitialState` method or another representation.

## Alternatives deferred

- Separate `InitialState(now)` method: cleaner lifecycle split, larger interface.
- Nil state means missing: fewer parameters, but nil interface handling is subtle in Go.

## Links

- [[D23 - Private Fixed Window Decision Helper]]
- [[D26 - Introduce Algorithm Interface]]
- [[G19 - Marker Interfaces and Opaque State]]

# D25 - Allow as Orchestrator

Back to [[Rate Limiter Learning Map]].

## Context

After V3, fixed-window decision logic lives in a private helper. The public method still owns validation, locking, map access, and persistence of updated state.

## Decision

Keep `Allow` as the orchestration layer.

## Responsibilities

- validate public input
- acquire and release the mutex
- read current state from the map
- call the deterministic decision helper
- write the returned state back to the map
- return the decision result and error

## Why

This creates a local version of the storage/algorithm boundary before introducing interfaces.

## Tradeoff

`Allow` still knows about both storage and algorithm details. That is acceptable until a second store or second algorithm creates real pressure for a stronger boundary.

## Links

- [[D23 - Private Fixed Window Decision Helper]]
- [[D13 - Delay Interfaces]]
- [[G14 - Pure Helper Functions]]

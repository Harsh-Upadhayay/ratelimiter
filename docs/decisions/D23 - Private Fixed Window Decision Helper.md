# D23 - Private Fixed Window Decision Helper

Back to [[Rate Limiter Learning Map]].

## Context

`Allow` currently validates input, locks state, reads the map, performs fixed-window decisions, mutates state, and builds the result.

## Decision

Introduce a private deterministic helper for fixed-window decision logic.

## Shape

The helper should receive the current time, limit, window duration, current state, and whether state exists.

It should return the decision result and the state to store after the decision.

The current parameter order groups inputs by category:

- time
- configuration
- state
- state-existence metadata

## Why

This separates algorithm behavior from map access and locking without introducing public interfaces too early.

## Tradeoff

The helper needs an `exists` input, which is slightly awkward but keeps map lookup outside the algorithm logic.

## Rejected-state rule

When a request is rejected, the returned state should be unchanged because rejected requests do not consume quota or move the window.

## Links

- [[D13 - Delay Interfaces]]
- [[D21 - Result Contract]]
- [[D25 - Allow as Orchestrator]]
- [[G14 - Pure Helper Functions]]
- [[G16 - Helper Parameter Ordering]]

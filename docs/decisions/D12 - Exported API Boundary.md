# D12 - Exported API Boundary

Back to [[Rate Limiter Learning Map]].

## Context

Go uses capitalization for package visibility. Exported names are part of the package API; unexported names are internal implementation detail.

## Decision

For V1, export the main limiter API and hide internal state.

## Shape

- Export `Limiter`
- Export `NewLimiter`
- Export `Allow`
- Keep per-key state unexported as `userState`

## Why

Callers need behavior, not storage details. Hiding `userState` leaves room to change the internal representation when algorithms and stores evolve.

## Tradeoff

Tests outside the package cannot directly inspect `userState`. They should test behavior through the public API.

## Links

- [[D13 - Delay Interfaces]]
- [[G01 - Structs for Grouped State]]

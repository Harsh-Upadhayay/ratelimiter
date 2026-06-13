# D52 - Default Memory Store Constructor

Back to [[Rate Limiter Learning Map]].

## Context

After introducing `StateStore`, one option was to require every caller and test to pass a store explicitly.

## Decision

Keep `NewLimiter(algo)` as the public constructor for now and create a default `MemoryStore` internally.

## Why

This keeps the public API and existing tests simple while still moving the internal implementation toward the storage boundary.

## Tradeoff

Store injection is deferred. A future Redis or alternate memory store will need a new constructor or options pattern, such as `NewLimiterWithStore`.

## Links

- [[D13 - Delay Interfaces]]
- [[D45 - Split Storage from Limiter]]
- [[D49 - MemoryStore Owns Runtime State]]
- [[G21 - Constructor Validation Ownership]]

# D59 - Store Injection Opened

Back to [[V7 Sharded MemoryStore Index]].

## Context

D52 deferred store injection and kept `NewLimiter(algo)` creating a default `MemoryStore` internally, to keep tests simple while the storage boundary was being established.

With `ShardedMemoryStore` now available as an alternative backend, the constructor was opened up.

## Decision

`NewLimiter(algo algorithm, store StateStore)` — the store is now an explicit parameter. Callers and tests choose the backend.

## Why

- Tests need a `newTestLimiter` helper anyway; passing the store there is no extra burden.
- Hardcoding a default backend inside the constructor would hide the choice and make it harder to swap.
- The `StateStore` interface already exists — store injection is the natural completion of that abstraction.

## Tradeoff

Callers must now construct a store explicitly. For a library, a convenience constructor (`NewLimiterWithDefaults`) could be added later without breaking this signature.

## Links

- [[D52 - Default Memory Store Constructor]]
- [[D45 - Split Storage from Limiter]]
- [[D57 - Sharded MemoryStore as Default Backend]]

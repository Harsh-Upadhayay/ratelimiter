# D62 - Redis as Limiter Not Store

Back to [[V7 Sharded MemoryStore Index]].

## Context

The natural instinct was to implement `RedisStore` as a `StateStore` — plugging Redis into the existing `Get`/`CompareAndSwap` interface so `Limiter` could use it transparently.

## Problem

`StateStore` operates on `algorithmState` — a Go interface holding concrete Go structs. Redis stores bytes. Bridging that gap requires serialization, which in turn requires type knowledge. Three options were considered:

- **A**: Store knows all algorithm types (coupling).
- **B**: Store holds a reference to the algorithm for marshal/unmarshal (coupling differently).
- **C**: Change `StateStore` to `[]byte`, serialize in `Limiter` — breaks in-process stores that pay unnecessary serialization cost.

All three options force a tradeoff that reveals a deeper issue: `StateStore` is an in-process abstraction. Remote stores have fundamentally different needs.

## Decision

`RedisStore` will not implement `StateStore`. Instead, a `RedisLimiter` will implement the `Allow(key, now)` contract directly, with algorithm logic living inside Redis Lua scripts. Redis atomicity replaces the CAS retry loop entirely — no `Get → Decide → CAS` cycle needed.

## Implications

- `Limiter` remains in-process only (`MemoryStore`, `ShardedMemoryStore`).
- `RedisLimiter` is a parallel implementation, not a backend swap.
- Both expose the same external behavior to callers.
- Lua scripts own the algorithm logic for the Redis path — each algorithm needs its own script.

## Tradeoff

Algorithm logic is now duplicated: once in Go (`FixedWindow`, `TokenBucket`), once in Lua. This is the cost of true distributed atomicity.

## Links

- [[D45 - Split Storage from Limiter]]
- [[D47 - StateStore Uses Get and CAS]]
- [[D51 - Bounded CAS Retry Loop]]
- [[D54 - Sharded Store Keeps StateStore Contract]]

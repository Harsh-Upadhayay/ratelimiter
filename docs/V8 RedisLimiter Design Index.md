# V8 RedisLimiter Design Index

Back to [[Rate Limiter Learning Map]].

This hub tracks V8: **breadth over the distributed axis**. The single-node depth axis (sharding,
contention, testing) closed with V7. V8 builds the distributed path — a parallel `RedisLimiter`
that implements `Allow` directly via Lua scripts, **not** a swappable `StateStore` backend.

## Pivot and direction

- [[decisions/D62 - Redis as Limiter Not Store]] — the pivot: Redis is a parallel limiter, not a store.
- Direction chosen 2026-06-14: **Redis-first** (port a simple algorithm to learn `EVAL` mechanics)
  over sliding-window-first.

## Decisions

- [[decisions/D63 - Fixed Window in Redis via Lua]]
- [[decisions/D64 - RedisLimiter Algorithm Interface]] — `RedisLimiter` holds a `Script()`/`Args()`
  algorithm interface; concrete `RedisFixedWindow` seals the typed params. Both decoupled and type-safe.

## Redis and Lua concepts

- [[redis/R01 - Key-Value Store and TTL]]
- [[redis/R02 - Atomic Operations with INCR]]
- [[redis/R03 - Lua Scripting for Atomicity]]
- [[redis/R04 - redis.call() and Command Execution]]
- [[redis/R05 - KEYS and ARGV Parameters]]
- [[redis/R06 - SET with NX and EX Flags]]
- [[redis/R07 - Levels of Atomicity in Redis]]
- Hub: [[Redis Concepts Index]]

## Status

Design and concepts logged. **No `RedisLimiter` code written yet** — next move is sketching the
`redisAlgorithm` interface (`Script()`/`Args()`, D64) and the Fixed Window Lua script behind an
`Allow(key)` contract. The Redis path drops the caller-supplied `now`: Redis owns the clock via
`TIME` (no cross-machine clock trust — [[redis/R07 - Levels of Atomicity in Redis]]).

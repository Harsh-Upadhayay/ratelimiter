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
- [[decisions/D64 - RedisLimiter Algorithm Interface]] — `RedisLimiter` holds a `script()`/`args()`
  algorithm interface; concrete `RedisFixedWindow` seals the typed params. Both decoupled and type-safe.
- [[decisions/D65 - Redis Adapter Returns Raw Result]] — Redis adapter returns `any`; `RedisLimiter.Allow`
  parses the script contract into `Result`.
- [[decisions/D66 - Redis Client Ownership Boundary]] — public constructor takes an existing Redis
  client; private adapter remains an internal seam.
- [[decisions/D67 - Sealed Redis Algorithms]] — Redis algorithm extension stays private; users use
  package-provided Redis algorithms for now.
- [[decisions/D68 - Shared Fixed Window Config]] — share stable fixed-window policy config while
  keeping in-memory and Redis execution separate.
- [[decisions/D69 - Backend Qualified Naming]] — backend-specific names use qualifier first:
  `MemoryLimiter`, `RedisLimiter`, `MemoryFixedWindow`, `RedisFixedWindow`.
- [[decisions/D70 - Redis Token Bucket Scaled Integer State]] — Redis Token Bucket stores scaled
  integer token units in a hash to preserve partial refill progress without float state.
- [[decisions/D71 - Redis Token Bucket Time and TTL Policy]] — Redis `TIME` owns the clock; TTL is
  full-refill idle cleanup and is refreshed on every script call.
- [[decisions/D72 - Redis Token Bucket Result Contract]] — Redis algorithms keep returning
  `{allowed, remaining, retryAfterSeconds}`; token bucket floors remaining and ceilings retry-after.
- [[decisions/D73 - Redis Token Bucket Scaling Boundary]] — scaling stays inside `RedisTokenBucket`
  and its Lua script, not in `RedisLimiter` or the public API.

## Redis and Lua concepts

- [[redis/R01 - Key-Value Store and TTL]]
- [[redis/R02 - Atomic Operations with INCR]]
- [[redis/R03 - Lua Scripting for Atomicity]]
- [[redis/R04 - redis.call() and Command Execution]]
- [[redis/R05 - KEYS and ARGV Parameters]]
- [[redis/R06 - SET with NX and EX Flags]]
- [[redis/R07 - Levels of Atomicity in Redis]]
- [[redis/R08 - Redis Hashes and HMGET]]
- [[redis/R09 - Redis TIME and Unit Conversion]]
- [[redis/R10 - Lua Scripts Embedded in Go Strings]]
- [[redis/R11 - Redis Token Bucket Arithmetic]]
- Hub: [[Redis Concepts Index]]

## Status

Design and concepts logged. Redis skeleton code exists for `RedisLimiter`, `RedisFixedWindow`,
`RedisTokenBucket`, and `goRedisAdapter`. Redis client ownership and sealed Redis algorithms are
intentionally decided in [[decisions/D66 - Redis Client Ownership Boundary]] and
[[decisions/D67 - Sealed Redis Algorithms]]. Backend-qualified naming is captured in
[[decisions/D69 - Backend Qualified Naming]]. Redis Token Bucket uses scaled integer units in a hash
([[decisions/D70 - Redis Token Bucket Scaled Integer State]]) and Redis `TIME` for server-owned clock
state ([[decisions/D71 - Redis Token Bucket Time and TTL Policy]]). The arithmetic note for reviewing
integer scaling, refill math, and float tradeoffs is [[redis/R11 - Redis Token Bucket Arithmetic]].

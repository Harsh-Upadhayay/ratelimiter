# D63 - Fixed Window in Redis via Lua

Back to [[V8 RedisLimiter Design Index]].

## Context

Starting V8 (distributed breadth, **Redis-first** over sliding-window-first). The first port is
Fixed Window into a `RedisLimiter`. The decision needs read-modify-write atomicity across
processes — the distributed analog of [[D46 - GetSet Race]] / [[D51 - Bounded CAS Retry Loop]].

## Options for atomicity

See [[R07 - Levels of Atomicity in Redis]] for the full ranking.

- **Primitives** — `SET key N NX EX window` (init-if-absent + TTL, one atomic command) then `DECR`
  (atomic, returns the new value to gate on).
- **MULTI/EXEC** — can't branch on an intermediate read; unusable for "if over limit, reject".
- **WATCH + MULTI/EXEC** — optimistic CAS with a client-side retry loop (Redis twin of D51).
- **Lua (EVAL)** — whole read-decide-write runs atomically server-side; can branch.

## Key analysis: does Fixed Window even need Lua?

- **Count-down** (`SET …NX EX` + `DECR`) is *already primitive-safe*: only the first `SET NX` wins,
  `DECR` is the atomic admission gate, and TTL is set atomically at creation. No over-admission.
- **Count-up** (`INCR` then `EXPIRE`) has a crash race — if the process dies between `INCR` and
  `EXPIRE`, the key never expires → permanent lockout. *That* version needs the two bundled.

So Fixed Window can be done with primitives alone. Lua is not strictly required here.

## Decision

Implement Fixed Window in **Lua anyway**, because:

1. **Cheap place to learn `EVAL`** (KEYS/ARGV, `redis.call`, return values) before Token Bucket
   *forces* Lua — TB must read tokens + last-refill timestamp, compute refill, conditionally
   deduct, write back: no single command does that and MULTI/EXEC can't branch.
2. **Uniform path** — all `RedisLimiter` algorithms express their logic the same way.
3. **Pessimistic single-shot** atomicity replaces the CAS retry loop entirely ([[D62 - Redis as Limiter Not Store]]);
   the Redis path has no `ErrCASConflict`, no retry.

## Tradeoff

- Algorithm logic is duplicated — once in Go (`FixedWindow`), once in Lua (the D62 cost).
- Lua scripts are less unit-testable from Go than the pure `Decide` functions; they need a real
  (or mocked) Redis to exercise.

## Links

- [[D62 - Redis as Limiter Not Store]]
- [[D46 - GetSet Race]]
- [[D51 - Bounded CAS Retry Loop]]
- [[R03 - Lua Scripting for Atomicity]]
- [[R06 - SET with NX and EX Flags]]
- [[R07 - Levels of Atomicity in Redis]]

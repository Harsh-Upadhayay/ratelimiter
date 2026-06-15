# D64 - RedisLimiter Algorithm Interface

Back to [[V8 RedisLimiter Design Index]].

## Context

`RedisLimiter` is the parallel, in-Redis path ([[D62 - Redis as Limiter Not Store]]). It needs
to run *some* algorithm (Fixed Window first, Token Bucket later) whose logic lives in a Lua
script ([[D63 - Fixed Window in Redis via Lua]]). The question: how does `RedisLimiter` stay
generic across algorithms without coupling to any one algorithm's parameters?

Each algorithm has a *different* parameter shape:

- Fixed Window → `limit`, `window`
- Token Bucket → `limit`, `rate`

## Options

- **Named fields on `RedisLimiter`** (`limit int`, `window time.Duration`) — type-safe, but
  hard-couples the generic limiter to Fixed Window's shape. Breaks for Token Bucket.
- **`args []string` on `RedisLimiter`** — generic, but loses compile-time type safety. A wrong
  value only fails at runtime inside Lua.
- **Per-algorithm limiter types** (`RedisFixedWindowLimiter`, `RedisTokenBucketLimiter`) — type-safe
  and decoupled, but duplicates the EVAL plumbing per algorithm; no single generic limiter.
- **Algorithm interface** (chosen) — mirrors the in-process `algorithm` split.

## Decision

`RedisLimiter` holds a Redis-side **algorithm interface**, not algorithm params. The interface
exposes what the limiter needs to call `EVAL`:

```go
type redisAlgorithm interface {
    Script() string   // the Lua source
    Args() []string   // algorithm-specific ARGV, converted from typed fields
}
```

The concrete `RedisFixedWindow` holds typed fields (`limit int`, `window time.Duration`) and
implements `Args()` by converting them to `[]string`. `RedisLimiter.Allow(key)` calls
`algo.Script()` and `algo.Args()`, passing `key` as `KEYS[1]` and the args as `ARGV` — never
knowing which algorithm it runs.

This is the same decoupling the in-process side already uses: `Limiter` holds `algorithm` and
knows nothing of `limit`/`window`; `FixedWindow` seals those inside the concrete type. The Redis
twin swaps the interface method from `Decide(now, state)` (compute in Go) to `Script()`/`Args()`
(compute server-side in Lua) — because on the Redis path the *script* is the algorithm.

## Why this beats the false dichotomy

The original tension was "coupling vs `[]string`". The resolution: the `[]string` was never the
enemy — letting it leak into the *generic* `RedisLimiter` was. By sealing the conversion inside
`RedisFixedWindow.Args()`:

- **Type safety** lives at the construction site — `RedisFixedWindow{limit: 100, window: time.Minute}`
  is compiler-checked.
- The unavoidable `[]string` conversion sits at the one boundary where it is forced anyway (Lua
  only speaks strings — [[R05 - KEYS and ARGV Parameters]]).
- `RedisLimiter` stays fully generic.

You get **both** decoupling and type safety.

## Tradeoff

- One more interface + concrete type per algorithm vs. just a struct literal. Worth it for the
  decoupling and for parity with the in-process design.
- The interface is intentionally minimal (`Script`/`Args`); if later algorithms need more KEYS
  than one, the contract may have to grow (e.g. `Keys()`), revisit then.

## Links

- [[D62 - Redis as Limiter Not Store]]
- [[D63 - Fixed Window in Redis via Lua]]
- [[R05 - KEYS and ARGV Parameters]]
- [[R03 - Lua Scripting for Atomicity]]

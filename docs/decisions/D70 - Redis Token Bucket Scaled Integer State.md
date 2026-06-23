# D70 - Redis Token Bucket Scaled Integer State

Back to [[V8 RedisLimiter Design Index]].

## Context

Redis Token Bucket needs to preserve fractional refill progress.

Example:

```text
refillRate = 2 tokens/sec
elapsed = 250ms
refill = 0.5 token
```

If Redis stores only whole tokens, that `0.5` progress is lost. If Redis stores floats, the script
and the stored state depend on floating-point representation.

## Options

- **Whole-token integers** - simple, but loses partial refill progress.
- **Float strings** - keeps precision at the domain level, but stores/parses floating values and
  spreads float behavior into Redis state.
- **Scaled integers** - represent one token as a fixed number of integer units.

## Decision

Use scaled integer token state for Redis Token Bucket.

```text
tokenScale = 1000

1 token   = 1000 units
0.5 token = 500 units
```

Redis stores:

```text
key -> {
  tokens: "...",
  last_refill_ms: "..."
}
```

`tokens` is stored in scaled units. One request costs `tokenScale` units.

## Why

Scaled integers preserve sub-token progress without storing floats in Redis.

The important check becomes:

```lua
tokens >= tokenScale
```

That means:

```text
Does this bucket have at least one whole request token available?
```

## Tradeoffs

- **Memory:** slightly more than a single string key because the token bucket uses a Redis hash with
  two fields. Still bounded to one key per active rate-limit key.
- **Latency:** `HMGET`/`HSET` are a little more structured than `GET`/`SET`, but the network round
  trip and Lua execution dominate.
- **Concurrency:** Lua owns the read-refill-decide-write sequence atomically, so scaled arithmetic
  does not introduce client-side races.
- **Correctness:** better than whole-token storage because partial refill progress is preserved.
- **Complexity:** higher because the code must track units: scaled token units, milliseconds, and
  seconds.

## Links

- [[D71 - Redis Token Bucket Time and TTL Policy]]
- [[D72 - Redis Token Bucket Result Contract]]
- [[D73 - Redis Token Bucket Scaling Boundary]]
- [[redis/R08 - Redis Hashes and HMGET]]
- [[redis/R09 - Redis TIME and Unit Conversion]]
- [[redis/R11 - Redis Token Bucket Arithmetic]]
- [[G40 - Unexported Package Constants]]

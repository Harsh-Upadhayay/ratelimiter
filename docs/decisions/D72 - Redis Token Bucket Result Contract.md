# D72 - Redis Token Bucket Result Contract

Back to [[V8 RedisLimiter Design Index]].

## Context

`RedisLimiter.Allow` parses Redis Lua output into the public Go `Result`.

The existing Redis fixed-window script returns:

```text
allowed, remaining, retryAfterSeconds
```

Redis Token Bucket could return more precise values because refill is continuous.

## Options

- **Return retry-after seconds** - keeps the Redis result contract uniform and simple.
- **Return retry-after milliseconds** - more precise for token bucket, but changes the current Redis
  contract.
- **Algorithm-specific output parsing** - each Redis algorithm can return its own shape, but
  `RedisLimiter` needs more dispatch/parsing complexity.

## Decision

Keep the Redis Lua result contract:

```text
allowed, remaining, retryAfterSeconds
```

For Redis Token Bucket:

- `remaining` is rounded down to whole tokens:

```lua
math.floor(tokens / tokenScale)
```

- `retryAfterSeconds` is rounded up:

```lua
math.ceil(deficit / refillRate)
```

## Why

The current learning value is Redis/Lua atomic state transition, not API precision tuning.

Rounding down `remaining` preserves the public `Result.Remaining int` contract and hides scaled
internal units. Rounding up `RetryAfter` is conservative: returning `0` for a partial-second wait
would invite an immediate retry when a full token is not available yet.

## Tradeoffs

- **Memory:** no difference.
- **Latency:** no difference.
- **Concurrency:** no difference; this is an API contract decision.
- **Client behavior:** seconds are coarse for high-rate token buckets. Milliseconds can be revisited
  if caller retry precision becomes important.
- **Maintainability:** uniform Redis result parsing stays simple.

## Links

- [[D65 - Redis Adapter Returns Raw Result]]
- [[D70 - Redis Token Bucket Scaled Integer State]]
- [[G13 - Exported Result Structs]]

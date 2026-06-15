# D65 - Redis Adapter Returns Raw Result

Back to [[V8 RedisLimiter Design Index]].

## Context

V8 introduces a Redis-backed limiter path. `RedisLimiter` runs a Lua script through a small
`redisAdapter` boundary, then returns the public `Result`.

The fixed-window Lua contract currently returns:

```text
allowed, remaining, retryAfterSeconds
```

The design question: should the adapter know that script output shape, or should it only execute
the script and return the raw Redis result?

## Options

- **Adapter returns `Result`** — simplest call site, but turns the Redis client adapter into
  rate-limiter domain logic.
- **Adapter returns `[]interface{}`** — validates that scripts return arrays, but still leaks
  script result shape into a boundary whose job is just Redis I/O.
- **Adapter returns `any`** — keeps the adapter generic; `RedisLimiter.Allow` owns parsing the
  limiter-specific Lua contract.
- **Algorithm parses output** — `RedisFixedWindow` could parse its own script result, but this is
  another interface method before there is enough variation to justify it.

## Decision

Make `redisAdapter.eval` return raw `any`.

`goRedisAdapter` should only translate the third-party Redis client call:

```text
Eval(ctx, script, keys, args) -> raw Redis result or error
```

`RedisLimiter.Allow` should interpret the script result and convert it into:

```go
Result{Allowed, Remaining, RetryAfter}
```

## Why this is the right boundary

The adapter is infrastructure plumbing. It should know how to call Redis, not what a rate-limiter
Lua script means.

`RedisLimiter.Allow` is the orchestrator for the Redis path. It already knows:

- the key being checked,
- which Redis algorithm is being run,
- the public `Result` contract,
- which output shape the Lua limiter script promised.

That makes `RedisLimiter.Allow` the natural place for type assertions and conversion from Redis
values into Go domain values.

## Tradeoffs

- **Memory:** no meaningful difference. Returning `any` does not add state; it only preserves the
  raw value returned by the Redis client.
- **Latency:** no meaningful difference. Type assertions are tiny compared with a Redis network
  round trip.
- **Concurrency:** no direct difference. This decision changes interpretation boundaries, not
  locking or atomicity. Redis script atomicity remains the concurrency guarantee.
- **Maintainability:** better. The adapter stays reusable for scripts with different return
  shapes.
- **Testability:** better for unit tests. Tests can mock `redisAdapter.eval` with raw Redis-like
  values and verify that `RedisLimiter.Allow` parses them correctly.
- **Cost:** `RedisLimiter.Allow` must do explicit shape validation. That is appropriate because
  the script contract is part of the limiter behavior, not the Redis transport behavior.

## Deferred

If future Redis algorithms produce meaningfully different output shapes, revisit an algorithm-owned
parser method, such as:

```go
parseResult(raw any) (Result, error)
```

Do not add that method yet. Fixed Window and the planned Token Bucket can both reasonably return
the same `{allowed, remaining, retryAfterSeconds}` shape.

## Links

- [[D62 - Redis as Limiter Not Store]]
- [[D63 - Fixed Window in Redis via Lua]]
- [[D64 - RedisLimiter Algorithm Interface]]
- [[G20 - Type Assertions]]
- [[redis/R03 - Lua Scripting for Atomicity]]

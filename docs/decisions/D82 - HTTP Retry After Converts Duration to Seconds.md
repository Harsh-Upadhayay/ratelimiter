# D82 - HTTP Retry After Converts Duration to Seconds

Back to [[V9 HTTP Middleware Index]].

## Context

Every limiter result uses the Go domain contract:

```go
Result.RetryAfter time.Duration
```

HTTP `Retry-After` uses seconds when represented as a delta value.

## Decision

Keep `Result.RetryAfter` as `time.Duration` everywhere inside Go.

Convert it to ceiling whole seconds only at the HTTP middleware boundary.

```text
Result.RetryAfter -> Retry-After header seconds
```

## Cross-limiter validation

- `MemoryFixedWindow` returns a `time.Duration` until the current window resets.
- `MemoryTokenBucket` returns a `time.Duration` until one full token is available.
- `MemoryStore` and `ShardedMemoryStore` do not change the result; they only store state.
- `RedisFixedWindow` Lua returns seconds; `RedisLimiter` converts those seconds into `time.Duration`.
- `RedisTokenBucket` Lua returns seconds; `RedisLimiter` converts those seconds into `time.Duration`.

So the HTTP middleware sees one consistent contract: `time.Duration`.

## Conversion rule

```text
RetryAfter <= 0       -> "0"
0 < RetryAfter <= 1s  -> "1"
1500ms                -> "2"
```

This avoids telling clients to retry immediately when a sub-second wait remains.

## Why

The domain API should not leak HTTP header formatting.

The middleware is the boundary that knows HTTP needs a string header value.

## Links

- [[D21 - Result Contract]]
- [[D77 - Rate Limit HTTP Headers]]
- [[G48 - Ceiling Duration Conversion]]

# C09 - V9 HTTP Middleware Checkpoint

Back to [[V9 HTTP Middleware Index]].

## Status

V9 is structurally complete for now.

The package now has an HTTP middleware boundary that turns a limiter into request-path behavior:

```text
HTTP request -> key function -> Limiter.Allow(ctx, key) -> pass or reject
```

## Current API shape

Because the package is named `ratelimiter`, the public middleware type uses the package qualifier
instead of repeating `RateLimiting` in the type name:

```go
ratelimiter.Middleware
ratelimiter.NewMiddleware(...)
```

Core pieces:

```text
Limiter interface    -> Allow(ctx, key) (Result, error)
Middleware           -> wraps http.Handler
KeyFunc              -> derives a key from *http.Request
FailurePolicy        -> FailOpen or FailClosed
WithFailurePolicy    -> optional functional option
```

## Runtime behavior

```text
empty key                 -> 400 Bad Request
limiter error + FailOpen  -> call next handler silently
limiter error + FailClosed -> 503 Service Unavailable
allowed result            -> set X-RateLimit-Remaining, call next handler
rejected result           -> set X-RateLimit-Remaining, set Retry-After, return 429
Wrap(nil)                 -> panic("ratelimiter: nil next handler")
```

The middleware writes status codes and headers only. It does not write response bodies.

## Retry-After contract

Inside Go, `Result.RetryAfter` remains a `time.Duration` for every limiter path.

At the HTTP boundary, it is converted to ceiling whole seconds for the `Retry-After` header.

## Verified

```text
go test ./...
```

passes at checkpoint time.

## Deferred

Intentionally deferred:

- middleware behavior tests,
- observability/logging/metrics,
- local HTTP example server,
- sliding window algorithms,
- memory adapter for the `Limiter` interface,
- delegate-on-error policy,
- Redis integration/runtime tests.

## Resume options

Pick based on learning goal:

- **Go testing:** add middleware tests with fake `Limiter` and `httptest`.
- **HTTP/system path:** add a local HTTP example server.
- **Redis/data structures:** design Sliding Window Log with sorted sets.
- **API design:** adapt `MemoryLimiter` to the `Limiter` interface through a clock-backed adapter.
- **Production hardening:** observability and delegate-on-error policy.

## Links

- [[D74 - HTTP Middleware Boundary]]
- [[D75 - Middleware Failure Policy]]
- [[D76 - Caller Provided Key Function]]
- [[D77 - Rate Limit HTTP Headers]]
- [[D78 - Functional Options for Middleware]]
- [[D79 - Behavior Named Middleware]]
- [[D80 - Required Dependencies Outside Functional Options]]
- [[D81 - Fail Open Is Silent Pass Through]]
- [[D82 - HTTP Retry After Converts Duration to Seconds]]
- [[G41 - net-http Middleware Pattern]]
- [[G49 - http HandlerFunc Adapter]]

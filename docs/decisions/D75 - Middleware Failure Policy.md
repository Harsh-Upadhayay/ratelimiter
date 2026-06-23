# D75 - Middleware Failure Policy

Back to [[V9 HTTP Middleware Index]].

## Context

When the limiter is in the HTTP request path, limiter infrastructure can fail.

For Redis-backed limiting, this includes Redis being unavailable, timing out, or returning an
operational error.

## Decision

Make limiter failure behavior configurable at the whole middleware level.

Use two policies:

```text
FailOpen   -> limiter error allows the request through
FailClosed -> limiter error returns 503 Service Unavailable
```

Default direction: `FailOpen`, because rate limiting is usually protective middleware rather than
the application's core feature.

Validate the policy in the middleware constructor. If an invalid policy exists, constructing the
middleware should fail.

## Why

`429 Too Many Requests` and `503 Service Unavailable` mean different things:

- `429`: limiter is healthy and says this key is over the limit.
- `503`: limiter infrastructure failed and this middleware chose fail-closed.

Fail-open protects availability. Fail-closed protects enforcement.

## Scope

The policy is per middleware instance, not per individual request.

If different routes need different behavior, create different middleware instances for those route
groups later.

## Tradeoff

- **Availability:** fail-open is best.
- **Abuse protection:** fail-closed is strongest.
- **Operational surprise:** configurable policy makes the caller's business choice explicit.
- **Complexity:** one enum-like config field and constructor validation.

## Links

- [[D74 - HTTP Middleware Boundary]]
- [[G43 - Iota Enum Pattern]]
- [[G21 - Constructor Validation Ownership]]

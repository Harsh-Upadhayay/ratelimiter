# D81 - Fail Open Is Silent Pass Through

Back to [[V9 HTTP Middleware Index]].

## Context

The middleware has to decide what to do when the limiter itself fails.

A question came up: should `FailOpen` notify the downstream handler that rate limiting failed, or
should it simply allow the request?

## Decision

`FailOpen` means silent pass-through.

```text
limiter error + FailOpen -> call next handler as if the request was allowed
```

The downstream handler is not told that the limiter failed.

## Why

This is the common meaning of fail-open: prioritize availability over enforcement.

If the application needs to know that rate limiting failed, that is a different policy, such as:

```text
DelegateOnLimiterError
```

That would require putting error information into `context.Context` or another explicit signal and
letting the downstream handler decide what to do.

## Deferred

Do not implement delegate-on-error now.

Revisit it when the project needs route-specific business decisions such as:

- login must fail closed,
- public feed can fail open,
- expensive AI endpoint wants app-owned fallback logic.

## Tradeoff

- **Availability:** fail-open gives the best availability.
- **Abuse protection:** limiter outages temporarily remove protection.
- **Simplicity:** no extra context values or downstream coupling.
- **Observability:** should come later from logs/metrics, not from the protected handler.

## Links

- [[D75 - Middleware Failure Policy]]
- [[D74 - HTTP Middleware Boundary]]
- [[G38 - context.Context]]

# D77 - Rate Limit HTTP Headers

Back to [[V9 HTTP Middleware Index]].

## Context

The limiter already returns:

```go
Result{Allowed, Remaining, RetryAfter}
```

HTTP callers can use that information.

## Decision

The middleware should write status codes and headers only, without response bodies.

Set:

```text
X-RateLimit-Remaining: <remaining>
```

after every limiter decision.

On rejection, return:

```text
429 Too Many Requests
Retry-After: <seconds>
```

`Retry-After` should be rounded up to seconds so a sub-second wait does not become `0`. The conversion is owned by the HTTP middleware boundary; `Result.RetryAfter` remains a `time.Duration`.

## Why

Headers are useful to clients, while response bodies are application-specific.

The middleware should not impose JSON, plain text, or any other body format on the caller's API.

## Tradeoff

- **Client usefulness:** better than status-only.
- **API neutrality:** preserved by not writing a body.
- **Precision:** `Retry-After` in whole seconds is simple and conservative.
- **Complexity:** small string conversions from `int` and `time.Duration`.

## Links

- [[D21 - Result Contract]]
- [[D82 - HTTP Retry After Converts Duration to Seconds]]
- [[G23 - Duration Seconds Conversion]]
- [[G48 - Ceiling Duration Conversion]]
- [[G45 - HTTP ResponseWriter Headers]]

# D74 - HTTP Middleware Boundary

Back to [[V9 HTTP Middleware Index]].

## Context

After the in-memory limiter and Redis limiter paths, the next learning step is to place the limiter
in an HTTP request path.

The project intentionally skipped more tests and another algorithm for now.

## Decision

Build an HTTP middleware boundary next.

Conceptually:

```text
incoming request
-> extract rate-limit key
-> call limiter
-> pass request or return HTTP status
```

## Why

This teaches new Go and systems concepts without repeating algorithm work:

- `net/http` middleware shape,
- request-scoped `context.Context`,
- request identity extraction,
- HTTP status codes,
- rate-limit response headers,
- infrastructure failure policy.

## Tradeoff

- **Learning value:** high, because it turns the library into a system boundary.
- **Complexity:** moderate; HTTP introduces user-facing behavior and error policy.
- **Correctness:** the middleware must distinguish rate-limit rejection from limiter infrastructure
  failure.

## Links

- [[D19 - Result and Error Return]]
- [[D21 - Result Contract]]
- [[D62 - Redis as Limiter Not Store]]
- [[G41 - net-http Middleware Pattern]]

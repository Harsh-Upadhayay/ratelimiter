# V9 HTTP Middleware Index

Back to [[Rate Limiter Learning Map]].

This hub tracks the move from a rate-limiter library to an HTTP integration boundary.

## Why this phase exists

The in-memory path taught Go state, mutexes, CAS, sharding, and benchmarks. The Redis path taught
server-side atomicity through Lua. The next useful learning step is not another algorithm; it is
putting the limiter into a real request path:

```text
HTTP request -> extract rate-limit key -> call limiter -> status code and headers
```

## Decisions

- [[decisions/D74 - HTTP Middleware Boundary]]
- [[decisions/D75 - Middleware Failure Policy]]
- [[decisions/D76 - Caller Provided Key Function]]
- [[decisions/D77 - Rate Limit HTTP Headers]]
- [[decisions/D78 - Functional Options for Middleware]]
- [[decisions/D79 - Behavior Named Middleware]]

## Go concepts

- [[go/G41 - net-http Middleware Pattern]]
- [[go/G42 - Functional Options]]
- [[go/G43 - Iota Enum Pattern]]
- [[go/G44 - Function Types as Callbacks]]
- [[go/G45 - HTTP ResponseWriter Headers]]

## Current direction

Use a behavior-oriented public name such as `RateLimitingMiddleware`, not `RedisMiddleware`.

The middleware's job is HTTP behavior:

- derive the key,
- call the limiter,
- decide pass-through vs reject,
- write status codes and headers.

The backing limiter may be Redis-backed first, but Redis should not dominate the middleware name.

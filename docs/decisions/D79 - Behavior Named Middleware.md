# D79 - Behavior Named Middleware

Back to [[V9 HTTP Middleware Index]].

## Context

The first middleware discussion used names like `RedisMiddleware` because the backing limiter path
is Redis-first.

The user corrected the direction: the middleware name should describe behavior, not the current
backend.

## Decision

Use behavior-oriented naming:

```text
RateLimitingMiddleware
```

Avoid public middleware names that contain `Redis`.

## Why

The middleware's public behavior is HTTP rate limiting:

- extract key,
- call limiter,
- pass or reject,
- set headers.

Redis is an implementation detail of the limiter dependency, not the purpose of the middleware.

## Important boundary

This naming decision does not force the project to introduce a fully generic limiter interface
immediately.

For the first implementation, the middleware may still depend on the Redis-backed limiter path if
that keeps the step small. The name is allowed to be more behavior-oriented than the first backend
because the HTTP behavior is not Redis-specific.

If that dependency starts to feel awkward, that pressure will be the natural reason to introduce a
small middleware-facing interface later.

## Tradeoff

- **API clarity:** better for users; they are adding rate limiting, not Redis.
- **Backend flexibility:** better later; the name does not need to change if memory or another
  backend is supported.
- **Risk:** a generic name can overpromise backend-agnostic behavior before the abstraction exists.
  Keep docs and constructor shape honest.

## Links

- [[D69 - Backend Qualified Naming]]
- [[D74 - HTTP Middleware Boundary]]
- [[D78 - Functional Options for Middleware]]

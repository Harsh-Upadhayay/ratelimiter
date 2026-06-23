# D78 - Functional Options for Middleware

Back to [[V9 HTTP Middleware Index]].

## Context

The middleware needs configuration:

- key extraction function,
- failure policy.

The user already understands config structs and chose to learn the functional-options pattern.

## Decision

Use functional options for middleware construction.

Conceptually:

```go
NewRateLimitingMiddleware(limiter, WithKeyFunc(fn), WithFailurePolicy(FailOpen))
```

Options should be plain mutators:

```go
type MiddlewareOption func(*middlewareConfig)
```

The constructor applies all options and then validates the final config.

## Required and default values

`KeyFunc` is required. If no key function is provided, constructor returns an error.

`FailurePolicy` defaults to `FailOpen`, but can be changed with an option.

Invalid failure policy values are rejected by constructor validation.

## Why

Functional options teach a common Go API pattern:

- defaults live in the constructor,
- optional settings are readable at the call site,
- future settings can be added without changing the constructor signature.

## Tradeoff

- **Readability:** good at call sites once the pattern is understood.
- **Extensibility:** better than positional constructor parameters.
- **Complexity:** more moving pieces than a config struct.
- **Validation:** final constructor validation is still required for missing required fields and
  invalid enum-like values.

## Links

- [[D75 - Middleware Failure Policy]]
- [[D76 - Caller Provided Key Function]]
- [[G42 - Functional Options]]
- [[G21 - Constructor Validation Ownership]]

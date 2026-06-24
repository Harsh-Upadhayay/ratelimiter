# D78 - Functional Options for Middleware

Back to [[V9 HTTP Middleware Index]].

## Context

The middleware has two required dependencies and one optional policy:

- required limiter behavior,
- required key extraction function,
- optional failure policy.

The user already understands config structs and chose to learn the functional-options pattern, but
we later clarified that functional options should not hide required dependencies.

## Decision

Use functional options only for optional middleware configuration.

Required values stay as explicit constructor parameters:

```go
NewMiddleware(limiter, keyFunc, WithFailurePolicy(FailClosed))
```

Optional settings use named option functions:

```go
type MiddlewareOption func(*Middleware)
```

For this middleware, options can mutate the partially built `*Middleware` directly. A separate
private config struct would duplicate the runtime object because the runtime fields and config fields
are currently the same.

## Required and default values

`Limiter` is required. If it is nil, constructor returns an error.

`KeyFunc` is required. If it is nil, constructor returns an error.

`FailurePolicy` defaults to `FailOpen`, but can be changed with an option.

Invalid failure policy values are rejected by constructor validation after options are applied.

## Why

Functional options are useful for optional configuration and named overrides. Required constructor
dependencies should stay visible at the call site.

This keeps the API honest:

```text
required behavior -> constructor parameters
optional tuning   -> functional options
```

## Tradeoff

- **Readability:** required dependencies are obvious; optional policy remains named.
- **Extensibility:** future optional middleware behavior can use more options.
- **Complexity:** lower than using both `Middleware` and a duplicate `middlewareConfig`.
- **Validation:** constructor still owns final validation after applying options.

## Links

- [[D75 - Middleware Failure Policy]]
- [[D76 - Caller Provided Key Function]]
- [[D80 - Required Dependencies Outside Functional Options]]
- [[G42 - Functional Options]]
- [[G47 - Functional Options on Runtime Structs]]
- [[G21 - Constructor Validation Ownership]]

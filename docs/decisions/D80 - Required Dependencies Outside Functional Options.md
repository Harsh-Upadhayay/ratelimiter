# D80 - Required Dependencies Outside Functional Options

Back to [[V9 HTTP Middleware Index]].

## Context

The first functional-options shape put `KeyFunc` behind `WithKeyFunc`.

That raised a design question: if `KeyFunc` is mandatory, should it be an option at all?

## Decision

Keep required middleware dependencies outside functional options.

Constructor shape:

```go
NewMiddleware(limiter Limiter, keyFunc KeyFunc, opts ...MiddlewareOption)
```

Options are reserved for optional configuration such as `WithFailurePolicy`.

## Why

A caller cannot construct useful middleware without both:

- something that can rate-limit (`Limiter`),
- a way to derive keys from requests (`KeyFunc`).

Hiding required values inside optional-looking `With...` functions makes the API less honest.

## Tradeoff

- **API clarity:** better; required values are visible in the constructor signature.
- **Compile-time guidance:** better; callers see required arguments immediately.
- **Extensibility:** optional settings can still be added through options.
- **Verbosity:** slightly more positional constructor arguments, but both are core dependencies.

## Links

- [[D76 - Caller Provided Key Function]]
- [[D78 - Functional Options for Middleware]]
- [[G42 - Functional Options]]

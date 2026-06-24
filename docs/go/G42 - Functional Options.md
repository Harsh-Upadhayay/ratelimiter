# G42 - Functional Options

Back to [[Go Concepts Index]].

## Concept

Functional options are a Go constructor pattern for optional configuration.

Instead of a long positional constructor:

```go
NewThing(a, b, c, d)
```

the caller passes named option functions for optional overrides:

```go
NewThing(required, WithTimeout(timeout), WithName(name))
```

The usual shape is:

```go
type Option func(*config)
```

An option is just a function that mutates some construction state.

## Required vs optional

Functional options are best for optional settings.

Required dependencies should usually stay as explicit constructor parameters:

```go
NewMiddleware(limiter, keyFunc, WithFailurePolicy(FailClosed))
```

This keeps required behavior visible at the call site.

## In this project

Middleware construction uses direct parameters for:

- `Limiter`,
- `KeyFunc`.

It uses functional options for optional settings:

- `WithFailurePolicy`.

Because the middleware's runtime fields are also its configuration fields, the option type can mutate
`*Middleware` directly instead of a duplicate private config struct.

## Why use it

Benefits:

- readable optional overrides,
- defaults can live in one place,
- adding optional settings does not change the constructor signature.

Costs:

- more indirection than setting a struct field,
- more code than positional parameters,
- final validation still belongs in the constructor.

## Links

- [[D78 - Functional Options for Middleware]]
- [[D80 - Required Dependencies Outside Functional Options]]
- [[G47 - Functional Options on Runtime Structs]]
- [[G21 - Constructor Validation Ownership]]

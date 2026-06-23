# G42 - Functional Options

Back to [[Go Concepts Index]].

## Concept

Functional options are a Go constructor pattern for optional configuration.

Instead of a long positional constructor:

```go
NewThing(a, b, c, d)
```

the caller passes named option functions:

```go
NewThing(required, WithTimeout(timeout), WithName(name))
```

The usual shape is:

```go
type Option func(*config)
```

Each option mutates a private config value during construction.

## In this project

Middleware construction uses functional options to configure:

- key extraction,
- failure policy.

The constructor applies options and then validates the final config.

## Why use it

Benefits:

- readable call sites,
- defaults can live in one place,
- adding new options does not change the constructor signature.

Costs:

- more indirection than a config struct,
- more code than positional parameters,
- validation still belongs in the constructor.

## Links

- [[D78 - Functional Options for Middleware]]
- [[G21 - Constructor Validation Ownership]]

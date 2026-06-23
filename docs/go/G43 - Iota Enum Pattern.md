# G43 - Iota Enum Pattern

Back to [[Go Concepts Index]].

## Concept

Go does not have native enums.

The common pattern is a named integer type plus constants:

```go
type FailurePolicy int

const (
    FailOpen FailurePolicy = iota
    FailClosed
)
```

`iota` starts at `0` in a `const` block and increments for each line.

## Why not bool

This call is unclear:

```go
NewMiddleware(true)
```

This call is readable:

```go
NewMiddleware(WithFailurePolicy(FailOpen))
```

## Validation

Named integer types are not closed sets. A caller can still write:

```go
FailurePolicy(99)
```

So constructors should validate enum-like config values when invalid values would create surprising
behavior.

## Links

- [[D75 - Middleware Failure Policy]]
- [[D78 - Functional Options for Middleware]]

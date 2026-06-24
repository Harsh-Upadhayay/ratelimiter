# G47 - Functional Options on Runtime Structs

Back to [[Go Concepts Index]].

## Concept

Functional options often mutate a private config struct:

```go
type Option func(*config)
```

But they can also mutate the final runtime object when the runtime object is effectively the config.

Example shape:

```go
type MiddlewareOption func(*Middleware)
```

## In this project

`Middleware` currently stores:

- limiter,
- key function,
- failure policy.

Those are also the construction fields. A separate `middlewareConfig` would duplicate the same data.

So using:

```go
type MiddlewareOption func(*Middleware)
```

is acceptable here.

## Safety boundary

The constructor still controls object publication:

```text
build object -> apply options -> validate -> return object
```

Even though options mutate the runtime object, callers never receive it unless validation succeeds.

## Tradeoff

- **Less ceremony:** no duplicate config struct.
- **More direct:** options modify the object they are configuring.
- **Less isolation:** options touch the partially built runtime object.
- **Still safe enough here:** the object is private until constructor returns.

## Links

- [[D78 - Functional Options for Middleware]]
- [[D80 - Required Dependencies Outside Functional Options]]
- [[G42 - Functional Options]]

# G46 - Method Sets and Interface Satisfaction

Back to [[Go Concepts Index]].

## Concept

A Go type satisfies an interface based on its **method set**.

The receiver type controls which method set gets the method.

## Value receiver

If a method has a value receiver:

```go
func (x T) Allow(ctx context.Context, key string) (Result, error)
```

then both `T` and `*T` have that method for interface satisfaction.

So both can satisfy:

```go
type Limiter interface {
    Allow(context.Context, string) (Result, error)
}
```

## Pointer receiver

If a method has a pointer receiver:

```go
func (x *T) Allow(ctx context.Context, key string) (Result, error)
```

then only `*T` satisfies the interface. The value type `T` does not.

This means:

```go
var _ Limiter = &T{} // ok
var _ Limiter = T{}  // does not compile
```

## In this project

`RedisLimiter.Allow` has a pointer receiver, so `*RedisLimiter` satisfies `Limiter`.

The middleware should accept the interface by value:

```go
func NewRateLimitingMiddleware(limiter Limiter, opts ...MiddlewareOption)
```

Passing a `*RedisLimiter` does not copy the limiter object. The interface value stores type
information plus the concrete pointer.

Do not use `*Limiter`. A pointer to an interface is usually a pointer to the interface box, not a
pointer to the concrete limiter behavior.

## Rule of thumb

- Accept interfaces by value.
- Use pointer receivers on concrete mutable/large types.
- Let pointer receiver method sets naturally require pointer values for those concrete types.

## Links

- [[G03 - Pointer Receivers]]
- [[G09 - Method Receivers]]
- [[D79 - Behavior Named Middleware]]

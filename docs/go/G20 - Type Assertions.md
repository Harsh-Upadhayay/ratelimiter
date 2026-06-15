# G20 - Type Assertions

Back to [[Rate Limiter Learning Map]].

## Go practice

When a value is held through an interface, a type assertion retrieves the expected concrete type.

Conceptually:

```go
state, ok := rawState.(fixedWindowState)
```

The same pattern appears at external boundaries:

```go
values, ok := result.([]interface{})
```

This means: "I have a value stored in an interface; try to view it as a `[]interface{}`."

- If the dynamic value really is a `[]interface{}`, `values` receives that slice and `ok` is `true`.
- If it is not, `values` receives the zero value for that type (`nil`) and `ok` is `false`.

This is the safe "comma-ok" form. The one-value form:

```go
values := result.([]interface{})
```

panics if the dynamic type is not `[]interface{}`.

## Why it matters

The limiter treats state as opaque, but each algorithm needs its own concrete fields to decide.

The Redis adapter has a similar boundary. `go-redis` returns script results as `any` because Redis
commands can return many shapes: integer, string, nil, error, or array. A Lua script returning
`{1, 4, 0}` usually arrives in Go as `[]interface{}{int64(1), int64(4), int64(0)}`. The limiter
must assert the outer array shape before converting each element into the public `Result`.

## Tradeoff

Type mismatches become runtime cases. Algorithms should handle mismatches deliberately instead of assuming they cannot happen.

At package boundaries, prefer the comma-ok assertion and return an error on unexpected shape. Use
the panic form only when the mismatch would mean an impossible internal invariant was violated and
crashing is the chosen policy.

## Links

- [[D28 - Marker Interface for Algorithm State]]
- [[G19 - Marker Interfaces and Opaque State]]
- [[D64 - RedisLimiter Algorithm Interface]]
- [[redis/R05 - KEYS and ARGV Parameters]]

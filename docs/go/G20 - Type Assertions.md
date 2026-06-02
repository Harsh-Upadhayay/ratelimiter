# G20 - Type Assertions

Back to [[Rate Limiter Learning Map]].

## Go practice

When a value is held through an interface, a type assertion retrieves the expected concrete type.

Conceptually:

```go
state, ok := rawState.(fixedWindowState)
```

## Why it matters

The limiter treats state as opaque, but each algorithm needs its own concrete fields to decide.

## Tradeoff

Type mismatches become runtime cases. Algorithms should handle mismatches deliberately instead of assuming they cannot happen.

## Links

- [[D28 - Marker Interface for Algorithm State]]
- [[G19 - Marker Interfaces and Opaque State]]

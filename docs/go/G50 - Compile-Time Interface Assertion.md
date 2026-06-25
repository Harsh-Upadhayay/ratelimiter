# G50 - Compile-Time Interface Assertion

Back to [[Go Concepts Index]].

## Concept

```go
var _ Limiter = (*MemoryLimiter)(nil)
var _ Limiter = (*RedisLimiter)(nil)
```

This is a **compile-time assertion** that a type satisfies an interface. Zero runtime cost; its
only job is to fail the build loudly, at the type's own definition site, if the type ever stops
implementing the interface.

## Token by token

- `var _ = ...` — the blank identifier declares a variable and discards it. Normally Go errors on
  an unused variable; `_` is the explicit exemption. The line type-checks, then vanishes.
- `var _ Limiter = ...` — declaring with explicit type `Limiter` forces the assignment to be
  valid, which is legal **only if the right-hand value's type implements `Limiter`**. That
  assignability rule is the entire mechanism.
- `(*MemoryLimiter)(nil)` — a **conversion expression**: convert `nil` to type `*MemoryLimiter`.
  The parentheses around `*MemoryLimiter` are required syntax (`(T)(v)`). The result is a typed
  nil pointer. Satisfaction is checked on the *type*, never by dereferencing — so a nil pointer is
  enough, with zero allocation and no need for a usable zero value.

## Why the pointer form

`Allow` has a pointer receiver, so only `*MemoryLimiter` carries it in its method set
(see [[G46 - Method Sets and Interface Satisfaction]]). `var _ Limiter = MemoryLimiter{}` would
**fail to compile** — the value type's method set lacks `Allow`.

## Why bother

Without it, "does `*MemoryLimiter` implement `Limiter`?" is only answered implicitly at the first
call site that passes one as a `Limiter` (e.g. `NewMiddleware`). Rename `Allow`, change its
signature, or edit the interface, and the failure appears *there*, in confusing terms. The
assertion pins the contract at the definition site: incompatible change -> build breaks here,
naming exactly which type fell out of conformance. It is executable documentation.

Common in the stdlib: `var _ io.Writer = (*bytes.Buffer)(nil)`.

## Links

- [[G46 - Method Sets and Interface Satisfaction]]
- [[G18 - Structural Interface Satisfaction]]
- [[D83 - MemoryLimiter Implements Limiter via Private Clock]]

# G51 - Interface Narrows the Method Set

Back to [[Go Concepts Index]].

## Concept

A value held through an interface exposes **only the methods the interface declares**, even though
the concrete type may have more. To reach the extra methods you need the **concrete type**, not the
interface.

## The project case

`MemoryLimiter.clock` has static type `clock` (the interface), which declares one method:

```go
type clock interface { now() time.Time }
```

The test double has more:

```go
func (c *testClock) now() time.Time         // satisfies clock
func (c *testClock) Set(t time.Time)        // extra
func (c *testClock) Advance(d time.Duration) // extra
```

Reaching for `Advance` through the interface does not compile:

```go
ml.clock.Advance(time.Minute) // compile error: clock has no method Advance
```

The `*testClock` is there at runtime, but the interface type hides everything beyond `now()`.

## Two ways to drive time in tests

| Option | Code | Verdict |
| --- | --- | --- |
| Type-assert each time | `ml.clock.(*testClock).Advance(d)` | works, ugly, panics if the concrete type changes |
| Return the concrete handle | helper returns `*testClock`; call `clk.Advance(d)` | clean, no assertion |

`newTestMemoryLimiter` returns the `*testClock` so tests hold the concrete type with its full
method set.

## The design point

`testClock` deliberately has **more** methods than the `clock` interface needs. The production
seam stays minimal — `now()` only — so the limiter can *read* time but never *control* it. The
test holds the wider surface (`Advance`/`Set`) because it holds the concrete type.

> Narrow interface for the consumer; richer concrete type for the controller.

Returning the clock is not a workaround — it is the clean expression of "the test owns time
control; the limiter only reads time."

## Links

- [[G46 - Method Sets and Interface Satisfaction]]
- [[G20 - Type Assertions]]
- [[G30 - Exported Interfaces With Unexported Types]]
- [[D83 - MemoryLimiter Implements Limiter via Private Clock]]

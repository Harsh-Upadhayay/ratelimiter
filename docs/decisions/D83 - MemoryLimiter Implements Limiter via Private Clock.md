# D83 - MemoryLimiter Implements Limiter via Private Clock

Back to [[V9 HTTP Middleware Index]].

## Context

`MemoryLimiter.Allow` was `Allow(key string, now time.Time)` — the caller supplied time, which
kept algorithms pure (plan step 2) and let tests time-travel by passing crafted `now` values.

The `Limiter` interface the HTTP middleware depends on is:

```go
type Limiter interface {
    Allow(ctx context.Context, key string) (Result, error)
}
```

There is no `now` parameter. The C09 checkpoint deferred this as a "clock-backed adapter". We
chose instead to make `MemoryLimiter` satisfy `Limiter` **directly** — no adapter type.

## Decision

`MemoryLimiter.Allow` now matches the interface signature. The lost `now` parameter is replaced
by an **injectable but private clock dependency**.

```go
type clock interface { now() time.Time }   // unexported
type realClock struct{}                     // wraps time.Now()

type MemoryLimiter struct {
    algo  memoryAlgorithm
    store StateStore
    clock clock
}
```

- `NewMemoryLimiter(algo, store)` signature is unchanged; it defaults `clock: realClock{}`.
- The clock is **not caller-configurable**. There is no `WithClock` and no public option.
- White-box tests inject a `*testClock` by setting the field directly via the
  `newTestMemoryLimiter` helper (see [[G51 - Interface Narrows the Method Set]]).

### Context handling

`Allow` honors `ctx` with a single cheap check at the top, then returns `ctx.Err()` raw:

```go
if err := ctx.Err(); err != nil {
    return Result{}, err   // context.Canceled / context.DeadlineExceeded
}
```

No per-CAS-attempt check: the loop is bounded and non-blocking, so there is nothing long-running
to interrupt. The only case worth catching is a caller who hands in an already-cancelled context.

### Clock sampled per attempt

`now := ml.clock.now()` is read inside the CAS retry loop, so a retry reflects the actual current
moment. The old code fixed `now` once per call (the caller sampled). Tests are unaffected because
`testClock` only changes on explicit `Advance`/`Set`.

## Why private, not a public option

This **reverses** the earlier application of [[D78 - Functional Options for Middleware]] /
[[D80 - Required Dependencies Outside Functional Options]] to the clock. The premise changed:
a functional option is for a setting the *caller* configures. The clock is not — it is purely an
internal time source whose only non-production use is testing.

Making it private is symmetric with `RedisLimiter`, whose clock is Redis `TIME` (server-side) —
also not a Go-side knob. The rule is now uniform: **a limiter owns its time source; callers do
not inject it.** "I need to test it" is not a reason to widen the public API; white-box tests
reach the private seam.

## Why no adapter

An adapter would add a wrapper type and an indirection for no benefit. `MemoryLimiter` already
had everything it needed; only the signature and the time source had to change. Direct
satisfaction keeps one limiter type per backend, asserted at compile time:

```go
var _ Limiter = (*MemoryLimiter)(nil)
```

See [[G50 - Compile-Time Interface Assertion]].

## Links

- [[D62 - Redis Is a Parallel Limiter Not a Store]]
- [[D78 - Functional Options for Middleware]]
- [[D80 - Required Dependencies Outside Functional Options]]
- [[D84 - Benchmarks Use Production Clock]]
- [[G50 - Compile-Time Interface Assertion]]
- [[G51 - Interface Narrows the Method Set]]
- [[G46 - Method Sets and Interface Satisfaction]]

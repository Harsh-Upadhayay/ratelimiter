# C10 - MemoryLimiter Implements Limiter Checkpoint

Back to [[V9 HTTP Middleware Index]].

## Status

Complete. `MemoryLimiter` now satisfies the `Limiter` interface **directly** — no adapter — so it
can be passed to `NewMiddleware` exactly like `RedisLimiter`. Build, tests, `-race`, and benchmarks
are all green.

This closes the C09 deferred branch "memory adapter for the `Limiter` interface", chosen as a
direct refactor rather than a wrapper.

## What changed

```text
Allow(key, now)            ->  Allow(ctx, key)     // matches Limiter
caller-supplied now        ->  private clock dependency (realClock default)
```

- New `clock.go`: unexported `clock` interface (`now() time.Time`) + `realClock`.
- `MemoryLimiter` gained a private `clock` field; `NewMemoryLimiter` signature unchanged, defaults
  to `realClock{}`. Clock is **not** caller-configurable ([[D83 - MemoryLimiter Implements Limiter via Private Clock]]).
- `Allow` honors `ctx` via a single `ctx.Err()` check; samples `clock.now()` per CAS attempt.
- Compile-time assertions in `limiter.go`:
  `var _ Limiter = (*MemoryLimiter)(nil)` / `(*RedisLimiter)(nil)`.
- Tests migrated to a white-box `*testClock` (mutex-guarded, advanceable) returned from
  `newTestMemoryLimiter`.
- Benchmarks build through the production constructor (lock-free `realClock`), not the test clock
  ([[D84 - Benchmarks Use Production Clock]]).

## Finding uncovered

Removing the accidental test-clock mutex from the hot path exposed real single-key CAS contention:
under sustained same-key load the bounded retry loop can return `ErrCASConflict`. Accepted as
expected; benchmarks tolerate it; mitigation deferred —
[[D85 - Hot Key CAS Exhaustion Is Expected]].

## Verified

```text
go build ./...
go test ./...
go test -race ./...
go test -run '^$' -bench BenchmarkAllow -benchmem ./...
```

All pass at checkpoint time.

## Deferred

Still open from C09 and unchanged:

- middleware behavior tests,
- observability / logging / metrics,
- local HTTP example server,
- sliding window algorithms,
- delegate-on-error policy,
- Redis integration / runtime tests,
- hot-key CAS mitigation ([[D85 - Hot Key CAS Exhaustion Is Expected]]).

## Decisions and concepts

- [[D83 - MemoryLimiter Implements Limiter via Private Clock]]
- [[D84 - Benchmarks Use Production Clock]]
- [[D85 - Hot Key CAS Exhaustion Is Expected]]
- [[G50 - Compile-Time Interface Assertion]]
- [[G51 - Interface Narrows the Method Set]]
- [[G52 - go test Compiles the Whole Package]]

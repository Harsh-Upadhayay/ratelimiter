# D85 - Hot Key CAS Exhaustion Is Expected

Back to [[V9 HTTP Middleware Index]].

## How this surfaced

[[D84 - Benchmarks Use Production Clock]] removed a shared test-clock mutex from the hot path.
That mutex had been accidentally serializing the same-key parallel benchmarks, so CAS rarely
conflicted. With the real lock-free clock, the same-key parallel benchmarks immediately failed:

```text
max number of CAS conflict attempts exhausted   (ErrCASConflict)
```

The fake clock had been **masking the limiter's true behavior under a hot key.**

## The finding

`MemoryLimiter.Allow` uses a bounded optimistic CAS loop: `Get -> Decide -> CompareAndSwap`,
up to 10 attempts, then `ErrCASConflict` (see [[G26 - Optimistic Concurrency with CAS]]).

Under sustained contention on a **single key**, `GOMAXPROCS` goroutines fight over **one record's
version**. CAS lets only one winner per round, so an unlucky goroutine can lose all 10 attempts
and the caller receives `ErrCASConflict` — a valid request rejected with an error.

Sharding does **not** help here: same key -> same shard -> same record. This is the classic
"hot key" problem (plan step 6), and it is intrinsic to optimistic concurrency under high
single-key write contention, not a bug in the loop.

## Decision

1. **Accept it as expected behavior** for now. The 10-retry bound stays. `ErrCASConflict` is a
   truthful signal that a key is too hot for optimistic CAS, not a defect to silence.
2. **Benchmarks tolerate it.** The parallel benchmark loops treat `ErrCASConflict` as a
   non-fatal outcome (`errors.Is(err, ErrCASConflict)`), so they measure throughput honestly
   instead of `b.Fatal`-ing on a legitimate result.

## Deferred — hot-key mitigation

Not built; logged as a real future thread. Candidate strategies, each with trade-offs to weigh:

- **Backoff / jitter between retries** — reduces livelock, adds latency.
- **Higher / unbounded retry bound** — fewer rejections, risks longer tail latency under load.
- **Per-key serialization for hot keys** — a lock manager instead of CAS for contended keys;
  trades throughput for guaranteed progress.
- **Move hot keys to the Redis path** — Lua's single-shot pessimistic atomicity has no CAS retry
  (see [[D63 - Fixed Window in Redis via Lua]]).

## Why this matters

This is arguably the most valuable thing the V9 clock refactor uncovered: a correctness-relevant
behavior under load that was invisible while a test artifact serialized the hot path. It is a
reminder that test doubles can hide the very properties a benchmark exists to reveal.

## Links

- [[D84 - Benchmarks Use Production Clock]]
- [[G26 - Optimistic Concurrency with CAS]]
- [[D63 - Fixed Window in Redis via Lua]]
- [[V7 Sharded MemoryStore Index]]

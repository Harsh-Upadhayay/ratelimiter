# Current Assessment — Interview/Résumé Readiness

**Date:** 2026-06-24  
**Baseline:** assume all deferred V9 tasks completed (middleware tests, observability, local HTTP example, memory `Limiter` adapter, delegate-on-error policy).

---

## The hidden design tension (read this first)

The `Limiter` interface and `MemoryLimiter` don't agree:

```
Limiter         Allow(ctx context.Context, key string)   // limiter.go
RedisLimiter    Allow(ctx context.Context, key string)   // satisfies Limiter
MemoryLimiter   Allow(key string, now time.Time)         // does NOT satisfy Limiter
```

The middleware calls `m.limiter.Allow(r.Context(), key)`, so today only `RedisLimiter` can go through the middleware. The deferred "memory adapter for `Limiter`" task is papering over a genuine philosophical fork: `RedisLimiter` threads `context` because it does I/O; `MemoryLimiter` threads an injected clock because it is pure. This is the most interview-revealing design decision left in the project — it isn't just labor, it's an architectural choice. See [[V7 - Sharded MemoryStore#ADR-0062 — Redis as a limiter, not a store (pivot)|ADR-0062]] for the earlier pivot that set up this fork.

---

## A. System-design assessment

**Verdict: well above median, including senior.**

### Genuine strengths
- **Storage/logic split + pure `Decide`** — textbook-correct factoring; most candidates never get there.
- **D62 pivot** — recognizing that "Redis as a `StateStore`" is wrong because in-process CAS passes live Go structs while Redis owns atomicity server-side is a Staff-level insight. The ability to narrate *why* the abstraction was abandoned carries an interview on its own.
- Clock injection + Redis `TIME` for skew, lock striping for hot-key contention, CAS-vs-Lua atomicity contrast — all real distributed-systems decisions with documented rationale.

### Gaps between "very strong" and "best"

| Gap | Why it matters |
|-----|----------------|
| Single-Redis-node correctness only | On failover to a stale replica, the key counter can be double-counted. No replication/cluster story yet. |
| Every request is a Redis round-trip | Production limiters (Stripe, Cloudflare, GCRA) lease tokens locally and reconcile asynchronously to avoid network-latency on the hot path. Sharding solves in-process contention; it doesn't touch cross-node latency. |
| Unbounded in-process memory | `MemoryStore` keys live forever; no eviction policy. Redis has TTL; memory does not. |
| GCRA / leaky bucket absent | The canonical "smooth" limiter and a common Staff-level interview flex. |

None of these are disqualifying. They mark the difference between "I built a correct rate limiter" and "I built one that survives Redis failover at scale." Don't claim the larger thing in an interview — the first follow-up ("what happens when the Redis primary fails over?") would expose it.

---

## B. Go concepts assessment

**Verdict: genuinely strong, and in places senior.**

### Genuine strengths
- **Sealed lowercase interfaces** (`memoryAlgorithm`, `redisAlgorithm`) — knowing what *not* to export is a senior tell.
- Opaque marker-interface state, pointer receivers to avoid copying mutexes (G33), sentinel errors, functional options, `net/http` middleware idiom.
- **Contract tests shared across implementations** (G37), race-detector discipline, benchmark-*driven* sharding decision (D48 → D53).

### Points an interviewer will poke

| Point | Detail |
|-------|--------|
| `context` inconsistency | Idiomatic Go threads `context` through any call that *might* do I/O. `MemoryLimiter` skipping it is defensible (it can't block), but the divergent signatures look like an oversight until explained. |
| Generics | Go 1.22 — why opaque interfaces for state instead of `Limiter[S State]`? Probably correct, but be ready for "why not." |
| Error wrapping | Sentinels are good; confirm `%w` + `errors.Is/As` is the documented contract. |
| `memory_token_bucket_tests.go` | Misnamed: Go never compiles it as tests. Has an open TODO list. A visible blemish on a résumé project. |

---

## Bottom line

**Résumé-ready: yes, comfortably** — more impressive than most because of the documented decision trail (D62 especially). The narrative is the differentiator.

**"Final best": no, and the gap is informative.** This is a *single-node-correct, contention-optimized, multi-algorithm* limiter, not a *failover-correct, latency-aware, internet-scale* one. That is a fine place to stop for learning.

### Open questions for the next conversation
1. Should `MemoryLimiter` adopt the `Limiter` signature (and where does its clock go?), or are these two different abstractions that shouldn't share one interface?
2. Of the four sysdesign gaps, which does a Staff interviewer reach for first — and is closing it a *learning* win or just a *résumé* win?

# V6 — Storage Boundary (ADR-0045 … ADR-0052)

> Back to [[Rate Limiter Learning Map]] · Prev [[V5 - Token Bucket]] · Next [[V7 - Sharded MemoryStore]] · Index [[README|ADR index]]

V6 pulls storage out of the limiter. The limiter should orchestrate without knowing whether state lives in a map, Redis, or elsewhere (ADR-0045). A plain `Get`/`Set` store has a lost-update race (ADR-0046), so the boundary is built on `Get` + compare-and-swap with versions (ADR-0047): the limiter reads with a version, decides, and commits only if the version still matches, retrying on conflict. A normal CAS conflict is an expected outcome, not a store error (ADR-0050), and the retry loop is bounded so the request path can't hang (ADR-0051). `MemoryStore` now owns the map, versions, and mutex (ADR-0049); the refactor is motivated by a benchmark that first proves the single-mutex bottleneck (ADR-0048). For now `NewLimiter(algo)` still creates a default store internally (ADR-0052) — a convenience V7 reverses.

---

## ADR-0045 — Split storage from the limiter
**Status:** Accepted

**Context:** `Limiter` still owns the state map and mutex, coupling orchestration, storage, and concurrency control.

**Decision:** Move storage behind a store boundary; the limiter orchestrates validate → load → decide → commit without knowing the backend.

**Consequences:** Another interface, and the design must now handle read-modify-write correctness explicitly (ADR-0047) instead of leaning on one process-local mutex.

## ADR-0046 — The Get/Set race
**Status:** Accepted (problem statement)

**Context:** A store exposing only `Get` and `Set`.

**Decision:** Recognize that two requests reading the same old state before either writes can both be allowed when only one should — `Get`/`Set` is insufficient without a shared atomic guard.

**Consequences:** Directly motivates CAS (ADR-0047). This is the same class of bug as the Redis check-then-act race (Redis Concepts: *SET with NX/EX*, *levels of atomicity*).

## ADR-0047 — StateStore uses Get + CAS
**Status:** Accepted

**Context:** The limiter must read, decide, and commit without losing concurrent updates.

**Decision:** A store boundary of `Get(key) → (state, version, …)` plus `CompareAndSwap(key, version, newState)`. On CAS failure, re-read and re-decide.

**Consequences:** Preserves read-decide-write correctness without a single global limiter mutex, matching the distributed direction. Adds versions, retries, and conflict handling — the optimistic-concurrency model (Go Concepts: *optimistic concurrency with CAS*).

## ADR-0048 — Benchmark before the storage refactor
**Status:** Accepted

**Context:** About to change the concurrency design.

**Decision:** Add concurrency tests and benchmarks first to expose the current bottleneck.

**Consequences:** Same-key and many-key parallel benchmarks both serialize behind one limiter mutex — unrelated keys can't progress independently. The refactor is now driven by observed pressure, and this benchmark becomes the baseline for V7's sharding win (and later the subject of ADR-0084).

## ADR-0049 — MemoryStore owns runtime state
**Status:** Accepted

**Context:** `Limiter` previously owned the map and mutex; V6 introduces the `StateStore` boundary.

**Decision:** `MemoryStore` owns the key→record map, per-key versions, and the mutex; `Limiter` owns orchestration only.

**Consequences:** Storage concerns separate from decision flow. Still one mutex internally — this creates the boundary, not yet the many-key concurrency win (V7).

## ADR-0050 — CAS conflict is not a store error
**Status:** Accepted

**Context:** `CompareAndSwap` can fail because another caller updated the key after `Get`.

**Decision:** Represent a normal conflict as `ok = false, err = nil`; reserve `err` for real store failures.

**Consequences:** A version mismatch is an expected concurrency outcome the limiter retries; operational failures stay distinct. Callers must check both return values.

## ADR-0051 — Bounded CAS retry loop
**Status:** Accepted

**Context:** The limiter may lose a race and need to retry.

**Decision:** Retry a bounded number of times (10) in `Allow`, then return `ErrCASConflict`.

**Consequences:** Tail latency stays finite and the caller gets an operational error it can fail-open/closed on. Under hot-key contention a valid request can be rejected with `ErrCASConflict` — accepted, and revisited in ADR-0085. (The Redis path avoids this entirely via single-shot Lua — ADR-0063.)

## ADR-0052 — Default memory store in the constructor
**Status:** Superseded by ADR-0059

**Context:** One option was to force every caller/test to pass a store.

**Decision:** Keep `NewLimiter(algo)` and create a default `MemoryStore` internally for now.

**Consequences:** Keeps the API and existing tests simple while the boundary settles. Defers store injection — a future Redis/alternate store needs a new constructor or options. V7 opens injection once a second backend exists.

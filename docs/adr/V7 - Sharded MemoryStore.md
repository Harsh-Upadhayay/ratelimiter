# V7 — Sharded MemoryStore (ADR-0053 … ADR-0062)

> Back to [[Rate Limiter Learning Map]] · Prev [[V6 - Storage Boundary]] · Next [[V8 - RedisLimiter]] · Index [[README|ADR index]]

V7 attacks the many-key bottleneck the V6 benchmark exposed: a sharded memory store using lock striping (ADR-0053). It slots in behind the unchanged `StateStore` interface (ADR-0054), proving the boundary — backend concurrency improves without touching algorithms or orchestration. Shards are a configurable count of reused `*MemoryStore` values (ADR-0055, ADR-0056), keys are hashed to shards, and internal calls trust their own invariants rather than bounds-checking (ADR-0058). Along the way store injection finally opens (ADR-0059), and `StateStore` gets a shared contract test with fresh-per-test isolation (ADR-0060, ADR-0061). V7 closes the single-node depth axis and ends with the pivot that reframes everything distributed: Redis will be a parallel limiter, not a store (ADR-0062).

---

## ADR-0053 — Shard the memory store next
**Status:** Accepted

**Context:** Another algorithm would mostly repeat the existing pattern; the memory backend still serializes all keys behind one mutex (ADR-0048).

**Decision:** Build a sharded memory store before adding more algorithms.

**Consequences:** Teaches partitioning ownership by key, driven by the many-key benchmark. Sharding helps spread traffic, not a single hot key, and adds hashing and shard-count tuning. (Chosen over a per-key lock manager and over another algorithm as the higher-learning move.)

## ADR-0054 — Sharded store keeps the StateStore contract
**Status:** Accepted

**Context:** `StateStore` already exposes `Get` and `CompareAndSwap`.

**Decision:** Implement sharding behind the existing interface, unchanged.

**Consequences:** The limiter doesn't care whether the backend has one mutex, many, or Redis — validating the V6 boundary. Sharding details become backend internals; inspecting shard distribution needs backend-specific test access.

## ADR-0055 — Configurable shard count
**Status:** Accepted

**Context:** A sharded store needs a fixed number of shards.

**Decision:** `NewShardedMemoryStore(shardCount)` with validation `1 ≤ shardCount ≤ MAXSHARDSIZE` (else `ErrInvalidShardCount`).

**Consequences:** Different counts make the memory-vs-concurrency trade visible for experiments. The constructor needs validation — a system boundary (contrast ADR-0058).

## ADR-0056 — Reuse MemoryStore internally
**Status:** Accepted

**Context:** Shards can be a dedicated type or composed from existing stores.

**Decision:** Shards are `[]*MemoryStore` — pointers, so mutex-containing values aren't copied.

**Consequences:** Reuses existing `Get`/`CAS`, avoiding duplicated logic. Slightly less semantically precise than a private `memoryShard`, but smaller for this step (see Go Concepts: *composition with pointer fields*).

## ADR-0057 — Sharded store as the default backend
**Status:** Superseded by ADR-0059

**Context:** Injection was still deferred (ADR-0052); `NewLimiter(algo)` created the default store.

**Decision:** Make `ShardedMemoryStore` the default backend, with a safe internal default shard count.

**Consequences:** Lets the new backend affect existing tests/benchmarks without exposing backend choice. Short-lived — ADR-0059 opens injection, removing the auto-default.

## ADR-0058 — Trust internal invariants
**Status:** Accepted

**Context:** `Get`/`CAS` use the index from the private `shardIndex`.

**Decision:** Don't bounds-check it. `shardIndex` always returns `h.Sum32() % len(shards)`, structurally in range; a defensive check would only guard your own code.

**Consequences:** Rule — validate at system boundaries (user input, constructor params per ADR-0055); trust internal contracts between methods on the same type. A check here would only swap a panic for a silent error without fixing a real bug.

## ADR-0059 — Store injection opened
**Status:** Accepted

**Context:** ADR-0052 deferred injection; now a real alternative backend (`ShardedMemoryStore`) exists.

**Decision:** `NewLimiter(algo, store)` — the store is an explicit parameter; callers and tests choose the backend. There is no auto-default.

**Consequences:** Reverses ADR-0052 and ADR-0057. Tests need a `newTestLimiter` helper anyway, so passing a store is no burden; hardcoding a default would hide the choice. A convenience `NewLimiterWithDefaults` could be added later without breaking this signature.

## ADR-0060 — Contract testing for StateStore
**Status:** Accepted

**Context:** Two backends must behave identically from the caller's view.

**Decision:** One shared helper `testStateStoreContract(t, store StateStore)`; each backend gets a thin `Test...` that passes its instance. The helper covers missing-key reads, CAS create/update/stale-version paths, and concurrent access.

**Consequences:** No duplicated test logic; a third backend (e.g. Redis) is one new `Test...`. The helper must use only the interface — no concrete type assertions (Go Concepts: *contract testing via interface parameter*).

## ADR-0061 — Test isolation with fresh stores
**Status:** Accepted

**Context:** A shared package-level store accumulates state across tests.

**Decision:** Every test builds its own store (`newTestLimiter` constructs fresh each call; contract tests build a new store per `Test...`).

**Consequences:** Eliminates order-dependent test pollution (Test A leaves version 1; Test B assumes 0 and gets a spurious CAS conflict).

## ADR-0062 — Redis as a limiter, not a store (pivot)
**Status:** Accepted — revises plan step 4

**Context:** The instinct was a `RedisStore` implementing `StateStore`. But `StateStore` passes live Go `algorithmState` structs; Redis stores bytes. Bridging needs serialization, which needs type knowledge — every option (store knows all types / store holds an algorithm reference / change the interface to `[]byte`) forces a bad trade and reveals that `StateStore` is fundamentally an in-process abstraction.

**Decision:** Redis will *not* implement `StateStore`. A parallel `RedisLimiter` implements `Allow` directly, with algorithm logic living in Redis Lua scripts. Redis server-side atomicity replaces the CAS retry loop entirely.

**Consequences:** `Limiter` stays in-process only (`MemoryStore`, `ShardedMemoryStore`); `RedisLimiter` is a parallel implementation, not a backend swap; both expose the same external behavior. Algorithm logic is duplicated — once in Go, once in Lua — the accepted cost of true distributed atomicity. This pivot drives all of V8.

> **Note (revised):** This is the original north-star plan, kept for reference.
> Steps 3–4 were **superseded by ADR-0062**: Redis is **not** a `StateStore`
> backend. `StateStore` is an in-process abstraction (it passes live Go state
> structs); Redis stores bytes and owns atomicity server-side. So the distributed
> path is a **parallel `RedisLimiter`** (single-shot Lua) satisfying the same
> `Limiter` interface — not a swappable store, and there is no `DistributedLimiter`.
> See [`docs/adr/README.md`](docs/adr/README.md) for the as-built design.

**1. Split storage from logic.**  
Define a `StateStore` interface (Get/CAS/PushEvent) so algorithms never touch Redis or memory directly.

**2. Make algorithms pure functions.**  
Each algorithm implements a single `Decide(now, state) → (result, newState)` method. No side effects, no `time.Now()`, no I/O.

**3. Write one distributed limiter.**  
A `DistributedLimiter` wraps any `Algorithm` + `StateStore`. It reads state, calls `Decide`, then commits with optimistic locking (CAS). On conflict, it retries automatically.

**4. Build two backends.**  
- `MemoryStore` with `sync.Mutex` for local dev/tests.  
- `RedisStore` using Lua scripts (for stateful algos like Token Bucket) and Sorted Sets (for event-based algos like Sliding Window Log).

**5. Add three algorithms.**  
Token Bucket, Sliding Window Counter, and Sliding Window Log. All plug into the same limiter without code changes.

**6. Harden.**  
Injectable clock for tests, Redis `TIME` to handle clock skew, key sharding for hot keys, and metrics tagged by algorithm/store.

---

## Learning notes

The V1 tradeoff log and Go practice notes now live in [[Rate Limiter Learning Map]].

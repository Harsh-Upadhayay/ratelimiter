# Rate Limiter — Learning Project

This is a **learning project**. Harsh is building a rate limiter in Go to learn both Go
syntax and distributed-systems design. The code reaching production is secondary to the
learning. Treat me (Claude) as a **Senior Staff Go Engineer and System Design Mentor**.

## Mentorship contract (HIGHEST PRIORITY — do not violate)

These rules override my usual "just implement it" instinct:

1. **Do NOT write code unless Harsh explicitly requests it.** Default to guiding, not doing.
2. **Do NOT give the final architectural answer upfront.** Reveal the design one piece at a
   time.
3. **Guide toward a simpler baseline first**, before going down the full plan.
4. **Use the Socratic method** — walk the 6-step plan one piece at a time, asking questions.
5. **On implementation questions:** explain the underlying Go concept, then ask Harsh how it
   should apply here.
6. **On design decisions:** present alternatives and ask Harsh to evaluate the tradeoffs
   (memory, latency, concurrency) before choosing.

When Harsh does ask for code, write it in the existing house style (early returns/guard
clauses, pointer receivers for mutating types, sentinel errors, pure helpers).

## Documentation conventions (keep these up to date as we work)

- **Design decisions** → `docs/decisions/Dxx - Title.md` (numbered, sequential; currently
  through D69).
- **Go concepts** → `docs/go/Gxx - Title.md` (currently through G39).
- **Redis/Lua concepts** → `docs/redis/Rxx - Title.md` (currently through R07; hub:
  `docs/Redis Concepts Index.md`).
- **Version index hubs** → `docs/Vn ... Index.md` linking the decisions/concepts for that
  iteration. Hub-and-spoke notes use `[[wikilinks]]` (Obsidian-style).
- Master map: `docs/Rate Limiter Learning Map.md`. Full narrative: `docs/Chat Export - Rate
  Limiter Learning Session.md`.
- `plan.md` = the north-star 6-step target plan.

## The 6-step target plan (from plan.md)

1. Split storage from logic via a `StateStore` interface.
2. Make algorithms pure functions: `Decide(now, state) → (result, newState)` — no side
   effects, no `time.Now()`, no I/O.
3. One distributed limiter wrapping any algorithm + store, committing with CAS, retrying on
   conflict.
4. Two backends: `MemoryStore` (mutex) and `RedisStore` (Lua + sorted sets).
5. Three more algorithms: Token Bucket, Sliding Window Counter, Sliding Window Log.
6. Harden: injectable clock, Redis `TIME` for skew, key sharding for hot keys, metrics.

> **Plan revised by D62:** step 4's "RedisStore as a `StateStore` backend" was abandoned.
> `StateStore` is an in-process abstraction (it passes live Go `algorithmState` structs;
> Redis stores bytes and owns atomicity server-side). Redis becomes a **parallel
> `RedisLimiter`** implementing `Allow` directly via Lua scripts — NOT a swappable store.
> So `Limiter` stays in-process only; the distributed path is a separate implementation.

## Architecture so far (V1 → V7, build green)

- **Package:** `ratelimiter` (module `github.com/Harsh-Upadhayay/ratelimiter`, Go 1.22.2).
- `Limiter` (`limiter.go`) holds `algo algorithm` + `store StateStore`. `NewLimiter(algo, store)`
  takes both — store injection is open (D59 reversed D52). Callers pick the backend; there is
  **no** auto-default to a particular store. Tests use `newTestLimiter(t, algo)` which wraps
  `NewLimiter` with a fresh `MemoryStore`.
- `Allow(key, now)` does `Get → Decide → CompareAndSwap` in a **bounded CAS retry loop**
  (10 attempts → `ErrCASConflict`). Retry re-runs `Decide`, not just CAS.
- `algorithm` + `algorithmState` (`types.go`) are **private interfaces**; `algorithmState`
  is a marker interface for opaque per-key state. `Result{Allowed, Remaining, RetryAfter}`
  is exported.
- `StateStore` (`state_store.go`): `Get(key) → (state, version, exists, error)` and
  `CompareAndSwap(key, version, state) → (ok, error)`. Missing key = version 0; CAS conflict
  is `ok=false, err=nil` (not an error).
- `MemoryStore` (`memory_store.go`): `map[string]record` + a **field** `sync.Mutex` (not a
  local var). Locks both reads and writes.
- Algorithms: `FixedWindow` and `TokenBucket` (exported constructors, own their validation).
  Token Bucket: lazy refill, `float64` tokens, clamps backward clock movement.
- `ShardedMemoryStore` (`sharded_memory_store.go`): lock striping over `shards []*MemoryStore`
  (pointers, so mutex-containing values aren't copied — G33). `NewShardedMemoryStore(shardCount)`
  validates `1 <= shardCount <= MAXSHARDSIZE` (1000) else `ErrInvalidShardCount` (D55).
  `shardIndex` uses `hash/fnv` 32a, `% len(shards)` (G32, G34, G35). `Get`/`CompareAndSwap`
  delegate to `shards[shardIndex(key)]`, reusing `*MemoryStore` (D56). Internal calls trust
  `shardIndex` — no bounds check (D58).
- Tests go **through the public API** to preserve refactoring freedom. `StateStore`
  implementations share one contract-test helper (D60, G37); every test builds a fresh store
  for isolation (D61, G36). Benchmarks in `limiter_benchmark_test.go` showed many-key parallel
  does NOT scale on a single `MemoryStore` mutex — the pressure that motivated sharding (D48).

## V7 — Sharded MemoryStore (COMPLETE, build + tests + race all green)

Lock striping fixed the many-key serialization bottleneck. Chosen over per-key lock manager
and over adding another algorithm (both lower learning value at the time — D53). The sharded
store is an available backend callers can inject; it is not auto-wired as a default.

## Current direction: V8 — breadth over the distributed axis

Decision made 2026-06-14: the single-node depth axis (sharding/contention/testing) is near
its learning ceiling. Next learning lives in **breadth** — the distributed path and the
sliding-window algorithm family. **D62 is the pivot that reframed this:** Redis is a parallel
`RedisLimiter`, not a `StateStore`.

Order of attack chosen (D63): **Redis-first** — port Fixed Window into a `RedisLimiter` to
learn `EVAL`/Lua mechanics on a simple algorithm, before Token Bucket *forces* Lua (read
tokens + timestamp → compute refill → conditionally deduct → write back). Sliding Window
Log/Counter (sorted sets, `ZADD`/`ZREMRANGEBYSCORE`) comes later. Redis/Lua concepts are
logged in `docs/redis/` (R01–R07); the atomicity reasoning lives in
[[R07 - Levels of Atomicity in Redis]] and [[D63 - Fixed Window in Redis via Lua]]. Lua's
pessimistic single-shot atomicity replaces the in-process CAS retry loop — the Redis path has
no `ErrCASConflict`. Hub: `docs/V8 RedisLimiter Design Index.md`. Redis skeleton code exists for
`RedisLimiter`, `RedisFixedWindow`, and `goRedisAdapter`; fixed-window Lua is sketched, but Redis
execution and result parsing are not fully wired yet.

Decision D65: `goRedisAdapter.eval` should return raw `any`; `RedisLimiter.Allow` owns parsing
the Lua contract `{allowed, remaining, retryAfterSeconds}` into `Result`. The adapter is Redis I/O
plumbing, not rate-limiter domain logic.

Decision D66: public Redis construction should use a caller-owned Redis client, while `redisAdapter`
and `goRedisAdapter` stay private. The app owns connection lifecycle/configuration; `RedisLimiter`
owns rate-limiting behavior; the adapter remains an internal test/swap seam.

Decision D67: Redis algorithms are sealed for now. `redisAlgorithm` stays private, so callers use
package-provided Redis algorithms while the Lua/result contract is still evolving.

Decision D68: share fixed-window policy config/validation between `FixedWindow` and
`RedisFixedWindow`, but keep execution behavior separate. Stable policy is shared; volatile
backend execution details stay isolated.

Decision D69: backend-specific concepts use qualifier-first names: `MemoryLimiter`, `RedisLimiter`,
`MemoryFixedWindow`, `RedisFixedWindow`, etc. Internal sealed interfaces remain lowercase
(`memoryAlgorithm`, `redisAlgorithm`) unless deliberately opened as public extension points.

## Known cleanup items (mention when relevant; don't fix unprompted)

- `token_bucket_tests.go` is misnamed — Go only compiles `_test.go` files, so this file is
  NOT run as tests. It also has an open `// TODO: Add tests` list (bucket-full start,
  rejection, refill, sub-token, cap, retry-after, backward time).

## Commands

- Build: `go build ./...`
- Test: `go test ./...` and `go test -race ./...`
- Benchmarks: `go test -bench=BenchmarkAllow.*FixedWindow -benchmem ./...`

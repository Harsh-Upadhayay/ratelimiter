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

Docs were consolidated 2026-06-25 from 165 atomic files (one note per file became unreadable)
into ADR logs + concept docs. Everything uses `[[wikilinks]]` (Obsidian-style):

- **ADRs** (`docs/adr/Vn - Title.md`) — architecture decisions `ADR-0001` … `ADR-0085` in
  standard Nygard format (**Status / Context / Decision / Consequences**), grouped one log per
  iteration (V1–V9), each opening with a narrative intro. ADR number == the old decision number
  (ADR-0038 was D38). Index: `docs/adr/README.md` (table of all ADRs + status). Currently
  through ADR-0085. Checkpoints (old C04/C09/C10) are folded into the end of their version log.
- **Concepts** (`docs/concepts/`) — `Go Concepts.md` (G-notes, through the old G52) and
  `Redis Concepts.md` (R-notes, through the old R11), condensed into themed prose. Each concept
  names the ADR(s) it supports.
- **Adding a decision**: append an `## ADR-00NN — Title` section (Status/Context/Decision/
  Consequences) to the **current** version log, add a row to `adr/README.md`, and refresh that
  log's narrative intro if needed. To reverse a past decision, add a new ADR and set the old
  one's status to `Superseded by 00NN`. Do **not** create a new file per decision. Start a new
  log only when a new version (V10+) begins.
- **Adding a concept**: append to the relevant theme section of `Go Concepts.md` /
  `Redis Concepts.md`, referencing the ADR(s) that use it.
- Cross-reference decisions inline as plain `ADR-00NN`; navigation lives in `adr/README.md`.
  Concept docs and ADR logs link each other with `[[wikilinks]]`.
- Master map / entry point: `docs/Rate Limiter Learning Map.md`. Full narrative:
  `docs/Chat Export - Rate Limiter Learning Session.md`.
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
> `StateStore` is an in-process abstraction (it passes live Go `memoryAlgorithmState` structs;
> Redis stores bytes and owns atomicity server-side). Redis becomes a **parallel
> `RedisLimiter`** implementing `Allow` directly via Lua scripts — NOT a swappable store.
> So `MemoryLimiter` stays in-process only; the distributed path is a separate implementation.

## Architecture so far (V1 → V7, build green)

- **Package:** `ratelimiter` (module `github.com/Harsh-Upadhayay/ratelimiter`, Go 1.22.2).
- `MemoryLimiter` (`memory_limiter.go`) holds `algo memoryAlgorithm` + `store StateStore`. `NewMemoryLimiter(algo, store)`
  takes both — store injection is open (D59 reversed D52). Callers pick the backend; there is
  **no** auto-default to a particular store. Tests use `newTestMemoryLimiter(t, algo)` which wraps
  `NewMemoryLimiter` with a fresh `MemoryStore`.
- `Allow(ctx, key)` does `Get → Decide → CompareAndSwap` in a **bounded CAS retry loop**
  (10 attempts → `ErrCASConflict`). Retry re-runs `Decide`, not just CAS. As of D83
  `MemoryLimiter` satisfies the `Limiter` interface **directly** (no adapter): `Allow` dropped the
  `now` parameter for a **private, non-injectable `clock`** (`realClock` default; symmetric with
  `RedisLimiter` owning Redis `TIME`). `ctx` is honored via one `ctx.Err()` check; `now` is sampled
  per CAS attempt. Compile-time assertions `var _ Limiter = (*MemoryLimiter)(nil)` /
  `(*RedisLimiter)(nil)` live in `limiter.go`. Under sustained single-key contention the CAS loop
  can legitimately return `ErrCASConflict` (hot-key problem, D85); benchmarks tolerate it.
- `memoryAlgorithm` + `memoryAlgorithmState` (`types.go`) are **private interfaces**; `memoryAlgorithmState`
  is a marker interface for opaque per-key state. `Result{Allowed, Remaining, RetryAfter}`
  is exported.
- `StateStore` (`state_store.go`): `Get(key) → (state, version, exists, error)` and
  `CompareAndSwap(key, version, state) → (ok, error)`. Missing key = version 0; CAS conflict
  is `ok=false, err=nil` (not an error).
- `MemoryStore` (`memory_store.go`): `map[string]record` + a **field** `sync.Mutex` (not a
  local var). Locks both reads and writes.
- Algorithms: `MemoryFixedWindow` and `MemoryTokenBucket` (exported constructors, own their validation).
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

Decision D68: share fixed-window policy config/validation between `MemoryFixedWindow` and
`RedisFixedWindow`, but keep execution behavior separate. Stable policy is shared; volatile
backend execution details stay isolated.

Decision D69: backend-specific concepts use qualifier-first names: `MemoryLimiter`, `RedisLimiter`,
`MemoryFixedWindow`, `RedisFixedWindow`, etc. Internal sealed interfaces remain lowercase
(`memoryAlgorithm`, `redisAlgorithm`) unless deliberately opened as public extension points.

Decisions D70-D73: Redis Token Bucket stores scaled integer token units in a Redis hash
(`tokens`, `last_refill_ms`) using package-level unexported `tokenScale = 1000`. Redis `TIME`
owns the script clock; TTL is full-refill idle cleanup and is refreshed on every script call.
Redis algorithms keep the common Lua result shape `{allowed, remaining, retryAfterSeconds}`;
token bucket floors `remaining` to whole tokens and ceilings `RetryAfter` to seconds. Scaling is
allowed inside `RedisTokenBucket` and its Lua script, but must not leak into `RedisLimiter`,
`Result`, or caller-facing API. The focused arithmetic review note is
`docs/redis/R11 - Redis Token Bucket Arithmetic.md`.

## Current direction: V9 — HTTP middleware

V9 turns the limiter into an HTTP request-path integration. The public middleware API is
behavior-oriented and package-qualified: `ratelimiter.Middleware` / `ratelimiter.NewMiddleware`,
not `RedisMiddleware`. This naming does not force a fully generic limiter interface immediately; it
only keeps Redis out of the middleware's public behavior name.

Decisions D74-D82 capture the current shape: HTTP middleware next, whole-middleware failure policy,
caller-provided `KeyFunc`, status-and-headers-only responses, functional options, and behavior
named middleware. Required middleware dependencies (`Limiter`, `KeyFunc`) are explicit constructor parameters; functional options are reserved for optional settings like failure policy. Fail-open is silent pass-through; delegate-on-error is deferred. `Retry-After` stays as `time.Duration` in Go and is converted to ceiling seconds only at the HTTP boundary. Go concepts G41-G49 capture `net/http` middleware, functional options, `iota`
enum-like constants, function callback types, `ResponseWriter` header/status behavior, and method
sets for interface satisfaction.

Checkpoints: `docs/checkpoints/C09 - V9 HTTP Middleware Checkpoint.md` and `C10 - MemoryLimiter
Implements Limiter Checkpoint.md` (C10: `MemoryLimiter` satisfies `Limiter` directly via a private
clock — D83/D84/D85, G50/G51/G52). Deferred branches: middleware tests, observability, local HTTP
example server, sliding window algorithms, delegate-on-error policy, Redis integration/runtime
tests, and hot-key CAS mitigation (D85). The "memory adapter for `Limiter`" branch is **done** (as
a direct refactor, not an adapter).

## Known cleanup items (mention when relevant; don't fix unprompted)

- `memory_token_bucket_tests.go` is misnamed — Go only compiles `_test.go` files, so this file is
  NOT run as tests. It also has an open `// TODO: Add tests` list (bucket-full start,
  rejection, refill, sub-token, cap, retry-after, backward time).

## Commands

- Build: `go build ./...`
- Test: `go test ./...` and `go test -race ./...`
- Benchmarks: `go test -bench=BenchmarkAllow.*MemoryFixedWindow -benchmem ./...`

# V8 — RedisLimiter (ADR-0063 … ADR-0073)

> Back to [[Rate Limiter Learning Map]] · Prev [[V7 - Sharded MemoryStore]] · Next [[V9 - HTTP Middleware]] · Index [[README|ADR index]]

V8 builds the distributed path the ADR-0062 pivot reframed: a parallel `RedisLimiter` whose algorithms are Lua scripts running atomically on the server. It starts Redis-first with fixed window in Lua — strictly to learn `EVAL` cheaply before token bucket *forces* Lua (ADR-0063). `RedisLimiter` stays generic via a Redis-side algorithm interface that exposes `script()`/`args()` (ADR-0064); the I/O adapter returns raw values and the limiter owns parsing (ADR-0065); the caller owns the Redis client while the adapter stays private (ADR-0066); and the algorithm interface stays sealed (ADR-0067). Fixed-window policy config is shared across memory and Redis while execution stays separate (ADR-0068), backend-qualified naming lands (ADR-0069), and Redis token bucket gets a full design — scaled-integer state, Redis `TIME` + idle TTL, a uniform result contract, and a contained scaling boundary (ADR-0070 … ADR-0073). The Redis/Lua mechanics behind these decisions live in [[Redis Concepts]].

---

## ADR-0063 — Fixed window in Redis via Lua
**Status:** Accepted

**Context:** The first port needs cross-process read-modify-write atomicity — the distributed twin of the Get/Set race (ADR-0046) and bounded CAS (ADR-0051). Notably, fixed window *can* be done with primitives alone: `SET …NX EX window` + `DECR` is already safe (only the first `SET NX` wins; `DECR` is the atomic admission gate). The crash-prone `INCR`-then-`EXPIRE` count-up is the version that would need bundling.

**Decision:** Implement fixed window in Lua anyway — as a cheap place to learn `EVAL` before token bucket forces it, for a uniform path across all Redis algorithms, and to get pessimistic single-shot atomicity that removes the CAS retry loop (no `ErrCASConflict` on the Redis path).

**Consequences:** Algorithm logic duplicated in Go and Lua (the ADR-0062 cost). Lua scripts are harder to unit-test than pure `Decide` functions — they need a real or mocked Redis.

## ADR-0064 — RedisLimiter algorithm interface
**Status:** Accepted

**Context:** `RedisLimiter` must run different algorithms (fixed window → `limit, window`; token bucket → `limit, rate`) without coupling to any one parameter shape. Options ranged from named fields (couples to fixed window) to a raw `args []string` (loses type safety) to per-algorithm limiter types (duplicates EVAL plumbing).

**Decision:** `RedisLimiter` holds a Redis-side algorithm interface — `script() string`, `args() []string` — not algorithm params. Concrete `RedisFixedWindow` holds typed fields and converts them to `[]string` in `args()`. `Allow` passes `key` as `KEYS[1]` and `args()` as `ARGV`, never knowing which algorithm runs.

**Consequences:** Both decoupling *and* type safety — typed fields are compiler-checked at construction, and the unavoidable `[]string` conversion sits at the one boundary Lua forces it (Lua speaks only strings). Mirrors the in-process split, swapping `Decide(now, state)` for `script()`/`args()`. One extra interface+type per algorithm; the minimal `script`/`args` contract may need to grow (e.g. `keys()`) if an algorithm needs more than one key.

## ADR-0065 — Redis adapter returns raw result
**Status:** Accepted

**Context:** `RedisLimiter` runs Lua through a small `redisAdapter`. Should the adapter understand the script's `{allowed, remaining, retryAfterSeconds}` output shape?

**Decision:** `redisAdapter.eval` returns raw `any`. `RedisLimiter.Allow` owns parsing the Lua contract into `Result`.

**Consequences:** The adapter is Redis I/O plumbing, reusable for any script shape; the limiter — which already knows the key, the algorithm, the `Result` contract, and the promised output shape — is the natural place for type assertions. Better testability (mock raw Redis-like values). If future algorithms diverge in output shape, revisit an algorithm-owned `parseResult` — not yet.

## ADR-0066 — Redis client ownership boundary
**Status:** Accepted

**Context:** What should the public constructor accept — a `*redis.Client`, a config, or a public adapter interface?

**Decision:** The public constructor takes an existing Redis client and wraps it internally with the private `goRedisAdapter`. The adapter stays an internal seam, not user-facing API.

**Consequences:** Idiomatic Go — the application owns connection lifecycle, config, pooling, TLS/auth, and can share one client; the library owns rate-limiting behavior. Users never touch Lua/`EVAL`/`KEYS`/`ARGV`. The public API couples to a Redis client type (accepted). A future client library means another internal adapter and likely another constructor, but `Allow` shouldn't change. Don't prematurely export a generic `RedisExecutor` interface.

## ADR-0067 — Sealed Redis algorithms
**Status:** Accepted

**Context:** The `redisAlgorithm` interface (`script`/`args`) is unexported, so outsiders can't define their own.

**Decision:** Keep Redis algorithms sealed; callers use package-provided ones. Defer custom extension.

**Consequences:** A Redis algorithm isn't just a Go function — it pins Lua source, KEYS/ARGV conventions, the result tuple, atomicity assumptions, and script error behavior. Exporting would freeze all of that as public API while the Lua contract is still young. Lower flexibility, higher stability (mirrors ADR-0030 for the in-process side).

## ADR-0068 — Shared fixed-window config
**Status:** Accepted

**Context:** `MemoryFixedWindow` and `RedisFixedWindow` implement the same policy (limit, window) through different substrates. Should one type implement both `Decide` and `script`/`args`?

**Decision:** Share a private `fixedWindowConfig` (one constructor, validation once); keep execution separate across the two sibling types.

**Consequences:** Rule — *share stable policy, separate volatile execution*. The policy ("allow N per window") is stable; execution (Go state/CAS vs Lua/TTL/KEYS/ARGV) is volatile. Avoids duplicate validation without making the pure Go algorithm know Redis details. Adds one private config type.

## ADR-0069 — Backend-qualified naming
**Status:** Accepted

**Context:** Bare `Limiter`/`algorithm` became ambiguous once a Redis backend existed.

**Decision:** Qualifier-first names for backend-specific concepts: `MemoryLimiter`, `RedisLimiter`, `MemoryFixedWindow`, `RedisFixedWindow`, `MemoryTokenBucket`, `RedisTokenBucket`. Internal interfaces stay lowercase (`memoryAlgorithm`, `redisAlgorithm`). Policy config stays backend-neutral (`fixedWindowConfig`). Prefer `Memory` over `InMem`.

**Consequences:** Type names say where execution happens — clearer once multiple backends exist, at the cost of verbosity and a broad rename. Leaves room for a generic `Limiter` interface later without a name clash — deliberately deferred until both backends prove common behavior (delivered in V9 as the `Limiter` interface).

## ADR-0070 — Redis token bucket: scaled-integer state
**Status:** Accepted

**Context:** Redis token bucket must preserve fractional refill progress (e.g. 0.5 token from 250 ms at 2/s). Whole integers lose it; floats spread float behavior into Redis state.

**Decision:** Store scaled integer units — `tokenScale = 1000`, so 1 token = 1000 units — in a Redis hash (`tokens`, `last_refill_ms`). One request costs `tokenScale`; the admission check is `tokens >= tokenScale`.

**Consequences:** Preserves sub-token progress without floats in Redis. Slightly more memory (a hash, two fields) and more code that must track units (scaled tokens, ms, seconds). Lua owns the atomic read-refill-decide-write, so scaled arithmetic introduces no client-side races. (Opposite trade to the in-process float choice, ADR-0039.)

## ADR-0071 — Redis token bucket: time and TTL policy
**Status:** Accepted

**Context:** Distributed callers can disagree on the clock; idle keys need cleanup.

**Decision:** Use Redis `TIME` inside Lua (converted to ms) as the one authoritative clock. Use TTL only as idle cleanup, `ttl = ceil(capacity / refillRate)` (scale cancels), refreshed on every script call including rejections.

**Consequences:** All callers share one clock and one atomic transition. After a full-refill idle period an empty bucket would be full again, so expiring the key changes nothing; refreshing on every call stops an active-but-rejected caller from letting the key expire and recreating a fresh full burst. One `EXPIRE` per call, inside the same Lua execution (no extra round trip).

## ADR-0072 — Redis token bucket: result contract
**Status:** Accepted

**Context:** `RedisLimiter.Allow` parses Lua output into the public `Result`. Token bucket could return finer values than fixed window.

**Decision:** Keep the uniform Lua contract `{allowed, remaining, retryAfterSeconds}`. Token bucket floors `remaining` to whole tokens (`floor(tokens/tokenScale)`) and ceilings `retryAfterSeconds` (`ceil(deficit/refillRate)`).

**Consequences:** Uniform parsing across Redis algorithms; current learning value is atomic state transition, not API precision. Flooring preserves the `Result.Remaining int` contract and hides scaled units; ceiling avoids telling a caller to retry before a full token exists. Seconds are coarse for high-rate buckets — milliseconds can be revisited later.

## ADR-0073 — Redis token bucket: scaling boundary
**Status:** Accepted

**Context:** How far should scaled-integer arithmetic spread?

**Decision:** Keep scaling inside `RedisTokenBucket` and its Lua script, behind a package-level unexported `const tokenScale = 1000`. Caller and `RedisLimiter` see only domain concepts (capacity, refillRate, whole `remaining`, seconds).

**Consequences:** The algorithm-specific adapter is allowed to know how its script wants arguments; scaled units never leak into public constructors, `Result`, `RedisLimiter`, or callers. The conversion lives in two places (`args()` and Lua), so names must stay clear (`capacityScaled`, `refillUnitsPerSecond`, `tokenScale`).

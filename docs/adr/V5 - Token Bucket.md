# V5 — Token Bucket (ADR-0035 … ADR-0044)

> Back to [[Rate Limiter Learning Map]] · Prev [[V4 - Algorithm Abstraction]] · Next [[V6 - Storage Boundary]] · Index [[README|ADR index]]

V5 adds token bucket as the second algorithm — partly because the project needs a second limiter, but mostly to prove the V4 boundary holds against a genuinely different shape (ADR-0035). The model: capacity (burst) separate from refill rate (sustained throughput) (ADR-0036), state of `(tokens, last-refill)` (ADR-0037), refilled lazily on access with no background worker (ADR-0038), using `float64` token math (ADR-0039). The public `Result` stays integer-clean — whole `Remaining`, a `RetryAfter` derived from the token deficit (ADR-0040, ADR-0041) — while edge cases like backward clock movement (ADR-0042) and cold-start bursts (ADR-0043) get explicit rules. With two algorithms plus shared types now crowding one file, V5 also splits the package by responsibility (ADR-0044).

---

## ADR-0035 — Token bucket as the second algorithm
**Status:** Accepted

**Context:** Fixed window sits behind the algorithm boundary; a second algorithm validates it.

**Decision:** Add token bucket.

**Consequences:** Different config and state from fixed window, still O(1) memory and decision cost — a good stress test of the boundary. Brings higher model complexity: refill math, capacity capping, fractional tokens, retry timing, clock-regression handling.

## ADR-0036 — Capacity and refill rate
**Status:** Accepted

**Context:** Token bucket separates burst size from sustained rate.

**Decision:** Configure with `capacity int` (max immediate burst) and `refillRate float64` (tokens per second).

**Consequences:** Deliberately unlike fixed window's limit/window, confirming that algorithm-specific config belongs with the algorithm constructor (ADR-0031).

## ADR-0037 — Token bucket state shape
**Status:** Accepted

**Context:** Token bucket needs different runtime state than fixed window.

**Decision:** Store available tokens + last-refill time.

**Consequences:** Tokens capture spendable capacity; last-refill enables lazy computation. The limiter stores this opaquely through the marker interface (ADR-0028) and never inspects it.

## ADR-0038 — Lazy token refill
**Status:** Accepted

**Context:** A token bucket is conceptually filling continuously.

**Decision:** Refill lazily on request, computing current tokens from stored tokens, last-refill, now, and rate — no background worker.

**Consequences:** CPU cost scales with traffic, not with the number of known keys. Inactive keys aren't physically updated until they receive traffic again (see Go Concepts: *lazy state materialization*).

## ADR-0039 — Floating-point token arithmetic
**Status:** Accepted

**Context:** Refill produces fractional tokens (`3/s × 1.5s = 4.5`).

**Decision:** Use `float64` internally for tokens and refill math.

**Consequences:** Simple and preserves fractional progress. Float precision has edge cases; revisit fixed-point integers if exact accounting matters. (The Redis path makes the opposite choice with scaled integers — ADR-0070.)

## ADR-0040 — Whole-request `Remaining`
**Status:** Accepted

**Context:** `Result.Remaining` is an `int`, but internal tokens are fractional.

**Decision:** Truncate internal tokens to whole requests for `Remaining`.

**Consequences:** The public field answers "how many complete one-token requests can pass now" (e.g. `5.9` tokens → `Remaining = 5`). Callers don't see exact internal state.

## ADR-0041 — Token bucket `RetryAfter`
**Status:** Accepted

**Context:** A rejected request needs to know when one token will exist.

**Decision:** `RetryAfter = (1 − availableTokens) / refillRate`.

**Consequences:** Preserves the public meaning of `RetryAfter` (minimum wait before retry). Internally differs from fixed window: token bucket waits for one token, fixed window waits for window reset.

## ADR-0042 — Clamp negative elapsed time
**Status:** Accepted

**Context:** Clock skew can make `now` appear earlier than the stored last-refill time.

**Decision:** Clamp negative elapsed to zero; if `now < lastRefill`, leave `lastRefill` unchanged.

**Consequences:** Resilient to clock regression — never subtracts tokens or misbehaves on backward time. Hides the event unless metrics are added later.

## ADR-0043 — Token bucket starts full
**Status:** Accepted

**Context:** A new key has no token-bucket state.

**Decision:** Initialize new keys with a full bucket, then consume one token for the first accepted request.

**Consequences:** New keys can burst up to capacity — expected token-bucket behavior, but possibly too permissive for abuse-sensitive flows.

## ADR-0044 — Split files by responsibility
**Status:** Accepted

**Context:** Fixed window, token bucket, shared result types, the algorithm contract, and sentinel errors now overflow the single file (resolving ADR-0015).

**Decision:** Split into responsibility-named files within one `ratelimiter` package: `types.go` (small shared contracts), `limiter.go` (orchestration), `fixed_window.go`, `token_bucket.go`, `errors.go`, `state_store.go`, `memory_store.go`.

**Consequences:** Concrete behavior sits next to the type that owns it; only small shared contracts are centralized. Watch that `types.go` doesn't become a junk drawer — keep concrete algorithm types out of it (file boundaries are not visibility boundaries; see Go Concepts: *package scope across files*).

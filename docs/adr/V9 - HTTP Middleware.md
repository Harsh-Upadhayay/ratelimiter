# V9 — HTTP Middleware (ADR-0074 … ADR-0085)

> Back to [[Rate Limiter Learning Map]] · Prev [[V8 - RedisLimiter]] · Index [[README|ADR index]]

V9 puts the limiter into a real `net/http` request path (ADR-0074). The middleware extracts a key via a caller-provided `KeyFunc`, calls the limiter, and translates the outcome into status codes and headers only — never a body (ADR-0076, ADR-0077). Failure of the limiter infrastructure is a configurable, whole-middleware policy distinct from rate-limit rejection: `429` means "over the limit", `503` means "limiter failed and we chose fail-closed", and fail-open is a silent pass-through (ADR-0075, ADR-0081). Required dependencies are explicit constructor params; only optional settings use functional options (ADR-0078, ADR-0080), and the type is named for behavior (`Middleware`), not its Redis backend (ADR-0079). `Retry-After` is converted from a `time.Duration` to ceiling seconds only at the HTTP edge (ADR-0082). The back half of V9 makes `MemoryLimiter` satisfy the same `Limiter` interface directly via a private clock (ADR-0083) — which forces benchmarks onto the production clock (ADR-0084) and, in doing so, uncovers real hot-key CAS exhaustion (ADR-0085). Checkpoints C09 and C10 capture the end state.

---

## ADR-0074 — HTTP middleware boundary
**Status:** Accepted

**Context:** After the memory and Redis paths, the next learning step is putting the limiter in an HTTP request path (more tests and another algorithm were deliberately skipped).

**Decision:** Build a `net/http` middleware: request → extract key → call limiter → pass through or write an HTTP status.

**Consequences:** Teaches middleware shape, request-scoped context, identity extraction, status codes, rate-limit headers, and failure policy without repeating algorithm work. The middleware must distinguish rate-limit rejection from limiter infrastructure failure (ADR-0075).

## ADR-0075 — Middleware failure policy
**Status:** Accepted

**Context:** In the request path, limiter infrastructure (e.g. Redis) can fail.

**Decision:** A whole-middleware `FailurePolicy` enum: `FailOpen` (limiter error → allow) or `FailClosed` (limiter error → `503`), validated in the constructor. Default `FailOpen` — rate limiting is usually protective middleware, not the core feature.

**Consequences:** `429` (healthy limiter, over limit) and `503` (limiter failed, fail-closed) mean different things. Per-instance, not per-request — different route groups get different middleware instances. Fail-open protects availability; fail-closed protects enforcement.

## ADR-0076 — Caller-provided key function
**Status:** Accepted

**Context:** The middleware can't know the app's identity model (IP, user, API key, tenant, route+user).

**Decision:** `type KeyFunc func(*http.Request) string`, called per request; an empty key → `400 Bad Request`. Because it's mandatory, `KeyFunc` is a direct constructor parameter, not an option.

**Consequences:** Identity stays application-owned; one middleware supports many key schemes. The caller is responsible for ensuring the key source is trustworthy (Go Concepts: *function types as callbacks*).

## ADR-0077 — Rate-limit HTTP headers
**Status:** Accepted

**Context:** `Result` already carries `Remaining` and `RetryAfter`.

**Decision:** Write status and headers only, no body. Set `X-RateLimit-Remaining` after every decision; on rejection return `429 Too Many Requests` with `Retry-After` in ceiling seconds (conversion owned by the HTTP boundary; `Result.RetryAfter` stays a `time.Duration`).

**Consequences:** Headers are useful to clients; bodies are application-specific, so the middleware imposes no JSON/text format. Small `int`/`Duration`→string conversions (ADR-0082).

## ADR-0078 — Functional options for the middleware
**Status:** Accepted

**Context:** Two required dependencies (limiter, key func) and one optional setting (failure policy).

**Decision:** Functional options (`type MiddlewareOption func(*Middleware)`) only for optional config; required values stay explicit constructor params. Options may mutate the partially built `*Middleware` directly since runtime fields and config fields coincide. The constructor validates after applying options.

**Consequences:** `NewMiddleware(limiter, keyFunc, WithFailurePolicy(FailClosed))`. Required behavior is visible at the call site; optional tuning stays named and extensible. Avoids a duplicate config struct (Go Concepts: *functional options on runtime structs*).

## ADR-0079 — Behavior-named middleware
**Status:** Accepted

**Context:** Early drafts used `RedisMiddleware` because the backing limiter is Redis-first.

**Decision:** Name for behavior, not backend. Since the package is `ratelimiter`, use `ratelimiter.Middleware` / `NewMiddleware` (no stutter), not `RedisMiddleware` or `RateLimitingMiddleware`.

**Consequences:** The middleware's public behavior is HTTP rate limiting; Redis is an implementation detail of its limiter dependency. The name doesn't force a fully generic limiter interface immediately, and won't need to change if the backend does. Keep docs honest so a generic name doesn't overpromise.

## ADR-0080 — Required dependencies outside functional options
**Status:** Accepted

**Context:** An early shape put `KeyFunc` behind `WithKeyFunc`.

**Decision:** Required deps stay out of options: `NewMiddleware(limiter Limiter, keyFunc KeyFunc, opts ...MiddlewareOption)`. Options are reserved for optional config like `WithFailurePolicy`.

**Consequences:** You can't build useful middleware without both a limiter and a key func, so hiding them in optional-looking `With...` calls would make the API dishonest. Required args are visible and compiler-enforced; slightly more positional arguments. (ADR-0083 later reverses applying this reasoning to the *clock*, which is not a caller-configured dependency.)

## ADR-0081 — Fail-open is silent pass-through
**Status:** Accepted

**Context:** On limiter failure under `FailOpen`, should the downstream handler be told?

**Decision:** `FailOpen` = silent pass-through; call `next` as if allowed, with no signal to the handler.

**Consequences:** The common meaning of fail-open — availability over enforcement, no extra context coupling. A "tell the handler" variant (`DelegateOnLimiterError`) is a different policy requiring an explicit context signal; deferred until route-specific business decisions need it (login fail-closed, public feed fail-open, expensive endpoint app-owned fallback). Observability should come from logs/metrics later, not the protected handler.

## ADR-0082 — HTTP Retry-After converts Duration to seconds
**Status:** Accepted

**Context:** Every limiter path uses `Result.RetryAfter time.Duration`; HTTP `Retry-After` as a delta uses seconds.

**Decision:** Keep `RetryAfter` a `time.Duration` everywhere in Go; convert to ceiling whole seconds only at the HTTP boundary (`≤0 → "0"`, `(0,1s] → "1"`, `1500ms → "2"`).

**Consequences:** The domain API never leaks HTTP formatting; the middleware is the one boundary that knows HTTP needs a string. All paths (memory fixed window/token bucket, both Redis algorithms) normalize to `time.Duration` before the boundary, so the middleware sees one contract. Ceiling avoids "retry immediately" when a sub-second wait remains (Go Concepts: *ceiling duration conversion*).

## ADR-0083 — MemoryLimiter implements Limiter via a private clock
**Status:** Accepted — supersedes ADR-0005

**Context:** `MemoryLimiter.Allow(key, now)` took caller time (kept algorithms pure, let tests time-travel), but the `Limiter` interface the middleware needs is `Allow(ctx context.Context, key string) (Result, error)` — no `now`. C09 had deferred this as a "clock adapter".

**Decision:** Make `MemoryLimiter` satisfy `Limiter` *directly* — no adapter. Replace the `now` parameter with an injectable but **private**, non-caller-configurable `clock` interface (`now() time.Time`), defaulting to `realClock{}`. `NewMemoryLimiter(algo, store)` is unchanged. `Allow` honors `ctx` with a single top-of-call `ctx.Err()` check and samples `clock.now()` per CAS attempt. White-box tests inject a `*testClock` via the field.

**Consequences:** This reverses applying ADR-0078/ADR-0080 to the clock: a functional option is for caller-configured settings; the clock is a pure internal time source whose only non-production use is testing. The rule is now uniform — a limiter owns its time source (symmetric with `RedisLimiter` owning Redis `TIME`); callers don't inject it, and "I need to test it" reaches the private seam, not the public API. Compile-time assertions `var _ Limiter = (*MemoryLimiter)(nil)` / `(*RedisLimiter)(nil)` pin the contract (Go Concepts: *compile-time interface assertion*, *interface narrows the method set*).

## ADR-0084 — Benchmarks use the production clock
**Status:** Accepted

**Context:** After ADR-0083, `Allow` reads `clock.now()` every call. The shared `*testClock` is mutex-guarded (it must be — concurrent tests read it). Routing benchmarks through the test-clock helper meant every `Allow` took that one mutex.

**Decision:** Benchmarks build limiters through the production constructor (defaults to lock-free `realClock`), not `newTestMemoryLimiter`. Benchmarks never advance time, so they have no reason to hold the controllable clock.

**Consequences:** A shared test-clock mutex is a global serialization point that flattened `ManyKeyParallel` toward `SameKeyParallel`, erasing the sharding win as a measurement artifact. Using the real clock restores the true sharding gap — and exposes real single-key CAS contention (ADR-0085). A benchmark must measure the production path; injecting a test-only sync primitive corrupts it.

## ADR-0085 — Hot-key CAS exhaustion is expected
**Status:** Accepted

**Context:** Removing the test-clock mutex (ADR-0084) made same-key parallel benchmarks immediately fail with `ErrCASConflict` — the fake clock had been masking the limiter's true behavior under a hot key.

**Decision:** Accept it as expected. The bounded 10-attempt CAS loop stays; under sustained single-key contention, `GOMAXPROCS` goroutines fight over one record's version and an unlucky one can lose all 10 attempts. `ErrCASConflict` is a truthful "this key is too hot for optimistic CAS" signal, not a defect. Benchmarks treat it as non-fatal (`errors.Is(err, ErrCASConflict)`).

**Consequences:** Sharding doesn't help (same key → same shard → same record); this is the classic hot-key problem (plan step 6), intrinsic to optimistic concurrency under high single-key write contention. Mitigation deferred — candidates: backoff/jitter, higher/unbounded retry bound, per-key serialization for hot keys, or routing hot keys to the Redis path (single-shot Lua has no CAS retry — ADR-0063). The broader lesson: test doubles can hide the very properties a benchmark exists to reveal.

---

## Checkpoint C09 — V9 HTTP middleware
V9 is structurally complete. Public shape: `ratelimiter.Middleware` / `NewMiddleware`, a `Limiter` interface (`Allow(ctx, key) (Result, error)`), `KeyFunc`, `FailurePolicy` (`FailOpen`/`FailClosed`), and `WithFailurePolicy`. Runtime behavior: empty key → `400`; limiter error + `FailOpen` → call next silently; limiter error + `FailClosed` → `503`; allowed → set `X-RateLimit-Remaining`, call next; rejected → set `X-RateLimit-Remaining` + `Retry-After`, return `429`; `Wrap(nil)` panics. Status/headers only, no body. `go test ./...` passed at checkpoint. Deferred: middleware behavior tests, observability, a local HTTP example server, sliding-window algorithms, delegate-on-error, Redis integration/runtime tests, and (then) the memory adapter for `Limiter`.

## Checkpoint C10 — MemoryLimiter implements Limiter
Complete. `MemoryLimiter` satisfies `Limiter` directly (no adapter), so it drops into `NewMiddleware` exactly like `RedisLimiter`; build, tests, `-race`, and benchmarks all green. Changes: `Allow(key, now)` → `Allow(ctx, key)`; caller `now` → private `clock` (default `realClock`, not caller-configurable — ADR-0083); a `clock.go` with the unexported `clock` interface + `realClock`; compile-time `var _ Limiter` assertions in `limiter.go`; tests on a white-box advanceable `*testClock`; benchmarks on the lock-free production clock (ADR-0084). This closed the C09 "memory adapter" branch as a direct refactor and uncovered hot-key CAS exhaustion (ADR-0085). Verified with `go build`, `go test`, `go test -race`, and `go test -bench BenchmarkAllow`. Still deferred: middleware tests, observability, example server, sliding-window algorithms, delegate-on-error, Redis integration tests, hot-key mitigation.

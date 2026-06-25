# Go Concepts

> Back to [[Rate Limiter Learning Map]] · See also [[Redis Concepts]] · Decisions in [[README|ADR index]]

The Go language and style ideas the rate limiter exercises, grouped by theme and condensed to the point. Each entry names the decision(s) that put it to use.

## State and data modeling

**Structs for grouped state.** Use a struct when several fields are one concept — request count and window start aren't independent, they *are* a key's window state. A named struct makes code communicate domain meaning instead of passing loose values around (ADR-0002).

**Maps, comma-ok, and value copies.** `v, ok := m[k]` distinguishes a missing key from a zero value. Crucially, reading a struct *value* from a map returns a **copy** — mutating it does nothing unless you assign it back. That shapes the limiter's read-decide-write path (ADR-0002, ADR-0006).

**Exported types with unexported fields.** Visibility is per-identifier (capitalization), not per-type: an exported `*Limiter` can still hide `limit`, `window`, `states`, `mu`. Returning a pointer does **not** grant outside packages access to lowercase fields. Expose the stable API; hide representation (ADR-0012, ADR-0017).

**Exported result structs.** If a public function returns a struct callers must read, export both the type *and* the inspected fields — an exported type with unexported fields would be useless to read. So `Result{Allowed, Remaining, RetryAfter}` is fully exported (ADR-0021).

## Methods and mutation

**Pointer receivers.** Use a pointer receiver when a method mutates the receiver or copying it would mislead/cost. The limiter mutates its own state, so `Allow` operates on the limiter itself — which makes shared mutation, and therefore a concurrency story, explicit (ADR-0003).

**Method receivers.** Go declares methods *outside* the struct body via a receiver before the name (unlike C++'s in-class methods). The receiver is the object the method attaches to; choose pointer vs value per mutation/copy cost (ADR-0017).

**Method sets and interface satisfaction.** A type satisfies an interface through its **method set**, and the receiver decides the set: a value-receiver method belongs to both `T` and `*T`; a pointer-receiver method belongs to **only** `*T`. So with `func (x *T) Allow(...)`, `var _ I = &T{}` compiles but `var _ I = T{}` does not. Accept interfaces *by value* (`limiter Limiter`, never `*Limiter` — a pointer-to-interface is almost always wrong) (ADR-0079, ADR-0083).

**Do not copy mutexes.** A `sync.Mutex` is synchronization bookkeeping; copying it duplicates the lock value, not the relationship. Copy a struct holding a mutex and a map and you can get two mutexes guarding one map — or a copy stuck locked forever. Prefer pointers for stateful mutex-bearing types (ADR-0017, ADR-0056).

## Concurrency

**Mutexes and critical sections.** `sync.Mutex` guards shared mutable state across goroutines. The critical section is the *whole* read-decide-write sequence, not just the map read — otherwise two goroutines read the same old count and both wrongly allow (ADR-0004).

**Defer-unlock pattern.** With multiple return branches, `defer mu.Unlock()` right after locking makes release function-scoped and unmissable. Tiny overhead; holds the lock until return — a fine clarity trade until a path is proven hot (ADR-0009).

**Optimistic concurrency with CAS.** Compare-and-swap commits a new value only if the stored version still matches the one read. It protects read-decide-write across concurrent updaters *without* a long-held lock, at the cost of retry loops and conflict handling (ADR-0047, ADR-0051, ADR-0085).

**Store-owned mutexes.** A mutex must live on the shared object it protects; a fresh local mutex per call protects nothing. `MemoryStore` owns the map, so it owns the mutex, and uses pointer receivers so copies don't fork the lock (ADR-0049).

**Race detector.** `go test -race` flags unsynchronized concurrent memory access. It finds *data* races, not *logical* ones — it can pass while a stale read-decide-write is still wrong; that's what CAS exists to protect (ADR-0048).

**Lock striping.** Split one big critical section into many: `key → shard → shard mutex`. Same key still serializes; different keys proceed in parallel when they land on different shards. More shards cut contention but add memory and hashing, and two hot keys can still collide on one shard (ADR-0053).

**Key hashing.** Hash a key to a stable number, then `hash(key) % shardCount` picks a shard. The *same* key must always map to the same shard, or `Get` and `CAS` would touch different records. Distribution quality drives contention (ADR-0053, ADR-0055).

**Composition with pointer fields.** Structs compose by holding other concrete types as fields; when the field contains a mutex, hold a **pointer** (`[]*MemoryStore`) so values aren't copied and synchronization stays intact (ADR-0056).

**String to byte slice.** `[]byte(key)` is a type conversion (not a call) that copies a string's bytes into a mutable slice — needed because hashers and `io.Writer` take `[]byte`, and Go strings are immutable (ADR-0053).

**Unsigned modulo safety.** `%` on a *signed* int can be negative in Go — an invalid index. Hashers return `uint32`; keep the modulo unsigned: `int(h.Sum32() % uint32(len(shards)))`. Casting to `int` *before* `%` risks a negative result when the high bit is set (ADR-0058).

## Time and control flow

**Time and duration boundaries.** `time.Time` for instants, `time.Duration` for spans. A half-open window includes its start and excludes its end, so the exact expiry instant belongs to the next window. Limiter bugs cluster at boundaries; explicit time makes them testable (ADR-0008).

**Duration → seconds conversion.** `d.Seconds()` yields a `float64`; token-bucket refill uses `elapsed.Seconds() * refillRate`. The public config speaks tokens-per-second while `Duration` stores nanoseconds internally (ADR-0036, ADR-0039).

**Ceiling duration conversion.** `int(d)` is raw nanoseconds, not seconds; `int(d.Seconds())` *truncates* (1500ms → 1, 500ms → 0). For `Retry-After` you want ceiling whole seconds: clamp `d ≤ 0 → 0`, else `(d + time.Second - 1) / time.Second`. Truncating would tell a client to retry before a token/window is ready (ADR-0082).

**Lazy state materialization.** Store enough to *compute* current state on access rather than updating continuously in the background. Token bucket computes "tokens now" from stored tokens, last-refill, and rate — no goroutines, no per-key timers, no work for idle keys. State may be stale between requests, but the decision is correct when read (ADR-0038).

**Package scope across files.** Files in one directory sharing a package name share scope, so splitting into `memory_limiter.go`, `fixed_window.go`, etc. needs no new exports. File boundaries organize for humans; they are not visibility boundaries (ADR-0044).

## Errors and API shape

**Multiple return values.** Go's `result, err := op()` convention: check the error, then use the result. For `Allow`, a rate-limit rejection is *not* an error — it's a successful decision with `Allowed: false`; only invalid input is an error (ADR-0019).

**Sentinel errors.** A package-level error value for a stable category, comparable with `errors.Is`. Good for known public branches like `ErrEmptyKey`; exported sentinels are public API, so don't export errors for temporary or overly specific internal failures (ADR-0022).

**Pure helper functions.** A deterministic helper takes all inputs as parameters and has no side effects — no map reads, no locking, no `time.Now`, no package state. The fixed-window helper only computes `(result, nextState)`, returning unchanged state on rejection. Pure helpers are easy to reason about and make later interfaces easy to discover (ADR-0023, ADR-0025).

**Public-API tests.** Same-package tests can still choose to exercise only exported behavior, keeping refactors safe — private internals change without breaking tests, as long as outcomes are observable through the public API (ADR-0024).

**Helper parameter ordering.** With no named arguments, parameter order *is* readability. Group by meaning — time, then config, then state, then metadata. Many params is a smell, but a config struct here would be premature (ADR-0023).

**Interfaces from real variation.** Interfaces are strongest when drawn from real behavioral variation. With one implementation an interface is a guess; with two (fixed window, token bucket) the shared contract and the differing details become clear (ADR-0013, ADR-0026).

**Structural interface satisfaction.** No `implements` keyword — a type satisfies an interface by having the methods. The limiter depends on algorithm *behavior* without depending on a concrete type. Keep interfaces small so their purpose stays legible (ADR-0026).

**Compile-time interface assertion.** `var _ Limiter = (*MemoryLimiter)(nil)` checks satisfaction at build time, zero runtime cost: `_` discards the variable, the explicit type forces assignability, and `(*T)(nil)` is a typed-nil conversion checked on the *type* (no deref). Use the pointer form when `Allow` has a pointer receiver. It pins the contract at the definition site so an incompatible change breaks the build *here*, not at some far call site — like `var _ io.Writer = (*bytes.Buffer)(nil)` (ADR-0083).

**Interface narrows the method set.** A value held through an interface exposes only the interface's methods, even if the concrete type has more. `MemoryLimiter.clock` (static type `clock`, just `now()`) can't reach `*testClock`'s `Advance`/`Set`. So the test helper *returns the concrete `*testClock`* rather than type-asserting. Design point: narrow interface for the consumer (the limiter only *reads* time), richer concrete type for the controller (the test *controls* it) (ADR-0083).

**Marker interfaces and opaque state.** A marker interface uses an (unexported) method only to identify a controlled set of valid types. It lets the limiter store algorithm state opaquely — never inspecting fields — while avoiding raw `any`. Adds ceremony and doesn't remove the per-algorithm type assertion (ADR-0028).

**Type assertions.** `v, ok := x.(T)` retrieves a concrete type from an interface value — the safe comma-ok form yields zero value + `false` on mismatch; the one-value form panics. The limiter treats state as opaque but each algorithm asserts its own concrete state; the Redis adapter likewise asserts go-redis's `any` (a Lua `{1,4,0}` arrives as `[]interface{}{int64,…}`). At package boundaries prefer comma-ok and return an error; reserve the panic form for impossible internal invariants (ADR-0028, ADR-0033, ADR-0065).

**Constructor validation ownership.** Constructors reject invalid config before returning a usable object, and validation belongs with the type that owns the invariant — algorithm constructors validate algorithm config; the limiter constructor validates limiter-level deps. Avoids repeated hot-path checks and prevents invalid components existing (ADR-0031).

**Nil interface guard.** A constructor taking an interface dependency should reject a missing one before storing it, so the hot path doesn't panic. Mind Go's typed-nil subtlety: a plain nil check catches a nil interface value but not a non-nil interface boxing a typed-nil pointer (ADR-0031).

**Exported interfaces with unexported types.** An exported interface can still be unimplementable from outside if its methods mention unexported types — `StateStore` is exported but its methods use `algorithmState`. That makes the boundary visible for learning while keeping state controlled inside the package; not a polished public boundary (ADR-0052).

**Iota enum pattern.** Go has no native enums; use a named int type plus `const ( FailOpen FailurePolicy = iota; FailClosed )`. Reads far better than a bare `bool` at the call site. Named int types aren't closed sets — a caller can still write `FailurePolicy(99)` — so validate enum-like config (ADR-0075).

**Function types as callbacks.** Functions are values, so a package can define `type KeyFunc func(*http.Request) string`, store it, and call it later. The middleware doesn't know the app's identity model, so the app supplies `request → KeyFunc → key`. The app owns identity; the middleware owns enforcement (ADR-0076).

**Functional options.** A constructor pattern for *optional* config: `NewThing(required, WithTimeout(t))` where `type Option func(*config)`. Best for optional settings and named overrides; required deps stay explicit positional params so they're visible. Defaults live in one place; adding options doesn't change the signature (ADR-0078, ADR-0080).

**Functional options on runtime structs.** When the runtime object *is* the config (the middleware's fields are its construction fields), options can mutate the final object directly — `type MiddlewareOption func(*Middleware)` — instead of a duplicate config struct. Safe because the constructor still builds → applies options → validates → only then returns; callers never see an unvalidated object (ADR-0078).

**Unexported package constants.** A package-scope `const tokenScale = 1000` is shared by every file in the package but hidden from callers (lowercase). It avoids duplicating the literal across the token-bucket constructor, arg conversion, and script meaning, while keeping the implementation detail off the public API (ADR-0073).

## HTTP integration

**net/http middleware pattern.** The server is built on `Handler` (`ServeHTTP(w, r)`); middleware is usually `func(next http.Handler) http.Handler` — a wrapper that inspects the request and either stops or calls `next.ServeHTTP`. The rate limiter sits before the app handler: allowed → call next; rejected → write status/headers and stop. Use middleware for cross-cutting concerns (auth, logging, tracing, rate limiting) (ADR-0074, ADR-0079).

**http.HandlerFunc adapter.** `HandlerFunc` adapts a plain `func(w, r)` into a `Handler` via its own `ServeHTTP` that just calls the function. `Wrap(next)` runs once at setup and returns `http.HandlerFunc(func(w, r){…})` that runs per request — clean mental split between build-time and request-time (ADR-0074).

**ResponseWriter headers.** Set headers with `w.Header().Set(k, v)` *before* `w.WriteHeader(status)`; once the body is written Go implicitly sends `200`, so a rejecting middleware must write its status explicitly first. Here: set `X-RateLimit-Remaining` and `Retry-After`, write `429`, no body (ADR-0077).

**context.Context.** A value threaded down a call stack carrying cancellation, deadlines, and (sparingly) request-scoped values, so the top of the chain can tell the bottom "stop waiting" — vital for blocking I/O like a Redis round-trip. Conventions: `ctx` is the first parameter; **never store it in a struct** (it models one operation's lifetime, not an object's). `RedisLimiter.Allow(ctx, key)` passes the caller's `ctx` straight into go-redis's `Eval` rather than inventing its own `Background()` (ADR-0066, ADR-0083).

## Testing and benchmarking

**Go benchmarks.** `Benchmark...(b *testing.B)` runs a `b.N` loop the tool sizes for stability. Put setup before `b.ResetTimer`, precompute keys, avoid incidental allocations, use `b.RunParallel` for concurrency with per-worker local counters to avoid measurement noise. The limiter's same-key vs many-key benchmarks revealed the global-mutex bottleneck (ADR-0048).

**Test helper pattern.** A shared setup/assertion helper taking `*testing.T` must call `t.Helper()` so failures point at the caller. Keep it unexported, let it own error handling via `t.Fatal` (return the value directly, not `(value, error)`), and have it build a *fresh* instance each call to avoid cross-test pollution (ADR-0060, ADR-0061).

**Contract testing via interface parameter.** When several types implement one interface, write one helper taking the interface and call it once per implementation. No duplicated assertions; a new backend is one new `Test...`. The helper must use *only* the interface — no concrete type assertions, or it tests implementation, not contract (ADR-0060).

**go test compiles the whole package.** The compilation unit is the package, not the file: `go test` builds *all* `*_test.go` in the package into one binary, then `-run` filters *execution* (an unanchored substring regex — anchor with `^...$`). A stale sibling test file fails the whole build → zero tests run. You can't run one file's tests; select by function name or package only. Escape hatch: temporarily rename an unmigrated file off the `_test.go` suffix (ADR-0083, ADR-0084).

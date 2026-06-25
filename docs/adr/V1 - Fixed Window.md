# V1 — Fixed Window (ADR-0001 … ADR-0018)

> Back to [[Rate Limiter Learning Map]] · Next [[V2 - API Evolution]] · Index [[README|ADR index]]

V1 builds the smallest useful rate limiter — a fixed-window counter — to learn Go state, maps, time, and locking before any abstraction. The whole thing lives in one file (ADR-0015) with interfaces deliberately postponed (ADR-0013): a `Limiter` struct owns the config, a `map[string]userState`, and one global mutex (ADR-0004, ADR-0016). Time is passed in explicitly so boundary cases are testable (ADR-0005), the window is half-open and resets at `now` (ADR-0007, ADR-0008), and the API stays deliberately thin — a boolean result, panics for bad input — because the goal is correct state transitions first, not a polished surface (ADR-0010, ADR-0011, ADR-0018). Several of these starter choices are intentionally provisional and get superseded as soon as real pressure appears.

---

## ADR-0001 — Fixed window first
**Status:** Accepted

**Context:** The target design spans token bucket and two sliding-window variants, but V1 needs something small enough to learn Go state, maps, time, and locking.

**Decision:** Start with a fixed-window limiter — one count and one window start per key.

**Consequences:** Smallest possible correct limiter. Accepts boundary bursts (a caller can spend quota at the end of one window and again at the start of the next); that flaw is the motivation to compare against sliding window and token bucket later.

## ADR-0002 — Per-key state as map of structs
**Status:** Accepted

**Context:** Each key needs enough state to decide the next request within the current window.

**Decision:** Keep a `map[string]userState`, where `userState` groups request count and window start.

**Consequences:** Fields that belong together travel together. Map struct values are copied on read, so the update path must read, decide, then write the value back — a recurring Go gotcha (see Go Concepts: *maps, comma-ok and value copies*).

## ADR-0003 — Pointer receiver for the mutating limiter
**Status:** Accepted

**Context:** The limiter owns mutable state: config, the per-key map, and a mutex.

**Decision:** Methods that mutate limiter state use pointer receivers.

**Consequences:** The method operates on the limiter itself, not a copy. Shared mutation becomes explicit, which forces clear locking rules.

## ADR-0004 — One global mutex
**Status:** Superseded by ADR-0045, ADR-0053

**Context:** Concurrent requests can read the same count, both decide "allow", and overwrite each other.

**Decision:** Protect the whole read-decide-write sequence with a single mutex.

**Consequences:** Simplest correct protection, but requests for *different* keys serialize behind one lock. That bottleneck — confirmed by benchmark in ADR-0048 — is exactly what the storage boundary and sharding later remove.

## ADR-0005 — Explicit time input
**Status:** Superseded by ADR-0083

**Context:** The limiter needs the current time to decide whether a window expired.

**Decision:** Pass `now` into the decision rather than reading the wall clock internally.

**Consequences:** Tests become deterministic and boundary instants are exactly reproducible. Call sites are slightly less ergonomic. (V9's `Limiter` interface has no `now`; ADR-0083 replaces this with a private injectable clock.)

## ADR-0006 — Count means accepted requests
**Status:** Accepted

**Context:** The stored count needs a stable meaning between requests.

**Decision:** Count = accepted requests in the current window (state at rest, not a "before this request" view).

**Consequences:** A policy that should reject without consuming quota must run before this limiter.

## ADR-0007 — Window resets at now
**Status:** Accepted

**Context:** When a request arrives after the previous window expired, the limiter must choose the next window start.

**Decision:** Reset window start to `now`.

**Consequences:** Simple per-key reasoning and clean idle-gap handling. Windows are not globally aligned — this is a per-key fixed window, not a wall-clock one. Alternatives (`oldStart + window`, rounded boundaries) are noted for the eventual sliding-window comparison.

## ADR-0008 — Half-open window interval
**Status:** Accepted

**Context:** Expiry needs exact boundary behavior.

**Decision:** Use `[WindowStart, WindowStart + window)` — the end instant belongs to the next window.

**Consequences:** Adjacent windows never overlap. Boundary behavior ("just before", "exactly at", "after" expiry) must be tested explicitly.

## ADR-0009 — Defer unlock immediately after lock
**Status:** Accepted

**Context:** The decision has several early-return branches.

**Decision:** `defer mu.Unlock()` right after locking.

**Consequences:** Every return path releases the lock. Tiny `defer` overhead and the lock is held until return — an acceptable clarity trade for V1; revisit only with evidence the path is hot.

## ADR-0010 — Boolean return for V1
**Status:** Superseded by ADR-0019, ADR-0021

**Context:** The decision could return just allow/reject, or richer operational detail.

**Decision:** Return only a boolean for V1.

**Consequences:** Focuses V1 on correct state transitions. A boolean can't feed HTTP headers or a precise `Retry-After`; the richer `Result` arrives in V2.

## ADR-0011 — Panic on invalid config
**Status:** Superseded by ADR-0031

**Context:** The limiter needs valid config before it can decide correctly.

**Decision:** Panic in the constructor on an invalid limit or window.

**Consequences:** Keeps V1 tiny and fails fast, but panic is harsh for library code. Algorithm constructors return typed errors instead from V4 onward.

## ADR-0012 — Exported API, hidden state
**Status:** Accepted

**Context:** Go controls visibility by capitalization.

**Decision:** Export `Limiter`, `NewLimiter`, `Allow`; keep per-key state unexported.

**Consequences:** Callers get behavior, not representation, leaving the internal state free to change as algorithms and stores evolve. External tests must go through the public API.

## ADR-0013 — Delay interfaces
**Status:** Accepted

**Context:** The target design has `Algorithm` and `StateStore` boundaries, but V1 has exactly one of each.

**Decision:** Do not introduce those interfaces yet.

**Consequences:** Interfaces describe known variation; with one implementation they would be a guess. Accepts later refactoring as the cost of not over-abstracting — interfaces arrive in ADR-0026 and ADR-0045 once a second implementation reveals the real contract.

## ADR-0014 — Package name `ratelimiter`
**Status:** Accepted

**Context:** The package name sets the import/usage style.

**Decision:** Use `ratelimiter`.

**Consequences:** Explicit and reusable. `ratelimiter.Limiter` slightly stutters but is clear; this naming context is later leaned on to avoid stutter in `Middleware` (ADR-0079).

## ADR-0015 — One file first
**Status:** Superseded by ADR-0044

**Context:** V1 is intentionally small.

**Decision:** Keep everything in one file.

**Consequences:** Makes real friction visible before splitting. The crowding becomes the signal to split by responsibility once algorithms, stores, and errors accumulate.

## ADR-0016 — Limiter owns runtime state
**Status:** Accepted

**Context:** V1 needs an owner for config, the per-key map, and its lock.

**Decision:** An exported `Limiter` struct owns limit, window, the state map, and the mutex.

**Consequences:** Multiple independent limiter instances become possible (e.g. one policy for `/search`, another for `/login`). More structure than package globals, but far easier to test.

## ADR-0017 — Constructor returns a pointer
**Status:** Accepted

**Context:** The limiter is stateful and contains a `sync.Mutex`.

**Decision:** `NewLimiter` returns `*Limiter`.

**Consequences:** Goroutines share one instance and coordinate through the same mutex and map. Copying a mutex-containing value would be a bug (see Go Concepts: *do not copy mutexes*).

## ADR-0018 — Panic on empty key
**Status:** Superseded by ADR-0022

**Context:** The key is the map key for per-key state.

**Decision:** Panic in `Allow` on an empty key.

**Consequences:** Treats an empty key as a caller bug and fails fast, consistent with ADR-0011. Harsh for library code; V2 replaces it with the `ErrEmptyKey` sentinel.

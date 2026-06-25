# V4 — Algorithm Abstraction (ADR-0026 … ADR-0034)

> Back to [[Rate Limiter Learning Map]] · Prev [[V3 - Decision Logic]] · Next [[V5 - Token Bucket]] · Index [[README|ADR index]]

V4 is where the delayed interface (ADR-0013) finally earns its keep: token bucket is coming, and it has different config and state from fixed window, so a real `algorithm` boundary now describes real variation (ADR-0026). Each algorithm owns its own state (ADR-0027), stored opaquely through a marker interface so the limiter never inspects algorithm-specific fields (ADR-0028). The interface stays private while it stabilises (ADR-0030), callers assemble built-in algorithms manually (ADR-0032), and validation and state-type handling move to where they belong (ADR-0031, ADR-0033, ADR-0034). The C04 checkpoint marks this boundary in place with token bucket not yet added.

---

## ADR-0026 — Introduce the algorithm interface
**Status:** Accepted

**Context:** V3 isolated fixed-window logic in a helper. Token bucket — different config, different state — is next.

**Decision:** Introduce an algorithm boundary now: the limiter orchestrates validation, locking, lookup, and writeback; algorithms own the decision rules.

**Consequences:** Real variation justifies the abstraction (resolving ADR-0013's wait). Adds indirection now in exchange for a clean limiter/algorithm split and easy future algorithms.

## ADR-0027 — Algorithm-owned state
**Status:** Accepted

**Context:** Fixed window needs (count, window-start); token bucket needs (tokens, last-refill). 

**Decision:** Each algorithm owns its own state shape rather than sharing one struct with every possible field.

**Consequences:** Avoids a growing "god state". The limiter must store state generically. Rejected alternatives: one concrete struct for all (weak as algorithms grow) and raw `any` (loses type safety).

## ADR-0028 — Marker interface for algorithm state
**Status:** Accepted

**Context:** The limiter must store per-key state without knowing whether it's fixed-window or token-bucket state. `any` is too loose; a shared struct accumulates irrelevant fields.

**Decision:** Use a marker interface (`algorithmState` with one unexported no-op method) that each concrete state type implements.

**Consequences:** The package controls which state types are valid; the limiter stores them opaquely. Adds ceremony and still needs a runtime type assertion inside each algorithm.

## ADR-0029 — Explicit `exists` flag for missing state
**Status:** Accepted

**Context:** Algorithms must tell a new key from an existing one with stored state.

**Decision:** Keep an explicit `exists` boolean in the decision method.

**Consequences:** Matches the existing helper and avoids subtle nil-interface behavior. Signature is slightly awkward; revisit if pressure builds for an `InitialState(now)` method instead.

## ADR-0030 — Keep the algorithm interface private
**Status:** Accepted

**Context:** The abstraction is new and will move as token bucket lands.

**Decision:** Keep the interface unexported for now.

**Consequences:** The package can refactor the contract while learning from two concrete algorithms without breaking external callers. External custom algorithms aren't possible yet — deferred until the contract stabilises (echoed for Redis in ADR-0067).

## ADR-0031 — Algorithm owns config validation
**Status:** Accepted

**Context:** Fixed window and token bucket have different config invariants.

**Decision:** Each algorithm constructor validates its own config (`NewFixedWindow` checks limit/window; `NewTokenBucket` checks capacity/rate); `NewLimiter` validates only limiter-level concerns. (This is where ADR-0011's constructor panic becomes a returned error.)

**Consequences:** The owner of an invariant validates it, at construction. Callers do two construction steps, but generic limiter code stays free of algorithm-specific settings.

## ADR-0032 — Manual assembly with built-in algorithms
**Status:** Accepted

**Context:** The algorithm interface is private but callers must still choose and configure an algorithm.

**Decision:** Callers assemble a configured package-provided algorithm with a limiter (`configured algorithm → limiter → Allow`).

**Consequences:** Exposes composition without committing to public custom-algorithm support. The package keeps freedom to refactor the private contract.

## ADR-0033 — `Decide` returns state errors
**Status:** Accepted

**Context:** Algorithms receive opaque state and must assert their concrete type.

**Decision:** Decision methods return an error alongside result and new state, so a wrong-algorithm state mismatch is reported deliberately rather than panicking.

**Consequences:** Larger signature; `Allow` must propagate algorithm errors before storing state.

## ADR-0034 — Initialize missing state before asserting
**Status:** Accepted

**Context:** A new key has no stored state; an existing key should carry the algorithm's concrete type.

**Decision:** Handle missing state (normal initialization) before the type assertion (which now means real wrong-type error).

**Consequences:** Each algorithm has an explicit init branch before its assertion, cleanly separating "cold start" from "corrupt state".

---

## Checkpoint C04 — Algorithm boundary in place
Fixed-window behavior now sits behind the private `algorithm` interface with a marker `algorithmState`, fixed-window config/state, generic limiter orchestration, an explicit state-mismatch error, and nil-algorithm validation. Token bucket not yet added. Known follow-up at the time: two constructor-validation tests masked `NewFixedWindow`'s error by then calling `NewLimiter` with a nil algorithm — they should assert the algorithm-constructor error directly. Next: fix those tests, add direct nil-algorithm coverage, then design token-bucket state and refill math (V5).

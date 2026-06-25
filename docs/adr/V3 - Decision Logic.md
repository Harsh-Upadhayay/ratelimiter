# V3 — Decision Logic (ADR-0023 … ADR-0025)

> Back to [[Rate Limiter Learning Map]] · Prev [[V2 - API Evolution]] · Next [[V4 - Algorithm Abstraction]] · Index [[README|ADR index]]

V3 is a local refactor, not a feature. `Allow` had grown to mix validation, locking, map access, the fixed-window math, and result building. V3 extracts the pure decision math into a private deterministic helper (ADR-0023) and leaves `Allow` as an orchestrator (ADR-0025) — a first, interface-free sketch of the storage/algorithm boundary that V4 and V6 formalise. Tests stay on the public API so the refactor is free (ADR-0024).

---

## ADR-0023 — Private fixed-window decision helper
**Status:** Accepted

**Context:** `Allow` mixed input validation, locking, map reads, fixed-window decisions, mutation, and result building.

**Decision:** Extract a private, deterministic helper that takes `(now, limit, window, state, exists)` and returns `(result, newState)` — grouped time → config → state → metadata. Rejected requests return unchanged state (rejections don't consume quota or move the window).

**Consequences:** Algorithm behavior separates from map access and locking without introducing public interfaces prematurely. The `exists` parameter is slightly awkward but keeps map lookup out of the algorithm.

## ADR-0024 — Test through the public API
**Status:** Accepted

**Context:** V3 adds a private helper but external behavior is unchanged.

**Decision:** Keep tests on the public `Limiter` API.

**Consequences:** Preserves refactoring freedom — the helper can change shape without rewriting tests. Some internal edge cases are slightly harder to target directly, but all important outcomes are observable through `Allow`.

## ADR-0025 — `Allow` as orchestrator
**Status:** Accepted

**Context:** With decision math in a helper, something must still own validation, locking, map access, and persistence.

**Decision:** `Allow` orchestrates: validate input → lock → read state → call the helper → write state back → return result/error.

**Consequences:** A local, interface-free version of the storage/algorithm boundary. `Allow` still knows both storage and algorithm details — acceptable until a second store or algorithm creates real pressure (V4, V6).

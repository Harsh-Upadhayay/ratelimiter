# V2 — API Evolution (ADR-0019 … ADR-0022)

> Back to [[Rate Limiter Learning Map]] · Prev [[V1 - Fixed Window]] · Next [[V3 - Decision Logic]] · Index [[README|ADR index]]

V2 grows the API from a bare boolean into something an HTTP caller can actually use. `Allow` now returns a `Result` struct plus an `error` (ADR-0019, ADR-0021): rate-limit rejection is a normal decision (`Allowed: false`, `nil` error), while invalid input is a real error carried by a sentinel (ADR-0022). The user-facing vocabulary also generalises from "user ID" to "rate-limit key" (ADR-0020), since real limiters key on IPs, tokens, tenants, or composites.

---

## ADR-0019 — Result-and-error return
**Status:** Accepted

**Context:** V1's boolean can't carry retry timing or quota detail, and can't distinguish a rejection from a failure.

**Decision:** `Allow` returns `(Result, error)`. Allowed → populated result, `nil` err. Rejected → rejected result, `nil` err. Invalid key → zero result, non-nil err.

**Consequences:** Rate limiting is modelled as a normal decision, not an error. For the in-memory limiter errors are mostly validation; later a distributed store produces real operational errors (Redis failure, CAS exhaustion). Callers and tests now handle two return values.

## ADR-0020 — Generic rate-limit key
**Status:** Accepted

**Context:** V1 spoke of `userID`, but real limiters key on user, IP, API token, tenant, route, or a composite.

**Decision:** Rename the public input concept from "user ID" to "rate-limit key".

**Consequences:** The API stops being tied to users; HTTP middleware can build keys from request attributes (ADR-0076). The abstraction needs clear examples so callers know what to key on.

## ADR-0021 — `Result` contract
**Status:** Accepted

**Context:** Callers need detail for HTTP headers and retry behavior.

**Decision:** Export `Result{Allowed, Remaining, RetryAfter}`. `Remaining` = requests left this window after the decision; `RetryAfter` = zero when allowed, time-until-reset when rejected. Fixed window computes `RetryAfter = windowStart + window − now`.

**Consequences:** A genuinely useful response, at the cost that every branch must maintain the result-field invariants. This contract stays stable all the way to the HTTP boundary (ADR-0077, ADR-0082).

## ADR-0022 — Sentinel error for empty key
**Status:** Accepted

**Context:** V1 panicked on an empty key; V2 has an error channel, so invalid input needn't crash the caller.

**Decision:** Define an exported sentinel `ErrEmptyKey`, comparable with `errors.Is`. Use generic key language, not `ErrEmptyUserID`.

**Consequences:** Empty key becomes a stable, branchable public error category. Exported errors are now part of the public API and must stay stable.

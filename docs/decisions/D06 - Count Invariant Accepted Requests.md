# D06 - Count Invariant Accepted Requests

Back to [[Rate Limiter Learning Map]].

## Context

The count field needs a stable meaning between requests.

## Decision

Stored count means accepted requests in the current window.

## Why

State should represent a fact at rest, not a temporary "before the current request" view.

## Tradeoff

If a higher-level policy should reject without consuming quota, it must run before this limiter.

## Revisit when

When composing multiple limiters or introducing global and per-user limits.

## Links

- [[D02 - State Shape Map String UserState]]
- [[D10 - Return Type]]

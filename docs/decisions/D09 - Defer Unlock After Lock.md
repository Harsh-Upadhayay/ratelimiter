# D09 - Defer Unlock After Lock

Back to [[Rate Limiter Learning Map]].

## Context

The allow decision has several branches that can return early.

## Decision

For V1, call deferred unlock immediately after locking.

## Why

Conceptually, this guarantees every return branch releases the lock.

## Tradeoff

`defer` has a small overhead and keeps the lock held until the function returns. For V1, clarity is worth it.

## Revisit when

Only after correctness is established and there is evidence that this path is hot enough to justify manual unlocks.

## Links

- [[D04 - One Global Mutex]]
- [[G05 - Defer Unlock Pattern]]

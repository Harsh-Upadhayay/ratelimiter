# G05 - Defer Unlock Pattern

Back to [[Rate Limiter Learning Map]].

## Go practice

When a function has multiple return branches, defer the unlock immediately after acquiring the lock.

## Why it matters

This makes lock release function-scoped and avoids missing an unlock in one branch.

## Caveat

`defer` has a tiny overhead and keeps the lock held until the function returns. For V1, this is a good clarity tradeoff. In measured hot paths, manual unlocks can be considered later.

## Links

- [[D09 - Defer Unlock After Lock]]
- [[G04 - Mutexes and Critical Sections]]

# D61 - Test Isolation with Fresh Stores

Back to [[V7 Sharded MemoryStore Index]].

## Decision

Every test gets its own freshly constructed store. No package-level `var store = NewMemoryStore()` shared across tests.

## Why

A shared store accumulates state across tests. Test A's `CompareAndSwap` leaves a key with version 1; Test B's setup assumes version 0 and gets a CAS conflict. Tests pass or fail depending on execution order — a classic test pollution problem.

## How

Algorithm tests use `newTestLimiter(t, algo)` which constructs a fresh store on each call. Store contract tests construct a new store at the top of each `Test...` function.

## Links

- [[D60 - Contract Testing for StateStore]]
- [[G36 - Test Helper Pattern]]

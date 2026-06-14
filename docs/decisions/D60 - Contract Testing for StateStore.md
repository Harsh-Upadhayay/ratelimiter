# D60 - Contract Testing for StateStore

Back to [[V7 Sharded MemoryStore Index]].

## Decision

Test the `StateStore` interface contract through a shared helper `testXxx(t *testing.T, store StateStore)` rather than writing separate test suites per implementation.

Two top-level `Test...` functions instantiate each backend and pass it to the same helpers:

```go
func TestMemoryStoreContract(t *testing.T)        { testXxx(t, NewMemoryStore()) }
func TestShardedMemoryStoreContract(t *testing.T) { testXxx(t, store) }
```

## Why

Any type that satisfies `StateStore` must behave identically from the caller's perspective. A shared contract test suite proves this without duplication. Adding a third backend (e.g. `RedisStore`) means passing it to the same helpers — no new test logic needed.

## What the contract covers

1. `Get` on missing key → `nil, 0, false, nil`
2. `Get` on existing key → correct state and version
3. CAS on missing key with version 0 → stores, returns `true`, version becomes 1
4. CAS on missing key with wrong version → returns `false`
5. CAS on existing key with correct version → updates, returns `true`, version increments
6. CAS on existing key with stale version → returns `false`
7. Concurrent access doesn't corrupt state

## Links

- [[D54 - Sharded Store Keeps StateStore Contract]]
- [[G37 - Contract Testing via Interface Parameter]]

# G36 - Test Helper Pattern

Back to [[Go Concepts Index]].

## Concept

A test helper is a function that does setup or assertions shared across multiple tests. In Go, helpers that call `t.Fatal` must accept `*testing.T` and call `t.Helper()` so failure line numbers point at the caller, not the helper.

## Shape

```go
func newTestMemoryLimiter(t *testing.T, algo memoryAlgorithm) *MemoryLimiter {
    t.Helper()
    store := NewMemoryStore()
    lim, err := NewMemoryLimiter(algo, store)
    if err != nil {
        t.Fatalf("failed to create limiter: %v", err)
    }
    return lim
}
```

- Unexported (`newTestMemoryLimiter`, not `NewTestMemoryLimiter`) — test helpers are not public API.
- Returns the value directly, not `(value, error)` — the helper owns error handling via `t.Fatal`.
- `t.Helper()` marks this frame so stack traces skip it and point to the test that called it.

## Why not a package-level var

A shared `var store = NewMemoryStore()` accumulates state across tests, causing test pollution. A helper function creates a fresh instance each call. See [[D61 - Test Isolation with Fresh Stores]].

## Related

- [[G15 - Public API Tests]]
- [[D61 - Test Isolation with Fresh Stores]]

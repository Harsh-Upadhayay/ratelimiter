# G37 - Contract Testing via Interface Parameter

Back to [[Go Concepts Index]].

## Concept

When multiple types implement the same interface, write one shared test helper that accepts the interface, then call it once per implementation. This is called contract testing.

## Shape

```go
func testStateStoreContract(t *testing.T, store StateStore) {
    t.Helper()
    // all behavioral assertions here
}

func TestMemoryStoreContract(t *testing.T) {
    testStateStoreContract(t, NewMemoryStore())
}

func TestShardedMemoryStoreContract(t *testing.T) {
    store, err := NewShardedMemoryStore(8)
    if err != nil { t.Fatal(err) }
    testStateStoreContract(t, store)
}
```

## Why

- No duplicated test logic across implementations.
- Adding a new backend (e.g. `RedisStore`) means one new `Test...` function — not rewriting all assertions.
- The helper proves the interface contract holds, not just that one implementation works.

## Key rule

The helper must only use the interface — no type assertions to concrete types, no access to internal fields. If it needs internal access, it's testing implementation, not contract.

## Related

- [[G17 - Interfaces From Real Variation]]
- [[G18 - Structural Interface Satisfaction]]
- [[D60 - Contract Testing for StateStore]]

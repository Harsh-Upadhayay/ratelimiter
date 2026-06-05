# D47 - StateStore Uses Get and CAS

Back to [[Rate Limiter Learning Map]].

## Context

The limiter needs to read state, run an algorithm decision, and commit updated state without losing concurrent updates.

## Decision

Use a store boundary based on `Get` plus compare-and-swap.

## Concept

```text
state, version = Get(key)
result, newState = Decide(state)
ok = CAS(key, version, newState)
```

If CAS fails, another request changed the state after the read. The limiter must retry.

## Why

This preserves read-decide-write correctness without relying on one process-local limiter mutex.

## Tradeoff

CAS adds versions, retry logic, and conflict handling. It is more complex than a simple in-memory mutex, but it matches the distributed-system direction.

## Links

- [[D46 - GetSet Race]]
- [[G26 - Optimistic Concurrency with CAS]]
- [[plan|Target 6-step plan]]

# G29 - Race Detector

Back to [[Rate Limiter Learning Map]].

## Concept

`go test -race` runs tests with Go's race detector enabled.

It looks for unsynchronized concurrent access to shared memory.

## Rate limiter use

The rate limiter has shared mutable state, first in `Limiter` and now in `MemoryStore`.

Running the race detector checks whether concurrent tests are hitting unsafe shared-memory access.

## Limits

The race detector finds data races, not all concurrency bugs.

It can pass even when the algorithm has a logical race, such as stale read-decide-write behavior. CAS exists to protect that logical invariant.

## Links

- [[D48 - Benchmark Before Storage Refactor]]
- [[G04 - Mutexes and Critical Sections]]
- [[G26 - Optimistic Concurrency with CAS]]

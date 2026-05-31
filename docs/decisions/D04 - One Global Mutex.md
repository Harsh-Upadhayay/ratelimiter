# D04 - One Global Mutex

Back to [[Rate Limiter Learning Map]].

## Context

Concurrent requests can read the same count, both decide to allow, and overwrite each other.

## Decision

Use one global mutex for the V1 limiter.

## Why

It is the simplest correct protection for the full read, decide, update sequence.

## Tradeoff

Requests for different users serialize, so throughput is lower than with per-key locking.

## Alternatives

- Per-key locks: better concurrency, more memory and lifecycle complexity.
- Sharded locks: middle ground, more complex than one lock but cheaper than one lock per key.

## Revisit when

When load testing or when learning lock granularity and hot-key behavior.

## Links

- [[G04 - Mutexes and Critical Sections]]
- [[D09 - Defer Unlock After Lock]]

# D49 - MemoryStore Owns Runtime State

Back to [[Rate Limiter Learning Map]].

## Context

`Limiter` previously owned the per-key state map and mutex.

V6 introduces a `StateStore` boundary.

## Decision

Move runtime state ownership into `MemoryStore`.

`MemoryStore` owns:

- the map from key to stored record
- the version for each key
- the mutex protecting the map

`Limiter` owns orchestration only.

## Why

This separates storage concerns from rate-limit decision flow.

The limiter can ask for state, run the algorithm, and attempt a commit without knowing how the state is stored.

## Tradeoff

The current memory store still uses one mutex internally, so this does not yet optimize many-key concurrency. It creates the boundary needed for future store implementations.

## Links

- [[D45 - Split Storage from Limiter]]
- [[D47 - StateStore Uses Get and CAS]]
- [[G27 - Store Owned Mutexes]]
- [[G26 - Optimistic Concurrency with CAS]]

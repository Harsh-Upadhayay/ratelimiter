# D54 - Sharded Store Keeps StateStore Contract

Back to [[Rate Limiter Learning Map]].

## Context

`StateStore` already exposes `Get` and `CompareAndSwap`.

## Decision

Implement sharding behind the existing `StateStore` interface.

The interface does not change.

## Why

The limiter should not care whether the backend has one mutex, many shard mutexes, Redis, or another storage strategy.

This validates the storage boundary: backend concurrency can improve without changing algorithms or limiter orchestration.

## Tradeoff

All sharding details become backend internals. This keeps the public boundary clean, but tests or benchmarks need backend-specific access if they want to inspect shard distribution.

## Links

- [[D47 - StateStore Uses Get and CAS]]
- [[D53 - Sharded MemoryStore Next]]
- [[G31 - Lock Striping]]

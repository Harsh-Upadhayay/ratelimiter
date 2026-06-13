# D53 - Sharded MemoryStore Next

Back to [[Rate Limiter Learning Map]].

## Context

The next algorithm would mostly repeat the existing algorithm pattern: config, state, constructor validation, and `Decide`.

The current memory backend still serializes all keys behind one mutex.

## Decision

Implement a sharded memory store before adding another algorithm.

## Why

This teaches a new concurrency and systems concept: partitioning ownership by key.

The many-key benchmark creates concrete pressure for this change because unrelated keys currently block each other.

## Tradeoff

Sharding does not help one hot key. It improves concurrency when traffic is spread across many keys.

It also adds hashing, shard-count tuning, and collision behavior.

## Links

- [[D48 - Benchmark Before Storage Refactor]]
- [[D54 - Sharded Store Keeps StateStore Contract]]
- [[G31 - Lock Striping]]
- [[G32 - Key Hashing]]

# D56 - Reuse MemoryStore Internally

Back to [[Rate Limiter Learning Map]].

## Context

A sharded store can either define a dedicated shard type or compose existing memory stores.

## Decision

Use internal `*MemoryStore` shards.

```text
ShardedMemoryStore
  shards []*MemoryStore
```

## Why

Each shard already needs the same behavior as `MemoryStore`: a map, a mutex, `Get`, and `CompareAndSwap`.

Reusing `MemoryStore` avoids duplicating CAS logic.

## Tradeoff

This is composition over a public concrete type. It is slightly less semantically precise than a private `memoryShard`, but it keeps the implementation smaller for this learning step.

Use pointers so the store does not copy mutex-containing `MemoryStore` values.

## Links

- [[D49 - MemoryStore Owns Runtime State]]
- [[D54 - Sharded Store Keeps StateStore Contract]]
- [[G33 - Composition with Pointer Fields]]
- [[G07 - Do Not Copy Mutexes]]

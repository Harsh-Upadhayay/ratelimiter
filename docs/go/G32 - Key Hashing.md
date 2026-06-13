# G32 - Key Hashing

Back to [[Rate Limiter Learning Map]].

## Concept

Hashing converts a key into a stable numeric value.

The store can map that value to a shard:

```text
hash(key) % shardCount
```

## Rate limiter use

The same key must always map to the same shard, otherwise `Get` and `CompareAndSwap` could operate on different records.

## Tradeoff

Hashing adds a small per-request cost, but it enables partitioned ownership and better many-key concurrency.

Hash distribution affects contention. Poor distribution can overload one shard.

## Links

- [[D53 - Sharded MemoryStore Next]]
- [[G31 - Lock Striping]]

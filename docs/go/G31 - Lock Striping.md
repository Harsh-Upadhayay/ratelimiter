# G31 - Lock Striping

Back to [[Rate Limiter Learning Map]].

## Concept

Lock striping splits one large critical section into multiple smaller lock-protected partitions.

Instead of:

```text
all keys -> one mutex
```

use:

```text
key -> shard -> shard mutex
```

## Rate limiter use

`ShardedMemoryStore` maps each key to one shard. Each shard has its own `MemoryStore` and therefore its own mutex.

Same-key requests still serialize. Different keys can proceed concurrently when they land on different shards.

## Tradeoff

More shards can reduce contention, but they add memory overhead and hash work.

Two unrelated hot keys can still block each other if they hash to the same shard.

## Links

- [[D53 - Sharded MemoryStore Next]]
- [[D54 - Sharded Store Keeps StateStore Contract]]
- [[G04 - Mutexes and Critical Sections]]

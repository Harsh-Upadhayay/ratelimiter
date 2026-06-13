# D55 - Configurable Shard Count

Back to [[Rate Limiter Learning Map]].

## Context

A sharded store needs a fixed number of shards.

## Decision

Use a configurable shard count:

```text
NewShardedMemoryStore(shardCount int)
```

If `shardCount <= 0`, return an error.

## Why

This keeps the implementation useful for learning and benchmarking. Different shard counts can show the tradeoff between memory overhead and concurrency.

## Tradeoff

The constructor needs validation and error handling. A hardcoded default would be simpler, but less useful for experiments.

## Links

- [[D53 - Sharded MemoryStore Next]]
- [[G21 - Constructor Validation Ownership]]
- [[G31 - Lock Striping]]

# D57 - Sharded MemoryStore as Default Backend

Back to [[Rate Limiter Learning Map]].

## Context

Store injection is deferred. `NewLimiter(algo)` currently creates the default store internally.

## Decision

After `ShardedMemoryStore` exists, make it the default backend used by `NewLimiter`.

## Why

This lets the new backend affect existing tests and benchmarks without introducing constructor injection yet.

It keeps the external API simple while improving the default local concurrency behavior.

## Tradeoff

Because `NewShardedMemoryStore` returns an error for invalid shard counts, the default constructor needs a safe internal default shard count.

This still does not expose backend choice to callers. Store injection remains deferred.

## Links

- [[D52 - Default Memory Store Constructor]]
- [[D55 - Configurable Shard Count]]
- [[D53 - Sharded MemoryStore Next]]

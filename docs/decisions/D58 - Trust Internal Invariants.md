---
Back to [[V7 Sharded MemoryStore Index]].

## Decision

`Get` and `CompareAndSwap` on `ShardedMemoryStore` do not bounds-check the index returned by `shardIndex`. They trust it directly.

## Reasoning

`shardIndex` is a private method on the same struct. Its correctness is a local, structural invariant: it always returns `h.Sum32() % uint32(len(s.shards))`, which is always in `[0, len-1]`. Adding a bounds check in `Get`/`CAS` would be defensive coding against your own internal code.

Defensive checks at internal call sites add noise without catching real bugs — if `shardIndex` were broken, a bounds check would just change the failure mode from a panic to a silent error. Neither outcome is better than fixing the root cause.

**Rule:** validate at system boundaries (user input, external APIs, constructor parameters). Trust internal contracts between methods on the same type.

## Contrast

This is different from `NewShardedMemoryStore` validating `shardCount <= 0` — that is a system boundary (caller-supplied input). Internal method calls are not.

## Related

- [[D55 - Configurable Shard Count]]
- [[G32 - Key Hashing]]

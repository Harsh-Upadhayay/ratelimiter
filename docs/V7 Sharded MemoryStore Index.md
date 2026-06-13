# V7 Sharded MemoryStore Index

Back to [[Rate Limiter Learning Map]].

This hub tracks the move from one memory-store mutex to a sharded in-memory backend.

## Core decisions

- [[D53 - Sharded MemoryStore Next]]
- [[D54 - Sharded Store Keeps StateStore Contract]]
- [[D55 - Configurable Shard Count]]
- [[D56 - Reuse MemoryStore Internally]]
- [[D57 - Sharded MemoryStore as Default Backend]]

## Prior pressure

- [[D48 - Benchmark Before Storage Refactor]]
- [[D49 - MemoryStore Owns Runtime State]]
- [[D52 - Default Memory Store Constructor]]

## Go and systems concepts used

- [[G31 - Lock Striping]]
- [[G32 - Key Hashing]]
- [[G33 - Composition with Pointer Fields]]

## Future bridge

- [[plan|Target 6-step plan]]

# G33 - Composition with Pointer Fields

Back to [[Rate Limiter Learning Map]].

## Concept

Go structs can compose behavior by storing other concrete types as fields.

When the composed type contains a mutex, prefer pointers to avoid copying mutex-containing values.

## Rate limiter use

`ShardedMemoryStore` can hold:

```text
[]*MemoryStore
```

Each shard delegates `Get` and `CompareAndSwap` to an existing `MemoryStore`.

## Why pointers

Copying a `MemoryStore` would copy its mutex field while still potentially sharing map state. That can break synchronization.

Pointers keep each shard as one shared object.

## Links

- [[D56 - Reuse MemoryStore Internally]]
- [[G07 - Do Not Copy Mutexes]]
- [[G27 - Store Owned Mutexes]]

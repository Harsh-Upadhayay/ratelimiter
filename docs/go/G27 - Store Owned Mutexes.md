# G27 - Store Owned Mutexes

Back to [[Rate Limiter Learning Map]].

## Concept

A mutex must be stored on the shared object it protects.

Creating a new local mutex inside each method call does not protect shared state because every call locks a different mutex.

## Rate limiter use

`MemoryStore` owns the state map, so `MemoryStore` owns the mutex.

Both reads and writes of a Go map must happen while holding the mutex when concurrent writes are possible.

## Pointer receivers

A type containing a mutex should use pointer receivers for mutating methods.

Copying a struct that contains a mutex can create multiple locks protecting the same underlying map, which defeats synchronization.

## Links

- [[D49 - MemoryStore Owns Runtime State]]
- [[G04 - Mutexes and Critical Sections]]
- [[G07 - Do Not Copy Mutexes]]

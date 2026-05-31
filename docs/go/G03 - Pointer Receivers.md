# G03 - Pointer Receivers

Back to [[Rate Limiter Learning Map]].

## Go practice

Use pointer receivers when a method mutates receiver state or when copying the receiver would be misleading or expensive.

## Why it matters

The limiter owns mutable state. Its allow method changes that state, so the method should operate on the limiter itself.

## Caveat

Pointer receivers make shared mutation possible. Shared mutation needs a clear concurrency story.

## Links

- [[D03 - Pointer Receiver for Mutating Limiter]]
- [[G04 - Mutexes and Critical Sections]]

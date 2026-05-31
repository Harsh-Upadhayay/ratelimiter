# D03 - Pointer Receiver for Mutating Limiter

Back to [[Rate Limiter Learning Map]].

## Context

The limiter owns mutable internal state: configuration, a user-state map, and a mutex.

## Decision

Methods that mutate limiter state should use pointer receivers.

## Why

The method is conceptually operating on the limiter itself, not a copy of it.

## Tradeoff

Pointer receivers make shared mutable state explicit, which means locking rules must be clear.

## Revisit when

When designing pure algorithm objects later. Pure algorithm methods may not need mutation at all.

## Links

- [[G03 - Pointer Receivers]]
- [[D04 - One Global Mutex]]

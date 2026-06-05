# D45 - Split Storage from Limiter

Back to [[Rate Limiter Learning Map]].

## Context

`Limiter` currently owns the per-key state map and mutex. That means orchestration, storage, and concurrency control are still coupled.

## Decision

Move storage behind a store boundary.

## Why

The limiter should orchestrate validation, state loading, algorithm decision, and state commit without knowing whether state lives in a map, Redis, or another backend.

## Tradeoff

This introduces another interface and forces the design to handle read-modify-write correctness explicitly.

## Links

- [[D46 - GetSet Race]]
- [[D47 - StateStore Uses Get and CAS]]
- [[G26 - Optimistic Concurrency with CAS]]

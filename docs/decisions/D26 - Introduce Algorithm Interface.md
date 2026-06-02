# D26 - Introduce Algorithm Interface

Back to [[Rate Limiter Learning Map]].

## Context

The project started with one fixed-window algorithm. V3 extracted fixed-window decision logic into a private helper. The next algorithm, token bucket, has different configuration and state.

## Decision

Introduce an algorithm boundary now.

## Why

There is now real variation. Fixed window and token bucket both answer the same decision question, but they use different state and configuration.

## Tradeoff

This adds design and syntax complexity now, but it should reduce duplication and make later algorithms easier to plug in.

## Boundary intent

The limiter should orchestrate validation, locking, key lookup, and state writeback. Algorithms should own the decision rules.

## Links

- [[D13 - Delay Interfaces]]
- [[D27 - Algorithm Owned State]]
- [[G17 - Interfaces From Real Variation]]

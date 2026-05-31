# D15 - One File First

Back to [[Rate Limiter Learning Map]].

## Context

The V1 implementation is intentionally small: one limiter, one algorithm, one storage mechanism.

## Decision

Start with one implementation file.

## Why

This makes real friction visible before introducing file splits or abstractions.

## Tradeoff

The file may become crowded as features are added. That crowding is useful signal for when to split responsibilities.

## Revisit when

When adding tests, richer result types, additional algorithms, or storage boundaries makes the file hard to navigate.

## Links

- [[D13 - Delay Interfaces]]
- [[D14 - Package Name]]

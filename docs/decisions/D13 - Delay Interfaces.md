# D13 - Delay Interfaces

Back to [[Rate Limiter Learning Map]].

## Context

The target design includes `Algorithm` and `StateStore` boundaries. V1 has only one algorithm and one storage mechanism.

## Decision

Do not introduce `Algorithm` or `StateStore` interfaces in V1.

## Why

Interfaces should describe behavior that is known to vary. With only one implementation, the interface would mostly be a guess.

## Tradeoff

There may be more refactoring later when the second algorithm or store appears.

## Future direction

Introduce interfaces after concrete implementations reveal the true shared contract.

## Links

- [[D12 - Exported API Boundary]]
- [[plan|Target 6-step plan]]

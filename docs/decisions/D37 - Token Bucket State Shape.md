# D37 - Token Bucket State Shape

Back to [[Rate Limiter Learning Map]].

## Context

Token bucket needs runtime state that differs from fixed-window state.

## Decision

Use token-bucket-specific state:

- available tokens
- last refill time

## Why

Available tokens captures current spendable capacity. Last refill time lets the algorithm lazily compute how many tokens should have accumulated since the previous request.

## Tradeoff

The limiter cannot inspect this state directly. It stores it opaquely through the algorithm-state boundary.

## Links

- [[D27 - Algorithm Owned State]]
- [[D38 - Lazy Token Refill]]
- [[D39 - Floating Point Token Arithmetic]]

# D35 - Token Bucket as Second Algorithm

Back to [[Rate Limiter Learning Map]].

## Context

Fixed window is implemented and now sits behind the internal algorithm boundary. The next step is adding a second algorithm to validate that boundary.

## Decision

Use token bucket as the second algorithm.

## Why

Token bucket has different configuration and state from fixed window while still keeping O(1) memory and O(1) request-time decision cost.

## Tradeoff

Compared with fixed window, token bucket has higher model complexity: refill arithmetic, capacity capping, fractional tokens, retry timing, and clock-regression handling.

## Links

- [[D26 - Introduce Algorithm Interface]]
- [[D36 - Token Bucket Capacity and Refill Rate]]
- [[D37 - Token Bucket State Shape]]

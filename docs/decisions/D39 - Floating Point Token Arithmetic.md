# D39 - Floating Point Token Arithmetic

Back to [[Rate Limiter Learning Map]].

## Context

Token refill can produce fractional tokens.

Example:

```text
3 tokens/second * 1.5 seconds = 4.5 tokens
```

## Decision

Use `float64` internally for available tokens and refill arithmetic.

## Why

This keeps the learning implementation simple and preserves fractional refill progress.

## Tradeoff

Floating-point arithmetic can have precision edge cases. If exact deterministic accounting becomes important, revisit fixed-point integer tokens.

## Links

- [[D36 - Token Bucket Capacity and Refill Rate]]
- [[D40 - Whole Request Remaining]]
- [[G23 - Duration Seconds Conversion]]

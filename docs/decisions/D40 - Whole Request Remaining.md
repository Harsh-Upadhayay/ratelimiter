# D40 - Whole Request Remaining

Back to [[Rate Limiter Learning Map]].

## Context

The public `Result.Remaining` field is an integer, but token bucket can have fractional internal tokens.

## Decision

Keep `Result.Remaining` as whole requests available now by truncating internal tokens.

## Why

The public API answers how many complete one-token requests can pass immediately.

## Tradeoff

Callers do not see exact internal bucket state. For example, `5.9` internal tokens becomes `Remaining = 5`.

## Links

- [[D21 - Result Contract]]
- [[D39 - Floating Point Token Arithmetic]]

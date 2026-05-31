# D05 - Explicit Time Input

Back to [[Rate Limiter Learning Map]].

## Context

The limiter needs the current time to decide whether a window has expired.

## Decision

Pass `now` into the allow decision for V1.

## Why

Tests become deterministic. Boundary cases can be tested exactly.

## Tradeoff

The call site is slightly less ergonomic than a method that reads the wall clock internally.

## Revisit when

When introducing an injectable clock abstraction or Redis `TIME` for distributed clock consistency.

## Links

- [[G06 - Time and Duration Boundaries]]
- [[D08 - Half Open Window Interval]]

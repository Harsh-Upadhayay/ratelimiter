# D08 - Half Open Window Interval

Back to [[Rate Limiter Learning Map]].

## Context

The limiter needs exact boundary behavior for window expiry.

## Decision

Use the interval `[WindowStart, WindowStart + window)`.

## Why

The end boundary belongs to the next window, so adjacent windows do not overlap.

## Tradeoff

Exact boundary behavior must be tested explicitly.

## Revisit when

When writing tests for "one instant before expiry", "exactly at expiry", and "after expiry".

## Links

- [[D05 - Explicit Time Input]]
- [[G06 - Time and Duration Boundaries]]

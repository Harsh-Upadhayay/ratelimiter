# D27 - Algorithm Owned State

Back to [[Rate Limiter Learning Map]].

## Context

Fixed window and token bucket need different state shapes.

Fixed window uses:

- consumed request count
- window start time

Token bucket uses:

- token count
- last refill time

## Decision

Prefer algorithm-owned state over one shared state struct with every possible field.

## Why

This avoids a growing "god state" where many fields only apply to some algorithms.

## Tradeoff

The design requires more thought now. The limiter must store state in a way that supports multiple algorithm-specific representations.

## Alternatives rejected for now

- One concrete state struct for all algorithms: simple but weak as algorithms grow.
- `any` state: flexible but loses type safety and pushes errors to runtime.

## Links

- [[D26 - Introduce Algorithm Interface]]
- [[G17 - Interfaces From Real Variation]]

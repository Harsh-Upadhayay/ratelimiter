# D07 - Window Reset Starts At Now

Back to [[Rate Limiter Learning Map]].

## Context

When a request arrives after the previous window expired, the limiter needs to choose the next window start.

## Decision

For V1, reset the window start to `now`.

## Why

This handles idle gaps simply and is easy to reason about per user.

## Tradeoff

Windows are not globally aligned. This is closer to a per-key fixed window than a wall-clock fixed window.

## Alternatives

- `oldStart + window`: preserves cadence but must handle skipped windows after idle gaps.
- Rounded boundary: aligns all users to clock boundaries but can amplify boundary bursts.

## Revisit when

During low-level design discussion about correctness, burst behavior, memory cleanup, and operational predictability.

## Links

- [[D01 - Fixed Window First]]
- [[D08 - Half Open Window Interval]]
- [[G06 - Time and Duration Boundaries]]

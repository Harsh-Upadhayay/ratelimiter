# D01 - Fixed Window First

Back to [[Rate Limiter Learning Map]].

## Context

The target design includes token bucket, sliding window counter, and sliding window log. For V1, we need a simpler implementation to learn Go state, maps, time, and locking.

## Decision

Start with a fixed window limiter.

## Why

Fixed window is the smallest useful rate limiter: one count and one window start per key.

## Tradeoff

It allows boundary bursts. A user can spend quota near the end of one window and again at the start of the next.

## Revisit when

After the V1 limiter is correct and tested, compare it with sliding window and token bucket.

## Links

- [[D07 - Window Reset Starts At Now]]
- [[D08 - Half Open Window Interval]]
- [[G06 - Time and Duration Boundaries]]

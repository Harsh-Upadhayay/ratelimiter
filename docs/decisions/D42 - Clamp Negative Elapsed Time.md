# D42 - Clamp Negative Elapsed Time

Back to [[Rate Limiter Learning Map]].

## Context

Token bucket computes refill from elapsed time since the last refill. In distributed systems, clock skew can make a later request appear to happen earlier than the stored last refill time.

## Decision

Clamp negative elapsed time to zero.

## Why

This avoids subtracting tokens or producing surprising refill behavior when time moves backward.

## Timestamp rule

If `now` is earlier than `lastRefillTime`, keep `lastRefillTime` unchanged.

## Tradeoff

This is resilient to clock regression, but it hides the event unless metrics are added later.

## Links

- [[D38 - Lazy Token Refill]]
- [[G06 - Time and Duration Boundaries]]

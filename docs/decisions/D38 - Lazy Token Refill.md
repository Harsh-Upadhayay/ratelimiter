# D38 - Lazy Token Refill

Back to [[Rate Limiter Learning Map]].

## Context

A token bucket can be imagined as continuously filling over time.

## Decision

Refill lazily when a request for a key arrives.

## Why

The algorithm can compute current tokens from stored tokens, last refill time, current time, and refill rate. No background worker is needed.

## Tradeoff

Token count is only materialized on access. This is efficient, but inactive keys are not physically updated until they receive traffic again.

## System-design impact

Lazy refill makes CPU cost scale with request traffic instead of with the number of known keys.

## Links

- [[D37 - Token Bucket State Shape]]
- [[G24 - Lazy State Materialization]]

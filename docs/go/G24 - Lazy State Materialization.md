# G24 - Lazy State Materialization

Back to [[Rate Limiter Learning Map]].

## Concept

Lazy state materialization means storing enough information to compute current state when needed instead of continuously updating it in the background.

## Rate limiter use

Token bucket does not continuously fill. On each request, it computes how many tokens should exist now from previous tokens, last refill time, and refill rate.

## Why it matters

This avoids background goroutines, per-key timers, and CPU work for inactive keys.

## Tradeoff

The stored token count may be stale between requests, but the computed decision remains correct when the key is accessed.

## Links

- [[D38 - Lazy Token Refill]]
- [[D37 - Token Bucket State Shape]]

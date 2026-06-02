# D43 - Token Bucket Starts Full

Back to [[Rate Limiter Learning Map]].

## Context

A new key has no existing token-bucket state.

## Decision

Initialize a new key with a full bucket, then consume one token for the first accepted request.

## Why

Capacity represents permitted burst size. Starting empty would reject or delay new keys until tokens accumulate.

## Tradeoff

New keys can burst immediately up to capacity. This is expected token-bucket behavior, but it may be too permissive for abuse-sensitive flows.

## Links

- [[D36 - Token Bucket Capacity and Refill Rate]]
- [[D38 - Lazy Token Refill]]

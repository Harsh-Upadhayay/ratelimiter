# R01 - Key-Value Store and TTL

Back to [[Redis Concepts Index]].

## Concept

Redis is an in-memory key-value store where each key can have an expiration time (TTL — Time To Live). Once a key's TTL expires, Redis automatically deletes it.

## Syntax

```
SET key value EX seconds
GET key
TTL key        # Returns seconds until expiry; -2 if expired/missing; -1 if no expiry set
```

## Why it matters for rate limiting

A fixed-window rate limiter needs the window to reset. Instead of storing window start time and checking it manually, you can:

1. Store the request count as the key's value.
2. Set its TTL to the window duration.
3. When the TTL expires, the key is automatically deleted — the next request starts fresh.

No background cleanup jobs needed; Redis handles it.

## Related

- [[R02 - Atomic Operations with INCR]]
- [[D62 - Redis as Limiter Not Store]]

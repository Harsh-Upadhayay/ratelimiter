# R02 - Atomic Operations with INCR

Back to [[Redis Concepts Index]].

## Concept

`INCR key` atomically increments the integer value stored at `key` by 1. If the key doesn't exist, it's created with value 0 before incrementing, so `INCR` on a missing key returns 1.

Because Redis is single-threaded, `INCR` is inherently atomic — no two clients can interleave their increments.

## Why separate INCR is not enough

For rate limiting, you need atomicity across *multiple* operations:

```
INCR counter       # increment
if counter == 1:
    EXPIRE counter 10   # set window TTL
```

Between `INCR` and `EXPIRE`, another client's request could hit the key. You need all operations to execute atomically as a unit.

## Solution

Lua scripting (`redis.call()` inside a Lua script) — the entire script runs atomically on the Redis server.

## Related

- [[R03 - Lua Scripting for Atomicity]]
- [[R04 - redis.call() and Command Execution]]

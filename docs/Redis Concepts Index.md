# Redis Concepts Index

Back to [[Rate Limiter Learning Map]].

This hub groups Redis and Lua concepts used by the distributed rate limiter.

## Core Redis concepts

- [[R01 - Key-Value Store and TTL]]
- [[R02 - Atomic Operations with INCR]]
- [[R03 - Lua Scripting for Atomicity]]
- [[R06 - SET with NX and EX Flags]]
- [[R07 - Levels of Atomicity in Redis]] — single command vs MULTI/EXEC vs Lua; optimistic vs pessimistic
- [[R08 - Redis Hashes and HMGET]]

## Lua in Redis

- [[R04 - redis.call() and Command Execution]]
- [[R05 - KEYS and ARGV Parameters]]
- [[R09 - Redis TIME and Unit Conversion]]
- [[R10 - Lua Scripts Embedded in Go Strings]]
- [[R11 - Redis Token Bucket Arithmetic]]

## Version bridges

- [[V8 RedisLimiter Design Index]]

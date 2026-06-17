# G39 - Go Naming Qualifier Order

Back to [[Go Concepts Index]].

## Concept

Go names usually read from general category to specific kind with the qualifier first:

```text
RedisLimiter
MemoryLimiter
HTTPClient
JSONEncoder
```

This is easier to scan than noun-first names:

```text
LimiterRedis
LimiterMemory
ClientHTTP
EncoderJSON
```

## Package context matters

Public names are read with the package name:

```go
ratelimiter.RedisLimiter
ratelimiter.MemoryLimiter
ratelimiter.RedisFixedWindow
```

Avoid names that repeat the package concept unnecessarily, but use qualifiers when they distinguish
real variants.

## Export only real API

Capitalization controls visibility. If an interface is not meant for outside packages to implement,
keep it lowercase:

```go
type redisAlgorithm interface { ... }
```

Exporting it as `RedisAlgorithm` would invite external implementations and freeze that method
contract as public API.

## Links

- [[D69 - Backend Qualified Naming]]
- [[G08 - Exported Types with Unexported Fields]]
- [[G30 - Exported Interfaces With Unexported Types]]

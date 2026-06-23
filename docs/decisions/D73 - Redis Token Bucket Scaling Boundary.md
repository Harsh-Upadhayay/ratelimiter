# D73 - Redis Token Bucket Scaling Boundary

Back to [[V8 RedisLimiter Design Index]].

## Context

Scaled integer arithmetic is an implementation detail of Redis Token Bucket. The question is how far
that detail should spread.

## Options

- **Go precomputes all Redis-ready values** - simpler Lua, stronger constructor-time validation, but
  Go knows more about script units.
- **Lua owns all scaling** - Go passes only domain config (`capacity`, `refillRate`), but scale-related
  validation happens later at script execution time unless duplicated in Go.
- **RedisTokenBucket owns the boundary** - callers and `RedisLimiter` see domain concepts, while the
  concrete Redis algorithm prepares the script arguments it needs.

## Decision

Keep scaling knowledge inside `RedisTokenBucket` and its Lua script.

```text
Caller       -> capacity, refillRate
RedisLimiter -> script, args, raw result parsing
RedisTokenBucket -> scaled Redis args
Lua script   -> scaled integer algorithm
```

Use a package-level unexported constant:

```go
const tokenScale = 1000
```

This is intentionally global inside the package but hidden from external callers.

## Why

`RedisTokenBucket` is the algorithm-specific adapter. It is allowed to know how its Lua script wants
arguments shaped.

The important boundary is that scaled units do not leak into:

- public constructors,
- `Result`,
- `RedisLimiter`,
- callers.

## Tradeoffs

- **Memory:** no runtime state difference.
- **Latency:** precomputing scaled args in Go keeps Lua slightly simpler.
- **Concurrency:** no direct difference; Redis Lua still owns the atomic transition.
- **Maintainability:** the conversion is in two places (`args()` and Lua), so names must stay clear:
  `capacityScaled`, `refillUnitsPerSecond`, and `tokenScale`.
- **API clarity:** callers keep using domain units.

## Links

- [[D70 - Redis Token Bucket Scaled Integer State]]
- [[D72 - Redis Token Bucket Result Contract]]
- [[G40 - Unexported Package Constants]]
- [[redis/R10 - Lua Scripts Embedded in Go Strings]]

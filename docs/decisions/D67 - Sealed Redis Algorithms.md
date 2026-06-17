# D67 - Sealed Redis Algorithms

Back to [[V8 RedisLimiter Design Index]].

## Context

`RedisLimiter` uses a private `redisAlgorithm` interface:

```go
type redisAlgorithm interface {
    script() string
    args() []string
}
```

The constructor can accept Redis algorithm values from this package, such as `*RedisFixedWindow`.
Because the interface and its methods are unexported, callers outside the package cannot define
their own Redis algorithms that satisfy it.

## Options

- **Keep `redisAlgorithm` private** — users can use only package-provided Redis algorithms.
- **Export `RedisAlgorithm`** — users can implement their own Lua-backed algorithms.
- **Avoid a generic Redis constructor** — expose per-algorithm constructors such as
  `NewRedisFixedWindowLimiter`.

## Decision

Keep Redis algorithms sealed for now.

External callers should use the package-provided Redis algorithms. Custom Redis algorithm extension
is deferred until there is real pressure for it.

## Why

The Redis algorithm boundary is lower-level and riskier than the in-process algorithm boundary.
A Redis algorithm is not just a Go function; it defines:

- Lua source,
- `KEYS` / `ARGV` conventions,
- result tuple shape,
- Redis atomicity assumptions,
- error behavior from scripts.

Exporting this interface would freeze those details as public API. Keeping it private lets the
package evolve the Lua contract while the Redis design is still young.

## Tradeoffs

- **Memory:** no difference.
- **Latency:** no difference.
- **Concurrency:** safer by default. Only package-owned scripts define atomic Redis behavior.
- **API flexibility:** lower. Users cannot bring custom Redis algorithms yet.
- **API stability:** higher. The package can change the private Lua interface without breaking
  external implementations.
- **Learning value:** better for now. We can focus on one correct Redis path before designing an
  extension system.

## Future pressure

Revisit this if callers need custom Redis-backed algorithms. At that point, compare:

- exported `RedisAlgorithm`,
- per-algorithm limiter constructors,
- a higher-level declarative algorithm config,
- separate packages for advanced Redis algorithms.

Do not export the interface just because it exists internally.

## Links

- [[D64 - RedisLimiter Algorithm Interface]]
- [[D65 - Redis Adapter Returns Raw Result]]
- [[D66 - Redis Client Ownership Boundary]]
- [[redis/R03 - Lua Scripting for Atomicity]]
- [[redis/R05 - KEYS and ARGV Parameters]]

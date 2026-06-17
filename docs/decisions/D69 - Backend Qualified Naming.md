# D69 - Backend Qualified Naming

Back to [[V8 RedisLimiter Design Index]].

## Context

The project now has two execution backends:

- in-process memory execution,
- Redis/Lua execution.

Using only `Limiter` and `algorithm` became ambiguous once `RedisLimiter` and Redis algorithms
were introduced.

## Decision

Use backend-qualified names for backend-specific concepts.

Conceptually:

```text
MemoryLimiter
RedisLimiter

MemoryAlgorithm
RedisAlgorithm

MemoryFixedWindow
RedisFixedWindow

MemoryTokenBucket
RedisTokenBucket
```

Because algorithm extension remains sealed for now ([[D67 - Sealed Redis Algorithms]]), the
interfaces should remain unexported in code unless we intentionally open them:

```text
memoryAlgorithm
redisAlgorithm
memoryAlgorithmState
```

Public concrete types can still be exported:

```text
MemoryLimiter
RedisLimiter
MemoryFixedWindow
RedisFixedWindow
MemoryTokenBucket
```

## Naming rule

Use qualifier before noun:

```text
MemoryLimiter
RedisLimiter
MemoryFixedWindow
RedisFixedWindow
```

Avoid noun-first names:

```text
LimiterMemory
LimiterRedis
FixedWindowRedis
AlgorithmMemory
```

Avoid short forms such as `InMem` in public API. Use `Memory` or `InMemory`. For this project,
prefer `Memory` because it matches `MemoryStore` and keeps names shorter.

## Config naming

Policy config should stay backend-neutral:

```text
fixedWindowConfig
tokenBucketConfig
```

These configs represent stable policy parameters, not execution backend details.

Backend behavior remains separate:

```text
MemoryFixedWindow -> Decide over Go state
RedisFixedWindow  -> Lua script and ARGV
```

## Tradeoffs

- **Readability:** better once multiple backends exist; type names say where execution happens.
- **Verbosity:** higher than short names like `Limiter` and `FixedWindow`.
- **API clarity:** better for users because memory and Redis paths are visibly different.
- **Refactoring cost:** medium; tests, docs, helpers, comments, and benchmarks must be renamed
  together.
- **Future flexibility:** leaves room for a generic `Limiter` interface later without name conflict.

## Deferred

Do not introduce a generic exported `Limiter` interface yet. Keep:

```text
MemoryLimiter
RedisLimiter
```

and revisit a shared interface once both backends have enough proven common behavior.

## Links

- [[D62 - Redis as Limiter Not Store]]
- [[D67 - Sealed Redis Algorithms]]
- [[D68 - Shared Fixed Window Config]]
- [[G39 - Go Naming Qualifier Order]]

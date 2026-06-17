# D68 - Shared Fixed Window Config

Back to [[V8 RedisLimiter Design Index]].

## Context

After adding `RedisFixedWindow`, the project has two implementations of the same fixed-window
policy:

- `MemoryFixedWindow` — in-process Go execution through `Decide`.
- `RedisFixedWindow` — Redis/Lua execution through `script()` and `args()`.

Both need the same policy configuration:

```text
limit
windowDuration
```

But they execute that policy through different substrates.

## Design question

Should one fixed-window type implement both the in-memory behavior and the Redis Lua behavior?

## Options

- **Single type implements both `Decide()` and `script()/args()`** — less duplication, but mixes
  pure Go decision logic with Redis/Lua execution details.
- **Fully separate types with duplicated config** — clean execution boundaries, but duplicate
  validation and config fields.
- **Shared private config, separate behavior** — one source of truth for stable policy config,
  while each backend keeps its own execution logic.

## Decision

Use shared private fixed-window config, but keep behavior separate.

Conceptually:

```text
fixedWindowConfig
  requestLimit
  windowDuration

MemoryFixedWindow
  config fixedWindowConfig
  Decide(...)

RedisFixedWindow
  config fixedWindowConfig
  script()
  args()
```

The shared config constructor owns validation once. The in-memory and Redis implementations stay
as sibling implementations of the fixed-window policy.

## Why

The stable part is the policy:

```text
allow N requests per window duration
```

The volatile part is execution:

- in-memory state structs,
- explicit `now`,
- CAS / mutex state commits,
- Redis TTL,
- Lua source,
- `KEYS` / `ARGV`,
- Redis result tuples,
- Redis script errors.

Sharing config avoids small duplication without making the pure Go algorithm know Redis details.

## Tradeoffs

- **Memory:** no meaningful difference; one nested struct replaces two direct fields.
- **Latency:** no meaningful difference; field access through config is trivial.
- **Concurrency:** no direct difference; execution guarantees remain backend-specific.
- **Maintainability:** validation becomes consistent across in-memory and Redis fixed window.
- **Boundary clarity:** better than one cross-backend fixed-window type because Redis execution
  remains isolated.
- **Complexity:** adds one private config type and constructor.

## Rule

Share stable policy. Separate volatile execution details.

## Links

- [[D31 - Algorithm Owns Config Validation]]
- [[D63 - Fixed Window in Redis via Lua]]
- [[D64 - RedisLimiter Algorithm Interface]]
- [[D67 - Sealed Redis Algorithms]]

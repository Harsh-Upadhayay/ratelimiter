# D66 - Redis Client Ownership Boundary

Back to [[V8 RedisLimiter Design Index]].

## Context

V8 has two layers:

- `RedisLimiter` — rate-limiting behavior.
- `redisAdapter` / `goRedisAdapter` — internal Redis execution plumbing.

After deciding that the adapter returns raw Redis output ([[D65 - Redis Adapter Returns Raw Result]]),
the next boundary question was: what should the public constructor accept?

## Options

- **Caller passes `*redis.Client`** — caller owns Redis connection setup and lifecycle; the limiter
  wraps it internally with `goRedisAdapter`.
- **Caller passes `RedisConfig`** — limiter creates the Redis client internally and hides the
  Redis library from the public API.
- **Caller passes a public adapter/interface** — caller implements a low-level `Eval`-style
  contract and hands that to the limiter.
- **Caller passes `goRedisAdapter`** — exposes internal plumbing as public API; not preferred.

## Decision

For now, the public constructor should accept an existing Redis client and wrap it internally:

```text
public constructor: existing Redis client + redis algorithm
internal field:     redisAdapter
internal wrapper:   goRedisAdapter
```

The adapter remains private. It is an internal seam, not the normal user-facing API.

## Why this is Go-idiomatic

Go libraries commonly let callers own heavyweight infrastructure clients:

- the application owns configuration,
- the application owns connection lifetime,
- the application can share one client across packages,
- the library focuses on the domain behavior it provides.

For this project:

```text
caller owns Redis connectivity
RedisLimiter owns rate limiting
goRedisAdapter hides go-redis Eval mechanics
```

This keeps Redis connection concerns out of the limiter without making users learn the limiter's
internal adapter boundary.

## System design framing

This decision is not "do we expose Redis complexity or hide it?" The real question is:

```text
Who owns the Redis connection lifecycle?
```

If the limiter creates the client from config, the limiter also becomes responsible for:

- closing the client,
- TLS/auth/DB/pool configuration,
- retry and timeout policy,
- cluster/sentinel support,
- observability hooks,
- avoiding accidental connection proliferation.

Those concerns belong naturally to the application boundary in mature services.

Passing an existing client keeps ownership explicit and avoids hidden infrastructure behavior.

## What the adapter still buys us

Keeping `redisAdapter` private still has value:

- `RedisLimiter.Allow` does not know the exact go-redis call shape.
- tests can fake Redis without running Redis.
- future client swaps are localized behind another adapter.

The important distinction:

```text
internal abstraction != public API
```

An internal adapter can be useful without becoming something users assemble directly.

## Future client support

If a future Redis client library is supported, add another internal adapter:

```text
goRedisAdapter
rueidisAdapter
redigoAdapter
```

Then expose a deliberate constructor for that client shape if needed.

Do not prematurely export a generic `RedisExecutor` interface. Once exported, that interface
becomes a public contract. It should only be introduced when there is real pressure from multiple
client implementations or tests that cannot use the private seam.

## Tradeoffs

- **Memory:** client sharing is efficient; the limiter does not create hidden extra connection
  pools.
- **Latency:** no meaningful overhead beyond one adapter method call.
- **Concurrency:** Redis client concurrency behavior remains owned by the client library and app
  configuration. `RedisLimiter` remains stateless apart from references to the client adapter and
  algorithm.
- **Coupling:** public API couples to a Redis client type. This is accepted for now to keep
  connection ownership explicit.
- **User experience:** users must create/configure Redis themselves, but they do not interact with
  Lua, `Eval`, `KEYS`, `ARGV`, or result parsing.
- **Extensibility:** supporting another Redis client later needs another adapter and likely another
  constructor, but `RedisLimiter.Allow` should not change.

## Links

- [[D62 - Redis as Limiter Not Store]]
- [[D63 - Fixed Window in Redis via Lua]]
- [[D65 - Redis Adapter Returns Raw Result]]
- [[G38 - context.Context]]
- [[redis/R03 - Lua Scripting for Atomicity]]
- [[redis/R05 - KEYS and ARGV Parameters]]

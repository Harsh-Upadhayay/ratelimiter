# ratelimiter

[![Go Reference](https://pkg.go.dev/badge/github.com/Harsh-Upadhayay/ratelimiter.svg)](https://pkg.go.dev/github.com/Harsh-Upadhayay/ratelimiter)
![Go 1.22+](https://img.shields.io/badge/go-1.22%2B-00ADD8)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A rate-limiting library for Go with pluggable algorithms, an in-process backend
(with optional lock-striping) and a Redis backend, plus a `net/http` middleware.

Built as a deliberate study of Go and distributed-systems design — every
non-trivial decision is recorded as an ADR (see [Design journal](#design-journal)).

> **Status:** active learning project. The in-process path is implemented,
> tested, race-checked, and benchmarked. The Redis path is implemented;
> integration tests, a runnable example, and CI are tracked in
> [`ROADMAP.md`](ROADMAP.md). See [Limitations](#limitations) before using in
> production.

## Features

- **One interface, two backends** — `Limiter.Allow(ctx, key) (Result, error)`,
  satisfied directly by both `MemoryLimiter` and `RedisLimiter`.
- **Two atomicity strategies** — in-process optimistic **CAS retry loop** for the
  memory backend; **single-shot Lua** scripts for Redis (server-side atomicity, no
  retry loop).
- **Pluggable algorithms** — Fixed Window and Token Bucket today; Sliding Window
  Counter/Log planned.
- **Hot-key sharding** — `ShardedMemoryStore` lock-strips across N `MemoryStore`
  shards (FNV-1a), added because benchmarks showed a single mutex doesn't scale
  under many-key load.
- **HTTP middleware** — drop-in `net/http` wrapper with a caller-supplied key
  function, fail-open / fail-closed policy, and `Retry-After` / `X-RateLimit-Remaining`
  headers.
- **Clock-injectable & skew-aware** — the memory backend uses an injectable clock
  for deterministic tests; the Redis backend uses Redis `TIME` as the script clock.

## Install

```sh
go get github.com/Harsh-Upadhayay/ratelimiter
```

Requires Go 1.22+.

## Quick start

### In-process (memory)

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Harsh-Upadhayay/ratelimiter"
)

func main() {
	// 100 requests per minute, per key, in-process.
	algo, err := ratelimiter.NewMemoryFixedWindow(100, time.Minute)
	if err != nil {
		panic(err)
	}
	limiter, err := ratelimiter.NewMemoryLimiter(algo, ratelimiter.NewMemoryStore())
	if err != nil {
		panic(err)
	}

	res, err := limiter.Allow(context.Background(), "user-123")
	if err != nil {
		panic(err)
	}
	if res.Allowed {
		fmt.Printf("allowed — %d remaining\n", res.Remaining)
	} else {
		fmt.Printf("rejected — retry after %s\n", res.RetryAfter)
	}
}
```

For high-cardinality / hot-key workloads, swap the store for a sharded one:

```go
store, _ := ratelimiter.NewShardedMemoryStore(256)
limiter, _ := ratelimiter.NewMemoryLimiter(algo, store)
```

### Redis (distributed)

```go
import "github.com/redis/go-redis/v9"

client := redis.NewClient(&redis.Options{Addr: "localhost:7777"})

// Token bucket: capacity 100, refill 10 tokens/sec.
algo, _ := ratelimiter.NewRedisTokenBucket(100, 10)
limiter, _ := ratelimiter.NewRedisLimiter(client, algo)

res, err := limiter.Allow(ctx, "user-123")
```

The app owns the Redis client's lifecycle and configuration; the limiter owns the
rate-limiting behavior. A local Redis for development is provided in
[`docker-compose.yml`](docker-compose.yml) (`docker compose up -d`, listens on
`:7777`).

### HTTP middleware

```go
mw, _ := ratelimiter.NewMiddleware(
	limiter,
	func(r *http.Request) string { return r.Header.Get("X-API-Key") }, // KeyFunc
	ratelimiter.WithFailurePolicy(ratelimiter.FailClosed),              // optional
)

http.Handle("/api/", mw.Wrap(apiHandler))
```

On a rejected request the middleware responds `429 Too Many Requests` with
`Retry-After`; every response carries `X-RateLimit-Remaining`. An empty key yields
`400`. On a limiter error the response is governed by the failure policy
(`FailOpen` passes the request through; `FailClosed` returns `503`).

## Algorithms

| Algorithm              | Memory | Redis | Notes                                              |
| ---------------------- | :----: | :---: | -------------------------------------------------- |
| Fixed Window           |   ✅   |  ✅   | Simplest; allows bursts at window edges.           |
| Token Bucket           |   ✅   |  ✅   | Smooths bursts; lazy refill, fractional tokens.    |
| Sliding Window Counter |   🔜   |  🔜   | Planned — see [`ROADMAP.md`](ROADMAP.md).          |
| Sliding Window Log     |   🔜   |  🔜   | Planned — Redis sorted sets.                       |

## Design

The core split is **storage from logic**:

- Algorithms are **pure functions** — `Decide(now, state, exists) → (Result, newState, error)`.
  No `time.Now()`, no I/O, no side effects, which makes them trivially testable.
- `StateStore` owns persistence — `Get` and `CompareAndSwap`. The memory backend
  implements it with a mutex; the sharded backend lock-strips it.
- `MemoryLimiter` orchestrates `Get → Decide → CompareAndSwap` in a **bounded
  optimistic-CAS retry loop**; on conflict it re-runs `Decide`, not just the swap.

Redis is **not** a `StateStore`. An early pivot ([ADR-0062](docs/adr/V7%20-%20Sharded%20MemoryStore.md))
established that the in-process abstraction passes live Go structs, while Redis
stores bytes and owns atomicity server-side. So `RedisLimiter` is a parallel
implementation that satisfies the same `Limiter` interface via **single-shot Lua
scripts** — pessimistic, server-side atomicity that replaces the in-process CAS
loop entirely.

```
        ┌──────────────┐        ┌──────────────────────────┐
HTTP ──▶ │  Middleware  │ ─────▶ │ Limiter.Allow(ctx, key)  │
        └──────────────┘        └────────────┬─────────────┘
                                  ┌───────────┴────────────┐
                                  ▼                        ▼
                       MemoryLimiter                 RedisLimiter
                 Get→Decide→CAS (retry loop)      single-shot Lua (EVAL)
                          │                                │
                   StateStore (mutex /              redis.Client
                    lock-striped shards)         (server-side atomic)
```

## Benchmarks

The sharded store exists *because* of measurement, not speculation: benchmarks
showed a single `MemoryStore` mutex does not scale under parallel many-key load
([ADR-0048](docs/adr/V7%20-%20Sharded%20MemoryStore.md)).

```sh
go test -bench=. -benchmem
```

Benchmarks live in `memory_limiter_benchmark_test.go` (same-key vs many-key,
serial vs parallel, for both algorithms).

## Limitations

This is a single-node-correct, contention-optimized, multi-algorithm limiter — not
a failover-correct, internet-scale one. Known boundaries (intentional — see
[`ROADMAP.md`](ROADMAP.md) non-goals):

- **Single Redis node only** — no replication/cluster story; a primary failover to a
  stale replica could double-count.
- **Every request is a Redis round-trip** — no local token leasing / async
  reconciliation (the GCRA-style hot-path optimization).
- **Unbounded memory** — `MemoryStore` keys live forever; no eviction (Redis has TTL).

## Testing

```sh
go build ./...
go test ./...
go test -race ./...
go vet ./...
```

## Design journal

The reasoning behind the architecture is documented as ADRs, one log per version:

- [`docs/adr/README.md`](docs/adr/README.md) — index of all decisions (ADR-0001…).
- [`docs/concepts/Go Concepts.md`](docs/concepts/Go%20Concepts.md) /
  [`docs/concepts/Redis Concepts.md`](docs/concepts/Redis%20Concepts.md) — the
  Go and Redis concepts each decision rests on.
- [`docs/Rate Limiter Learning Map.md`](docs/Rate%20Limiter%20Learning%20Map.md) — entry point.
- [`plan.md`](plan.md) — the original north-star plan (since revised; see the note at top).

## Roadmap & status

What's left before this is considered "done" lives in [`ROADMAP.md`](ROADMAP.md).

## License

[MIT](LICENSE) © Harsh Upadhayay

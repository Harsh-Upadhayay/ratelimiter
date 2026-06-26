# Roadmap

This is the **definition of done** for the project — what's left before it's
marked complete as a portfolio piece. Items use GitHub task-list syntax so they
render as checkboxes here and can be lifted into Issues + a Milestone if desired.

Priorities:

- **P0 — Blockers:** things that actively misrepresent the project to a reader.
- **P1 — Required:** the credibility floor (tests, example, CI).
- **P2 — Differentiators:** pick at least one.
- **Non-goals:** deliberately out of scope; documented so the scoping reads as a choice.

The bar for "complete": **every P0 and P1 box checked, at least one P2 box checked,
and the non-goals left explicit.**

The P1/P2 items below are mirrored as GitHub issues
[#1–#14](https://github.com/Harsh-Upadhayay/ratelimiter/issues) under the
[`v0.1.0`](https://github.com/Harsh-Upadhayay/ratelimiter/milestone/1) milestone.

---

## P0 — Blockers

- [x] Rewrite `README.md` as a real Go-library front door (replace the old AWS/k6
      "phased plan" that described a project that was never built).
- [x] Reconcile stale docs with shipped code:
  - [x] `docs/current_assessment.md` — the `MemoryLimiter` / `Limiter` "tension" is
        resolved (it now satisfies `Limiter` directly via a private clock).
  - [x] `plan.md` — flag that steps 3–4 were revised by ADR-0062 (Redis is a
        parallel `RedisLimiter`, not a `StateStore`).
- [x] Rename `memory_token_bucket_tests.go` → `memory_token_bucket_test.go` so Go
      actually compiles it (its two tests now run).
- [x] Add a `LICENSE` file (MIT) and update the README License section.

## P1 — Required (credibility floor)

### Tests

- [ ] **Token-bucket behavior tests** (fill the file's existing TODO list): starts
      full, rejects when empty, refills over elapsed time, caps at capacity,
      sub-token rejection, retry-after math, backward clock grants nothing.
- [ ] **Redis integration tests** against the `docker-compose` Redis, behind a build
      tag / `-short` skip so unit runs stay fast. Both algorithms: allow/deny,
      refill, TTL expiry, retry-after, and a concurrency hammer proving Lua
      atomicity. *(Expect to surface a real bug in the token-bucket Lua on first run.)*
- [ ] **Middleware tests** for `Wrap`: allowed pass-through, 429 + `Retry-After`,
      `X-RateLimit-Remaining`, empty-key 400, fail-open vs fail-closed on limiter
      error, nil-next panic.
- [ ] **`state_store_test.go` TODO list** — the 5 CAS / concurrency cases.

### Runnable example

- [ ] `examples/` (or `cmd/server/`) — a small `net/http` server using
      `NewMiddleware`, wired to the docker-compose Redis, with README instructions
      to hit it and see `429`s.
- [ ] **k6 load script** producing a **Memory vs Redis latency comparison**; drop the
      numbers/graph into the README Benchmarks section.
- [ ] **Logical architecture diagram** in the README (the ASCII one can be replaced
      with a rendered image).

### CI & hygiene

- [ ] `.github/workflows/ci.yml` — `build`, `vet`, `test`, `test -race`, with a Redis
      service container for the integration tests. Add the green badge to the README.
- [ ] `golangci-lint` (or at least `gofmt -l`) clean.
- [ ] All gates green locally and in CI:
      `go build ./... && go vet ./... && go test ./... && go test -race ./...`.
- [ ] `CHANGELOG.md` (Keep a Changelog format) and a first tagged release (`v0.1.0`).

## P2 — Differentiators (pick ≥ 1)

- [ ] **Sliding Window Counter** (memory + Redis) — closes `plan.md` step 5; the
      canonical interview algorithm beyond the warm-up pair.
- [ ] **Metrics seam** — a small `Observer` / `Metrics` interface (allowed / denied /
      conflict counts); closes `plan.md` step 6.
- [ ] **Error-wrapping audit** — confirm `%w` + `errors.Is/As` is the documented
      contract; capture it in a `doc.go` package comment.
- [ ] **Package doc (`doc.go`)** so `pkg.go.dev` renders a proper overview.

## Non-goals (explicitly out of scope)

These were considered and deliberately cut — listed so the scope reads as a
decision, not an omission.

- [ ] ~~AWS deployment / CDK / ECS / ElastiCache / ALB~~ — re-proves cloud skills the
      résumé already covers; wrong form factor for a library. (The old README's plan.)
- [ ] ~~Redis Cluster / failover correctness~~ — single-node only.
- [ ] ~~Local token leasing / latency-aware hot path (GCRA)~~.
- [ ] ~~In-process eviction policy for `MemoryStore`~~ — unbounded by design.
- [ ] ~~OpenTelemetry → CloudWatch~~ — a local metrics seam (P2) is sufficient.

---

> Tip: to make this GitHub-native, convert each P0/P1 item into an Issue and group
> them under a `v0.1.0` Milestone; this file then becomes the high-level summary.

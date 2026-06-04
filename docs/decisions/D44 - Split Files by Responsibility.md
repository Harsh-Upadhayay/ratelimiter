# D44 - Split Files by Responsibility

Back to [[Rate Limiter Learning Map]].

## Context

The implementation started in one file to make basic Go syntax and design pressure visible. With fixed window, token bucket, shared result types, algorithm contracts, and sentinel errors, the single file now mixes too many responsibilities.

## Decision

Split files by responsibility while keeping one `ratelimiter` package.

## Proposed layout

```text
types.go
limiter.go
fixed_window.go
token_bucket.go
errors.go
```

## Ownership

- `types.go`: small shared contracts such as `Result`, `algorithm`, and `algorithmState`.
- `limiter.go`: generic orchestration with `Limiter`, `NewLimiter`, and `Allow`.
- `fixed_window.go`: fixed-window config, state, constructor, and decision logic.
- `token_bucket.go`: token-bucket config, state, constructor, and decision logic.
- `errors.go`: sentinel errors.

## Why

This keeps concrete behavior close to the types that own it while centralizing only the small shared contracts.

## Tradeoff

`types.go` can become a junk drawer if too much is added. Keep concrete algorithm types out of it.

## Links

- [[D15 - One File First]]
- [[D26 - Introduce Algorithm Interface]]
- [[D35 - Token Bucket as Second Algorithm]]
- [[G25 - Package Scope Across Files]]

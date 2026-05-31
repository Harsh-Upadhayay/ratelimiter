# D14 - Package Name

Back to [[Rate Limiter Learning Map]].

## Context

The package name sets the import and usage style for the library.

## Decision

Use package name `ratelimiter`.

## Why

This project is a reusable learning library, not just an executable demo. The name is explicit and matches the repository purpose.

## Tradeoff

`ratelimiter.Limiter` is slightly repetitive, but clear.

## Alternatives

- `limiter`: shorter, but less specific outside this repository.
- `main`: useful for an executable demo, not for a reusable package.

## Links

- [[D12 - Exported API Boundary]]

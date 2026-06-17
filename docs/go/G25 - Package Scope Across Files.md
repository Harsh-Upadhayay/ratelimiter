# G25 - Package Scope Across Files

Back to [[Rate Limiter Learning Map]].

## Go practice

Files in the same directory with the same package name share package scope.

## Why it matters

Moving `MemoryFixedWindow` to `memory_fixed_window.go` and `MemoryLimiter` to `memory_limiter.go` does not require exporting internal types as long as both files use `package ratelimiter`.

## Rate limiter use

Private interfaces and state types can remain unexported after the file split.

## Tradeoff

File boundaries organize code for humans; they do not create visibility boundaries inside the same package.

## Links

- [[D44 - Split Files by Responsibility]]
- [[G08 - Exported Types with Unexported Fields]]

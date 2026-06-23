# G40 - Unexported Package Constants

Back to [[Go Concepts Index]].

## Concept

In Go, a constant declared at package scope is visible to every file in the same package.

```go
const tokenScale = 1000
```

Because the name starts with a lowercase letter, it is unexported. Code outside the package cannot
access it as `ratelimiter.tokenScale`.

## Why This Fits Here

`tokenScale` is an implementation detail shared by Redis Token Bucket methods:

- constructor validation,
- Redis argument conversion,
- Lua script argument meaning.

Package-level scope avoids duplicating the literal in every method while still keeping it hidden from
callers.

## Tradeoff

- **Local constant:** narrowest scope, but duplicate declarations if multiple methods need the same
  value.
- **Package-level unexported constant:** shared inside the package, hidden from callers.
- **Exported constant:** visible to callers, but exposes an implementation detail as public API.

The chosen shape is package-level and unexported.

## Links

- [[D73 - Redis Token Bucket Scaling Boundary]]
- [[G25 - Package Scope Across Files]]

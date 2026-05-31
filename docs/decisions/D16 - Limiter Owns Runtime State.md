# D16 - Limiter Owns Runtime State

Back to [[Rate Limiter Learning Map]].

## Context

V1 needs somewhere to keep shared configuration, all per-user state, and the lock protecting that state.

## Decision

Use an exported `Limiter` type that owns limit, window duration, user-state map, and mutex.

## Why

This allows multiple independent limiter instances, such as one policy for `/api/search` and another for `/login`.

## Tradeoff

This is more structure than package-level globals, but much easier to test and reason about.

## Alternatives

- Package-level globals: smallest code, poor test isolation and only one policy per process.
- Function with all state passed in: explicit, but awkward and easy to misuse.
- Per-user object only: still needs an owner for the outer map and creation synchronization.

## Links

- [[D12 - Exported API Boundary]]
- [[D15 - One File First]]
- [[G04 - Mutexes and Critical Sections]]

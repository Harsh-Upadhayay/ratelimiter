# D34 - Initialize Missing State Before Assertion

Back to [[Rate Limiter Learning Map]].

## Context

For a new key, map lookup returns no stored algorithm state. For an existing key, the algorithm expects a concrete state type.

## Decision

Handle missing state before asserting the concrete state type.

## Why

Missing state is normal initialization. Existing state of the wrong concrete type is an error.

## Tradeoff

Each algorithm has an explicit initialization branch before its type assertion.

## Links

- [[D29 - Exists Flag for Missing State]]
- [[D33 - Decide Returns State Errors]]
- [[G20 - Type Assertions]]

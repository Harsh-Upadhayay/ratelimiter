# D19 - Result and Error Return

Back to [[Rate Limiter Learning Map]].

## Context

The V1 `Allow` method returns only a boolean. That is enough for a local decision, but not enough for callers that need retry timing or quota details.

## Decision

For V2, evolve `Allow` to return a result struct and an error.

## Why

A result struct can carry rate-limit decision details, while an error can represent invalid input or operational failure.

## Tradeoff

This is a larger API change than returning only a result. It forces callers and tests to handle two return values.

## Design pressure

For an in-memory limiter, errors are mostly validation errors. Later, a distributed store can produce real operational errors such as Redis failures or CAS retry exhaustion.

## Contract

- Allowed request: return populated result and `nil` error.
- Rate-limited request: return rejected result and `nil` error.
- Invalid key: return zero result and non-nil error.

Rate limiting is a normal decision, not an error.

## Links

- [[D10 - Return Type]]
- [[D18 - Empty Key Validation]]
- [[D21 - Result Contract]]
- [[D22 - Sentinel Error for Empty Key]]
- [[G11 - Multiple Return Values]]
- [[G01 - Structs for Grouped State]]

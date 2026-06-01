# D18 - Empty Key Validation

Back to [[Rate Limiter Learning Map]].

## Context

The limiter uses the user ID as the map key for per-user state.

## Decision

For V1, panic if `Allow` receives an empty user ID.

## Why

An empty key is almost certainly a caller bug. Failing immediately keeps the V1 behavior explicit and consistent with constructor validation.

## Tradeoff

Panic is harsh for reusable library code. A production API may prefer returning an error or treating key extraction failure before calling the limiter.

## Revisit when

When adding HTTP middleware or a richer result/error API.

## V2 update

V2 replaces this panic with a returned sentinel error. See [[D22 - Sentinel Error for Empty Key]].

## Links

- [[D10 - Return Type]]
- [[D11 - Constructor Validation]]
- [[D12 - Exported API Boundary]]
- [[D22 - Sentinel Error for Empty Key]]

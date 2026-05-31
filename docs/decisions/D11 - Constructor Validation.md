# D11 - Constructor Validation

Back to [[Rate Limiter Learning Map]].

## Context

The limiter needs valid configuration before it can make correct decisions.

## Decision

For V1, panic if the constructor receives an invalid limit or window.

## Why

This keeps the implementation small while still preserving consistent behavior. Invalid configuration fails immediately instead of creating a broken limiter.

## Tradeoff

Panic is harsh for reusable library code because callers cannot handle the error normally.

## Future direction

In the next iteration, return a typed error or sentinel error for invalid configuration.

## Links

- [[D10 - Return Type]]
- [[G01 - Structs for Grouped State]]

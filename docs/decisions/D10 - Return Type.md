# D10 - Return Type

Back to [[Rate Limiter Learning Map]].

## Context

The allow decision can return only allowed/rejected, or it can return operational detail such as remaining quota and retry timing.

## Decision

For V1, return only a boolean.

## Why

The first implementation should focus on correct state transitions, locking, and time-window behavior.

## Tradeoff

A boolean does not give HTTP middleware enough information to produce helpful headers or a precise `Retry-After` value.

## Revisit when

When adding an HTTP integration or when comparing algorithms that can calculate different retry timing semantics.

## Future likely shape

A richer result may include:

- allowed/rejected
- remaining quota
- retry-after duration
- reset time

## Previous question

Should V1 return only a boolean, or a richer result object?

## Alternatives to evaluate

- Boolean only: smallest API, easy to learn, little information for callers.
- Result object: more useful to HTTP callers, more fields and invariants to maintain.

## Tradeoff dimensions

- API simplicity
- Test clarity
- HTTP response usefulness
- Future compatibility with multiple algorithms

## Links

- [[D06 - Count Invariant Accepted Requests]]
- [[G01 - Structs for Grouped State]]

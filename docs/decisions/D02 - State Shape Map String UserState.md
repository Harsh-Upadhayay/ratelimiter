# D02 - State Shape Map String UserState

Back to [[Rate Limiter Learning Map]].

## Context

Each user needs enough state to decide whether another request should be allowed in the current window.

## Decision

Use a map keyed by user ID, where each value stores request count and window start.

## Why

The user ID is naturally represented as a string, and the per-user fields belong together as one state object.

## Tradeoff

Map values are copied when read. If the value is updated locally, it must be written back to the map.

## Revisit when

When comparing `map[string]State` with `map[string]*State`, especially around mutation clarity, memory use, and lock ownership.

## Links

- [[G01 - Structs for Grouped State]]
- [[G02 - Maps Comma Ok and Value Copies]]
- [[D06 - Count Invariant Accepted Requests]]

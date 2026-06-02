# D30 - Keep Algorithm Interface Private

Back to [[Rate Limiter Learning Map]].

## Context

The algorithm abstraction is new and will evolve as token bucket is added.

## Decision

Keep the algorithm interface unexported for now.

## Why

The package can refactor the interface while learning from two concrete algorithms without breaking external callers.

## Tradeoff

External packages cannot provide custom algorithms yet.

## Revisit when

After multiple internal algorithms stabilize the contract and there is a real external-extension use case.

## Links

- [[D26 - Introduce Algorithm Interface]]
- [[D27 - Algorithm Owned State]]
- [[G18 - Structural Interface Satisfaction]]

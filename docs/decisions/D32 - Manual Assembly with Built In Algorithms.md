# D32 - Manual Assembly with Built In Algorithms

Back to [[Rate Limiter Learning Map]].

## Context

The algorithm interface remains private while it stabilizes, but callers should explicitly choose and assemble a configured built-in algorithm with a limiter.

## Decision

Use manual assembly with package-provided algorithms only.

## Conceptual flow

```text
configured built-in algorithm -> generic limiter -> Allow
```

## Why

This exposes composition without committing to public custom-algorithm support yet.

## Tradeoff

External packages cannot implement custom algorithms. The package keeps freedom to refactor the private contract while fixed window and token bucket reveal its final shape.

## Links

- [[D30 - Keep Algorithm Interface Private]]
- [[D31 - Algorithm Owns Config Validation]]

# G18 - Structural Interface Satisfaction

Back to [[Rate Limiter Learning Map]].

## Go practice

Go interfaces describe method sets.

A type satisfies an interface automatically when it has the required methods. There is no `implements` keyword.

## Why it matters

The limiter can depend on algorithm behavior without depending on a concrete algorithm type.

## Tradeoff

Interface compatibility is implicit. Keep interfaces small so their purpose remains easy to understand.

## Links

- [[D26 - Introduce Algorithm Interface]]
- [[G17 - Interfaces From Real Variation]]

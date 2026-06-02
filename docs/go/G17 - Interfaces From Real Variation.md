# G17 - Interfaces From Real Variation

Back to [[Rate Limiter Learning Map]].

## Go practice

In Go, interfaces are strongest when they come from real variation in behavior.

## Why

With one implementation, an interface is usually a guess. With two implementations, the shared behavior and differing details become clearer.

## Rate limiter use

Fixed window and token bucket both make rate-limit decisions, but they use different state and configuration. That creates real pressure for an algorithm interface.

## Tradeoff

Introducing an interface increases abstraction and indirection. The benefit should be a cleaner boundary between limiter orchestration and algorithm decision logic.

## Links

- [[D13 - Delay Interfaces]]
- [[D26 - Introduce Algorithm Interface]]

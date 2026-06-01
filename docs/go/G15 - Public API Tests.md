# G15 - Public API Tests

Back to [[Rate Limiter Learning Map]].

## Go practice

Tests can live in the same package and still choose to exercise only exported behavior.

## Why

Testing through the public API makes refactors safer because private implementation details can change without breaking tests.

## Tradeoff

Private helpers may be harder to target directly, but behavior remains covered if all important outcomes are observable through the public API.

## Links

- [[D24 - Test Through Public API]]
- [[G14 - Pure Helper Functions]]

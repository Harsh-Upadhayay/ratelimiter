# G16 - Helper Parameter Ordering

Back to [[Rate Limiter Learning Map]].

## Go practice

When a helper takes several parameters, group them by meaning so the call site is readable.

## Current shape

The fixed-window helper uses this order:

- time input
- configuration inputs
- state input
- state metadata

## Why it matters

Go does not have named arguments. Parameter order is part of readability, so mixed ordering can make call sites harder to audit.

## Tradeoff

Many parameters can be a smell. For now this is acceptable because introducing a config struct or interface would be premature.

## Links

- [[D23 - Private Fixed Window Decision Helper]]
- [[G14 - Pure Helper Functions]]

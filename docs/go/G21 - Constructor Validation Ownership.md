# G21 - Constructor Validation Ownership

Back to [[Rate Limiter Learning Map]].

## Go practice

Constructors should reject invalid configuration before returning a usable object.

## Design rule

Validation belongs with the type that owns the invariant.

## Rate limiter use

- Algorithm constructors validate algorithm-specific configuration.
- The generic limiter constructor validates limiter-level dependencies.

## Why it matters

This avoids repeated hot-path validation and prevents invalid configured components from existing.

## Links

- [[D31 - Algorithm Owns Config Validation]]
- [[D32 - Manual Assembly with Built In Algorithms]]

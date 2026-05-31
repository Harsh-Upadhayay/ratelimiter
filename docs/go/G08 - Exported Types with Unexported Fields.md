# G08 - Exported Types with Unexported Fields

Back to [[Rate Limiter Learning Map]].

## Go practice

In Go, package visibility is controlled by capitalization.

An exported type can still have unexported fields.

## Pointer note

Returning a pointer to a value does not let callers from another package access unexported fields.

For example, code outside the package can use exported methods on `*Limiter`, but it cannot access lowercase fields such as `limit`, `window`, `states`, or `mu`.

## Why it matters

This lets a package expose behavior while hiding representation.

For the limiter, callers can ask whether a request is allowed, but they cannot replace the state map, change the limit directly, or manipulate the mutex.

## Design implication

Export the stable API surface and keep internal state private.

## Links

- [[D12 - Exported API Boundary]]
- [[D17 - Constructor Returns Pointer]]

# G10 - Early Returns and Guard Clauses

Back to [[Rate Limiter Learning Map]].

## Go practice

Go code commonly uses early returns and guard clauses.

## Why

If a branch has reached a final decision, returning immediately keeps the rest of the function flat and easier to audit.

## Why this fits Go

- Go has explicit error handling, so guard clauses are common.
- `defer` handles cleanup across multiple return paths.
- Go style generally favors simple, flat control flow over nested branches.

## Tradeoff

Single-return functions can be useful when one final result is assembled through several steps, but they can also introduce extra mutable variables and nesting.

## Rate limiter example

The allow decision is naturally branch-based:

- missing user means allow
- expired window means allow
- available quota means allow
- exhausted quota means reject

Each branch can return its decision directly.

## Links

- [[G05 - Defer Unlock Pattern]]
- [[D10 - Return Type]]

# D22 - Sentinel Error for Empty Key

Back to [[Rate Limiter Learning Map]].

## Context

V1 panicked for an empty key. V2 introduces an error return, so invalid input can be represented without crashing the caller.

## Decision

Define an exported sentinel error for empty keys.

## Why

Empty key is a stable public error category. A sentinel lets callers and tests compare with `errors.Is`.

## Tradeoff

Exported errors become part of the public API and should remain stable.

## Naming

Use generic key language, such as `ErrEmptyKey`, instead of user-specific language like `ErrEmptyUserID`.

## Links

- [[D18 - Empty Key Validation]]
- [[D20 - Generic Rate Limit Key]]
- [[G12 - Sentinel Errors]]

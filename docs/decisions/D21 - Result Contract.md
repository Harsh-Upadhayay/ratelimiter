# D21 - Result Contract

Back to [[Rate Limiter Learning Map]].

## Context

The boolean V1 API only says allow or reject. Callers need more detail for HTTP headers and retry behavior.

## Decision

Introduce an exported `Result` with `Allowed`, `Remaining`, and `RetryAfter`.

## Semantics

- `Allowed`: whether this request was accepted.
- `Remaining`: allowed requests left in the current window after this decision.
- `RetryAfter`: zero for allowed requests; time until reset for rejected requests.

## Fixed-window retry calculation

For a rejected request:

```text
RetryAfter = windowStartTime + windowDuration - now
```

## Tradeoff

The API becomes more useful but every branch must maintain result-field invariants.

## Links

- [[D19 - Result and Error Return]]
- [[G13 - Exported Result Structs]]
- [[G06 - Time and Duration Boundaries]]

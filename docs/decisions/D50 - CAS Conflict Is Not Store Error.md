# D50 - CAS Conflict Is Not Store Error

Back to [[Rate Limiter Learning Map]].

## Context

`CompareAndSwap` can fail because another caller updated the key after `Get`.

## Decision

Represent a normal CAS conflict as:

```text
ok = false
err = nil
```

Reserve `err` for real store failures.

## Why

A version mismatch is an expected concurrency outcome. The limiter can retry by reading fresh state and running the algorithm again.

## Tradeoff

Callers must check both return values. This is slightly more verbose, but it keeps operational failures distinct from retryable conflicts.

## Links

- [[D47 - StateStore Uses Get and CAS]]
- [[D51 - Bounded CAS Retry Loop]]
- [[G11 - Multiple Return Values]]
- [[G26 - Optimistic Concurrency with CAS]]

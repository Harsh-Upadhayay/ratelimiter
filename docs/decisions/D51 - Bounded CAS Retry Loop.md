# D51 - Bounded CAS Retry Loop

Back to [[Rate Limiter Learning Map]].

## Context

With `Get` plus `CompareAndSwap`, the limiter may lose a race to another request and need to retry.

## Decision

Use a bounded retry loop in `Allow`.

After retry exhaustion, return `ErrCASConflict`.

## Why

Request-path code should not retry forever under hot-key contention.

A bounded retry keeps tail latency finite and gives the caller an operational error it can handle with a fail-open or fail-closed policy.

## Tradeoff

Bounded retries can return an unknown operational result even when capacity may exist. Infinite retries are simpler logically, but can hang a request under heavy contention.

## Links

- [[D47 - StateStore Uses Get and CAS]]
- [[D50 - CAS Conflict Is Not Store Error]]
- [[G26 - Optimistic Concurrency with CAS]]
- [[G12 - Sentinel Errors]]

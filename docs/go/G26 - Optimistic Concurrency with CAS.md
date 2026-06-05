# G26 - Optimistic Concurrency with CAS

Back to [[Rate Limiter Learning Map]].

## Concept

CAS means compare and swap.

It commits a new value only if the stored version still matches the version that was read.

## Rate limiter use

CAS protects the read-decide-write sequence when multiple requests may update the same key concurrently.

## Why it matters

Without CAS, two callers can both make decisions from stale state and overwrite each other.

## Tradeoff

CAS avoids holding a long lock, but it introduces retry loops and conflict handling.

## Links

- [[D47 - StateStore Uses Get and CAS]]
- [[D46 - GetSet Race]]

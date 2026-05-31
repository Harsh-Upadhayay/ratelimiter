# G07 - Do Not Copy Mutexes

Back to [[Rate Limiter Learning Map]].

## Go practice

Do not copy synchronization primitives such as `sync.Mutex` after first use.

## Why

A mutex contains synchronization bookkeeping. Copying it duplicates the lock value but not a meaningful shared synchronization relationship.

## Failure mode

If a struct with a mutex and a map is copied, both structs may still point to the same map but have different mutexes.

That means two goroutines can lock different mutexes and mutate the same map concurrently.

## Locked-copy risk

If a mutex is copied while locked, the copy may remain locked forever because the unlock happens on the original mutex, not the copy.

## Design implication

For stateful types containing a mutex, prefer pointer usage and avoid APIs that encourage copying.

## Links

- [[D17 - Constructor Returns Pointer]]
- [[G04 - Mutexes and Critical Sections]]

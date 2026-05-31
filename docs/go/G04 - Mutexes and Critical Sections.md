# G04 - Mutexes and Critical Sections

Back to [[Rate Limiter Learning Map]].

## Go practice

Use `sync.Mutex` to protect shared mutable state accessed by multiple goroutines.

## Critical section

The critical section is not just the map read. It includes the full read, decide, update sequence.

## Why it matters

If two goroutines both read the same old count before either writes, they can both allow a request that should not both pass.

## Links

- [[D04 - One Global Mutex]]
- [[G05 - Defer Unlock Pattern]]

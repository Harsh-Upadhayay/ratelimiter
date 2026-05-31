# V1 Learning Mindmap

Back to [[Rate Limiter Learning Map]].

```mermaid
mindmap
  root((Rate Limiter V1))
    Goal
      In-memory fixed window limiter
      Learn Go syntax through a small stateful system
      Delay distributed design until basics are clear
      package ratelimiter
      one file first
      panic on invalid constructor input and empty key
    State
      Limiter owns runtime state
      exported type with unexported fields
      map from user ID to per-user state
      count of accepted requests
      window start time
      constructor initializes map
    Concurrency
      one global mutex
      critical section covers read, decide, update
      defer unlock after lock
      do not copy mutexes
      constructor returns pointer
    Time
      pass now into the method
      reset expired windows to now
      half-open interval
    Future Iterations
      Algorithm interface
      StateStore interface
      richer result type
      constructor errors
      per-key locks
      Redis and CAS
```

## Linked notes

- [[D01 - Fixed Window First]]
- [[D04 - One Global Mutex]]
- [[D05 - Explicit Time Input]]
- [[D09 - Defer Unlock After Lock]]
- [[D10 - Return Type]]

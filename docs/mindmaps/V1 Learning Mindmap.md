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
    V2 API
      Allow returns Result and error
      generic key instead of user ID
      ErrEmptyKey sentinel error
      rate-limit rejection is not an error
      RetryAfter only matters when rejected
    V3 Refactor
      private pure decision helper
      Allow remains orchestration
      test through public API
      helper parameters grouped by meaning
    V4 Algorithm Boundary
      second algorithm creates real variation
      introduce Algorithm interface
      algorithm owns state shape
      marker interface narrows valid states
      explicit exists flag for missing state
      keep algorithm interface private initially
      algorithm owns config validation
      manual assembly with built-in algorithms
      Decide returns state mismatch errors
    V5 Token Bucket
      capacity is burst size
      refill rate is sustained rate
      lazy refill on request
      float64 tokens internally
      Remaining exposes whole requests
      clamp negative elapsed time
      new buckets start full
      split files by responsibility
    V6 Storage Boundary
      limiter should not own map storage
      simple Get and Set can race
      Get plus CAS preserves read decide write
      retry on version conflict
```

## Linked notes

- [[V1 - Fixed Window#ADR-0001 — Fixed window first|ADR-0001]]
- [[V1 - Fixed Window#ADR-0004 — One global mutex|ADR-0004]]
- [[V1 - Fixed Window#ADR-0005 — Explicit time input|ADR-0005]]
- [[V1 - Fixed Window#ADR-0009 — Defer unlock immediately after lock|ADR-0009]]
- [[V1 - Fixed Window#ADR-0010 — Boolean return for V1|ADR-0010]]
- [[V2 - API Evolution#ADR-0019 — Result-and-error return|ADR-0019]]
- [[V2 - API Evolution#ADR-0020 — Generic rate-limit key|ADR-0020]]
- [[V2 - API Evolution#ADR-0021 — `Result` contract|ADR-0021]]
- [[V2 - API Evolution#ADR-0022 — Sentinel error for empty key|ADR-0022]]

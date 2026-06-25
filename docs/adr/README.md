# Architecture Decision Records

> Back to [[Rate Limiter Learning Map]]

Architecture decisions for the rate limiter, in [Nygard ADR format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) (Status / Context / Decision / Consequences). ADRs are grouped into one log per iteration (V1–V9); each log opens with a short narrative of that iteration. ADR numbers are stable and append-only — a decision is never edited away, only superseded by a later ADR.

The Go and Redis/Lua *concepts* these decisions rely on live separately in [[Go Concepts]] and [[Redis Concepts]].

## Logs

- [[V1 - Fixed Window]] — ADR-0001 … ADR-0018
- [[V2 - API Evolution]] — ADR-0019 … ADR-0022
- [[V3 - Decision Logic]] — ADR-0023 … ADR-0025
- [[V4 - Algorithm Abstraction]] — ADR-0026 … ADR-0034
- [[V5 - Token Bucket]] — ADR-0035 … ADR-0044
- [[V6 - Storage Boundary]] — ADR-0045 … ADR-0052
- [[V7 - Sharded MemoryStore]] — ADR-0053 … ADR-0062
- [[V8 - RedisLimiter]] — ADR-0063 … ADR-0073
- [[V9 - HTTP Middleware]] — ADR-0074 … ADR-0085

## Index

| ADR | Title | Status | Log |
|-----|-------|--------|-----|
| 0001 | Fixed window first | Accepted | [[V1 - Fixed Window]] |
| 0002 | Per-key state as map of structs | Accepted | [[V1 - Fixed Window]] |
| 0003 | Pointer receiver for the mutating limiter | Accepted | [[V1 - Fixed Window]] |
| 0004 | One global mutex | Superseded by 0045, 0053 | [[V1 - Fixed Window]] |
| 0005 | Explicit time input | Superseded by 0083 | [[V1 - Fixed Window]] |
| 0006 | Count means accepted requests | Accepted | [[V1 - Fixed Window]] |
| 0007 | Window resets at now | Accepted | [[V1 - Fixed Window]] |
| 0008 | Half-open window interval | Accepted | [[V1 - Fixed Window]] |
| 0009 | Defer unlock immediately after lock | Accepted | [[V1 - Fixed Window]] |
| 0010 | Boolean return for V1 | Superseded by 0019, 0021 | [[V1 - Fixed Window]] |
| 0011 | Panic on invalid config | Superseded by 0031 | [[V1 - Fixed Window]] |
| 0012 | Exported API, hidden state | Accepted | [[V1 - Fixed Window]] |
| 0013 | Delay interfaces | Accepted | [[V1 - Fixed Window]] |
| 0014 | Package name `ratelimiter` | Accepted | [[V1 - Fixed Window]] |
| 0015 | One file first | Superseded by 0044 | [[V1 - Fixed Window]] |
| 0016 | Limiter owns runtime state | Accepted | [[V1 - Fixed Window]] |
| 0017 | Constructor returns a pointer | Accepted | [[V1 - Fixed Window]] |
| 0018 | Panic on empty key | Superseded by 0022 | [[V1 - Fixed Window]] |
| 0019 | Result-and-error return | Accepted | [[V2 - API Evolution]] |
| 0020 | Generic rate-limit key | Accepted | [[V2 - API Evolution]] |
| 0021 | `Result` contract | Accepted | [[V2 - API Evolution]] |
| 0022 | Sentinel error for empty key | Accepted | [[V2 - API Evolution]] |
| 0023 | Private fixed-window decision helper | Accepted | [[V3 - Decision Logic]] |
| 0024 | Test through the public API | Accepted | [[V3 - Decision Logic]] |
| 0025 | `Allow` as orchestrator | Accepted | [[V3 - Decision Logic]] |
| 0026 | Introduce the algorithm interface | Accepted | [[V4 - Algorithm Abstraction]] |
| 0027 | Algorithm-owned state | Accepted | [[V4 - Algorithm Abstraction]] |
| 0028 | Marker interface for algorithm state | Accepted | [[V4 - Algorithm Abstraction]] |
| 0029 | Explicit `exists` flag for missing state | Accepted | [[V4 - Algorithm Abstraction]] |
| 0030 | Keep the algorithm interface private | Accepted | [[V4 - Algorithm Abstraction]] |
| 0031 | Algorithm owns config validation | Accepted | [[V4 - Algorithm Abstraction]] |
| 0032 | Manual assembly with built-in algorithms | Accepted | [[V4 - Algorithm Abstraction]] |
| 0033 | `Decide` returns state errors | Accepted | [[V4 - Algorithm Abstraction]] |
| 0034 | Initialize missing state before asserting | Accepted | [[V4 - Algorithm Abstraction]] |
| 0035 | Token bucket as the second algorithm | Accepted | [[V5 - Token Bucket]] |
| 0036 | Capacity and refill rate | Accepted | [[V5 - Token Bucket]] |
| 0037 | Token bucket state shape | Accepted | [[V5 - Token Bucket]] |
| 0038 | Lazy token refill | Accepted | [[V5 - Token Bucket]] |
| 0039 | Floating-point token arithmetic | Accepted | [[V5 - Token Bucket]] |
| 0040 | Whole-request `Remaining` | Accepted | [[V5 - Token Bucket]] |
| 0041 | Token bucket `RetryAfter` | Accepted | [[V5 - Token Bucket]] |
| 0042 | Clamp negative elapsed time | Accepted | [[V5 - Token Bucket]] |
| 0043 | Token bucket starts full | Accepted | [[V5 - Token Bucket]] |
| 0044 | Split files by responsibility | Accepted | [[V5 - Token Bucket]] |
| 0045 | Split storage from the limiter | Accepted | [[V6 - Storage Boundary]] |
| 0046 | The Get/Set race | Accepted | [[V6 - Storage Boundary]] |
| 0047 | StateStore uses Get + CAS | Accepted | [[V6 - Storage Boundary]] |
| 0048 | Benchmark before the storage refactor | Accepted | [[V6 - Storage Boundary]] |
| 0049 | MemoryStore owns runtime state | Accepted | [[V6 - Storage Boundary]] |
| 0050 | CAS conflict is not a store error | Accepted | [[V6 - Storage Boundary]] |
| 0051 | Bounded CAS retry loop | Accepted | [[V6 - Storage Boundary]] |
| 0052 | Default memory store in the constructor | Superseded by 0059 | [[V6 - Storage Boundary]] |
| 0053 | Shard the memory store next | Accepted | [[V7 - Sharded MemoryStore]] |
| 0054 | Sharded store keeps the StateStore contract | Accepted | [[V7 - Sharded MemoryStore]] |
| 0055 | Configurable shard count | Accepted | [[V7 - Sharded MemoryStore]] |
| 0056 | Reuse MemoryStore internally | Accepted | [[V7 - Sharded MemoryStore]] |
| 0057 | Sharded store as the default backend | Superseded by 0059 | [[V7 - Sharded MemoryStore]] |
| 0058 | Trust internal invariants | Accepted | [[V7 - Sharded MemoryStore]] |
| 0059 | Store injection opened | Accepted | [[V7 - Sharded MemoryStore]] |
| 0060 | Contract testing for StateStore | Accepted | [[V7 - Sharded MemoryStore]] |
| 0061 | Test isolation with fresh stores | Accepted | [[V7 - Sharded MemoryStore]] |
| 0062 | Redis as a limiter, not a store (pivot) | Accepted | [[V7 - Sharded MemoryStore]] |
| 0063 | Fixed window in Redis via Lua | Accepted | [[V8 - RedisLimiter]] |
| 0064 | RedisLimiter algorithm interface | Accepted | [[V8 - RedisLimiter]] |
| 0065 | Redis adapter returns raw result | Accepted | [[V8 - RedisLimiter]] |
| 0066 | Redis client ownership boundary | Accepted | [[V8 - RedisLimiter]] |
| 0067 | Sealed Redis algorithms | Accepted | [[V8 - RedisLimiter]] |
| 0068 | Shared fixed-window config | Accepted | [[V8 - RedisLimiter]] |
| 0069 | Backend-qualified naming | Accepted | [[V8 - RedisLimiter]] |
| 0070 | Redis token bucket: scaled-integer state | Accepted | [[V8 - RedisLimiter]] |
| 0071 | Redis token bucket: time and TTL policy | Accepted | [[V8 - RedisLimiter]] |
| 0072 | Redis token bucket: result contract | Accepted | [[V8 - RedisLimiter]] |
| 0073 | Redis token bucket: scaling boundary | Accepted | [[V8 - RedisLimiter]] |
| 0074 | HTTP middleware boundary | Accepted | [[V9 - HTTP Middleware]] |
| 0075 | Middleware failure policy | Accepted | [[V9 - HTTP Middleware]] |
| 0076 | Caller-provided key function | Accepted | [[V9 - HTTP Middleware]] |
| 0077 | Rate-limit HTTP headers | Accepted | [[V9 - HTTP Middleware]] |
| 0078 | Functional options for the middleware | Accepted | [[V9 - HTTP Middleware]] |
| 0079 | Behavior-named middleware | Accepted | [[V9 - HTTP Middleware]] |
| 0080 | Required dependencies outside functional options | Accepted | [[V9 - HTTP Middleware]] |
| 0081 | Fail-open is silent pass-through | Accepted | [[V9 - HTTP Middleware]] |
| 0082 | HTTP Retry-After converts Duration to seconds | Accepted | [[V9 - HTTP Middleware]] |
| 0083 | MemoryLimiter implements Limiter via a private clock | Accepted (supersedes 0005) | [[V9 - HTTP Middleware]] |
| 0084 | Benchmarks use the production clock | Accepted | [[V9 - HTTP Middleware]] |
| 0085 | Hot-key CAS exhaustion is expected | Accepted | [[V9 - HTTP Middleware]] |

## Adding a new ADR

Append a new `## ADR-00NN — Title` section (Status / Context / Decision / Consequences) to the **current** version's log, add a row here, and update that log's narrative intro if the decision changes the iteration's story. To reverse a past decision, add a new ADR and set the old one's status to `Superseded by 00NN`. Start a new log only when a new version (V10+) begins.

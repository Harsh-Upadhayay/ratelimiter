# Rate Limiter Learning Map

This is the Obsidian entry point for the learning version of the rate limiter.

## Roadmap

- [[plan|Target 6-step plan]]
- [[mindmaps/V1 Learning Mindmap|V1 learning mindmap]]
- [[V1 Decision Index]]
- [[V2 API Evolution Index]]
- [[Go Concepts Index]]

## V1 decisions

- Hub: [[V1 Decision Index]]
- [[decisions/D01 - Fixed Window First]]
- [[decisions/D02 - State Shape Map String UserState]]
- [[decisions/D03 - Pointer Receiver for Mutating Limiter]]
- [[decisions/D04 - One Global Mutex]]
- [[decisions/D05 - Explicit Time Input]]
- [[decisions/D06 - Count Invariant Accepted Requests]]
- [[decisions/D07 - Window Reset Starts At Now]]
- [[decisions/D08 - Half Open Window Interval]]
- [[decisions/D09 - Defer Unlock After Lock]]
- [[decisions/D11 - Constructor Validation]]
- [[decisions/D12 - Exported API Boundary]]
- [[decisions/D13 - Delay Interfaces]]
- [[decisions/D14 - Package Name]]
- [[decisions/D15 - One File First]]
- [[decisions/D16 - Limiter Owns Runtime State]]
- [[decisions/D17 - Constructor Returns Pointer]]
- [[decisions/D18 - Empty Key Validation]]

## V2 decisions

- Hub: [[V2 API Evolution Index]]
- [[decisions/D19 - Result and Error Return]]
- [[decisions/D20 - Generic Rate Limit Key]]
- [[decisions/D21 - Result Contract]]
- [[decisions/D22 - Sentinel Error for Empty Key]]

## V3 decisions

- Hub: [[V3 Decision Logic Index]]
- [[decisions/D23 - Private Fixed Window Decision Helper]]
- [[decisions/D24 - Test Through Public API]]
- [[decisions/D25 - Allow as Orchestrator]]

## V4 decisions

- Hub: [[V4 Algorithm Abstraction Index]]
- [[decisions/D26 - Introduce Algorithm Interface]]
- [[decisions/D27 - Algorithm Owned State]]
- [[decisions/D28 - Marker Interface for Algorithm State]]
- [[decisions/D29 - Exists Flag for Missing State]]
- [[decisions/D30 - Keep Algorithm Interface Private]]
- [[decisions/D31 - Algorithm Owns Config Validation]]
- [[decisions/D32 - Manual Assembly with Built In Algorithms]]
- [[decisions/D33 - Decide Returns State Errors]]
- [[decisions/D34 - Initialize Missing State Before Assertion]]
- Checkpoint: [[checkpoints/C04 - V4 Algorithm Boundary Checkpoint]]

## V5 decisions

- Hub: [[V5 Token Bucket Design Index]]
- [[decisions/D35 - Token Bucket as Second Algorithm]]
- [[decisions/D36 - Token Bucket Capacity and Refill Rate]]
- [[decisions/D37 - Token Bucket State Shape]]
- [[decisions/D38 - Lazy Token Refill]]
- [[decisions/D39 - Floating Point Token Arithmetic]]
- [[decisions/D40 - Whole Request Remaining]]
- [[decisions/D41 - Token Bucket RetryAfter]]
- [[decisions/D42 - Clamp Negative Elapsed Time]]
- [[decisions/D43 - Token Bucket Starts Full]]
- [[decisions/D44 - Split Files by Responsibility]]

## V6 decisions

- Hub: [[V6 Storage Boundary Index]]
- [[decisions/D45 - Split Storage from Limiter]]
- [[decisions/D46 - GetSet Race]]
- [[decisions/D47 - StateStore Uses Get and CAS]]
- [[decisions/D48 - Benchmark Before Storage Refactor]]
- [[decisions/D49 - MemoryStore Owns Runtime State]]
- [[decisions/D50 - CAS Conflict Is Not Store Error]]
- [[decisions/D51 - Bounded CAS Retry Loop]]
- [[decisions/D52 - Default Memory Store Constructor]]

## V7 decisions

- Hub: [[V7 Sharded MemoryStore Index]]
- [[decisions/D53 - Sharded MemoryStore Next]]
- [[decisions/D54 - Sharded Store Keeps StateStore Contract]]
- [[decisions/D55 - Configurable Shard Count]]
- [[decisions/D56 - Reuse MemoryStore Internally]]
- [[decisions/D57 - Sharded MemoryStore as Default Backend]]
- [[decisions/D58 - Trust Internal Invariants]]
- [[decisions/D59 - Store Injection Opened]]
- [[decisions/D60 - Contract Testing for StateStore]]
- [[decisions/D61 - Test Isolation with Fresh Stores]]
- [[decisions/D62 - Redis as Limiter Not Store]] — **pivot**: revises plan step 4; Redis is a parallel `RedisLimiter`, not a `StateStore`

## V8 — RedisLimiter (breadth over the distributed axis)

V7 closed the single-node depth axis. Next learning is breadth: the distributed path and the
sliding-window algorithm family. Trailhead is [[decisions/D62 - Redis as Limiter Not Store]].
Direction chosen 2026-06-14: **Redis-first**, starting with Fixed Window in Lua.

- Hub: [[V8 RedisLimiter Design Index]]
- [[decisions/D63 - Fixed Window in Redis via Lua]]
- [[decisions/D64 - RedisLimiter Algorithm Interface]]
- [[decisions/D65 - Redis Adapter Returns Raw Result]]
- [[decisions/D66 - Redis Client Ownership Boundary]]
- [[decisions/D67 - Sealed Redis Algorithms]]
- [[decisions/D68 - Shared Fixed Window Config]]
- [[decisions/D69 - Backend Qualified Naming]]
- Redis/Lua concepts hub: [[Redis Concepts Index]] (R01–R07)

## Go practice notes

- Hub: [[Go Concepts Index]]
- [[go/G01 - Structs for Grouped State]]
- [[go/G02 - Maps Comma Ok and Value Copies]]
- [[go/G03 - Pointer Receivers]]
- [[go/G04 - Mutexes and Critical Sections]]
- [[go/G05 - Defer Unlock Pattern]]
- [[go/G06 - Time and Duration Boundaries]]
- [[go/G07 - Do Not Copy Mutexes]]
- [[go/G08 - Exported Types with Unexported Fields]]
- [[go/G09 - Method Receivers]]
- [[go/G10 - Early Returns and Guard Clauses]]
- [[go/G11 - Multiple Return Values]]
- [[go/G12 - Sentinel Errors]]
- [[go/G13 - Exported Result Structs]]
- [[go/G14 - Pure Helper Functions]]
- [[go/G15 - Public API Tests]]
- [[go/G16 - Helper Parameter Ordering]]
- [[go/G17 - Interfaces From Real Variation]]
- [[go/G18 - Structural Interface Satisfaction]]
- [[go/G19 - Marker Interfaces and Opaque State]]
- [[go/G20 - Type Assertions]]
- [[go/G21 - Constructor Validation Ownership]]
- [[go/G22 - Nil Interface Guard]]
- [[go/G23 - Duration Seconds Conversion]]
- [[go/G24 - Lazy State Materialization]]
- [[go/G25 - Package Scope Across Files]]
- [[go/G26 - Optimistic Concurrency with CAS]]
- [[go/G27 - Store Owned Mutexes]]
- [[go/G28 - Go Benchmarks]]
- [[go/G29 - Race Detector]]
- [[go/G30 - Exported Interfaces With Unexported Types]]
- [[go/G31 - Lock Striping]]
- [[go/G32 - Key Hashing]]
- [[go/G33 - Composition with Pointer Fields]]
- [[go/G34 - String to Byte Slice]]
- [[go/G35 - Unsigned Modulo Safety]]
- [[go/G36 - Test Helper Pattern]]
- [[go/G37 - Contract Testing via Interface Parameter]]
- [[go/G38 - context.Context]]
- [[go/G39 - Go Naming Qualifier Order]]

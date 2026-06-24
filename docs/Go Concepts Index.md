# Go Concepts Index

Back to [[Rate Limiter Learning Map]].

This hub groups Go language and style concepts used by the rate limiter.

## State and data modeling

- [[G01 - Structs for Grouped State]]
- [[G02 - Maps Comma Ok and Value Copies]]
- [[G08 - Exported Types with Unexported Fields]]
- [[G13 - Exported Result Structs]]

## Methods and mutation

- [[G03 - Pointer Receivers]]
- [[G09 - Method Receivers]]
- [[G46 - Method Sets and Interface Satisfaction]]
- [[G07 - Do Not Copy Mutexes]]

## Concurrency

- [[G04 - Mutexes and Critical Sections]]
- [[G05 - Defer Unlock Pattern]]
- [[G26 - Optimistic Concurrency with CAS]]
- [[G27 - Store Owned Mutexes]]
- [[G29 - Race Detector]]
- [[G31 - Lock Striping]]
- [[G32 - Key Hashing]]
- [[G33 - Composition with Pointer Fields]]
- [[G34 - String to Byte Slice]]
- [[G35 - Unsigned Modulo Safety]]

## Time and control flow

- [[G06 - Time and Duration Boundaries]]
- [[G10 - Early Returns and Guard Clauses]]
- [[G23 - Duration Seconds Conversion]]
- [[G48 - Ceiling Duration Conversion]]
- [[G24 - Lazy State Materialization]]
- [[G25 - Package Scope Across Files]]

## Errors and API shape

- [[G11 - Multiple Return Values]]
- [[G12 - Sentinel Errors]]
- [[G39 - Go Naming Qualifier Order]]
- [[G40 - Unexported Package Constants]]
- [[G42 - Functional Options]]
- [[G47 - Functional Options on Runtime Structs]]
- [[G43 - Iota Enum Pattern]]
- [[G44 - Function Types as Callbacks]]
- [[G14 - Pure Helper Functions]]
- [[G15 - Public API Tests]]
- [[G16 - Helper Parameter Ordering]]
- [[G17 - Interfaces From Real Variation]]
- [[G18 - Structural Interface Satisfaction]]
- [[G19 - Marker Interfaces and Opaque State]]
- [[G20 - Type Assertions]]
- [[G21 - Constructor Validation Ownership]]
- [[G22 - Nil Interface Guard]]
- [[G30 - Exported Interfaces With Unexported Types]]

## HTTP integration

- [[G41 - net-http Middleware Pattern]]
- [[G49 - http HandlerFunc Adapter]]
- [[G45 - HTTP ResponseWriter Headers]]

## Testing and benchmarking

- [[G28 - Go Benchmarks]]
- [[G36 - Test Helper Pattern]]
- [[G37 - Contract Testing via Interface Parameter]]

## Decision hubs

- [[V1 Decision Index]]
- [[V2 API Evolution Index]]
- [[V3 Decision Logic Index]]
- [[V4 Algorithm Abstraction Index]]
- [[V5 Token Bucket Design Index]]
- [[V6 Storage Boundary Index]]
- [[V7 Sharded MemoryStore Index]]
- [[V9 HTTP Middleware Index]]

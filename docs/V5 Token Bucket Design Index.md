# V5 Token Bucket Design Index

Back to [[Rate Limiter Learning Map]].

This hub tracks the token bucket design before implementation.

## Core decisions

- [[D35 - Token Bucket as Second Algorithm]]
- [[D36 - Token Bucket Capacity and Refill Rate]]
- [[D37 - Token Bucket State Shape]]
- [[D38 - Lazy Token Refill]]
- [[D39 - Floating Point Token Arithmetic]]
- [[D40 - Whole Request Remaining]]
- [[D41 - Token Bucket RetryAfter]]
- [[D42 - Clamp Negative Elapsed Time]]
- [[D43 - Token Bucket Starts Full]]

## Related earlier decisions

- [[D26 - Introduce Algorithm Interface]]
- [[D27 - Algorithm Owned State]]
- [[D31 - Algorithm Owns Config Validation]]
- [[D32 - Manual Assembly with Built In Algorithms]]

## Go concepts used

- [[G06 - Time and Duration Boundaries]]
- [[G21 - Constructor Validation Ownership]]
- [[G23 - Duration Seconds Conversion]]
- [[G24 - Lazy State Materialization]]

## Future bridge

- [[plan|Target 6-step plan]]

# D36 - Token Bucket Capacity and Refill Rate

Back to [[Rate Limiter Learning Map]].

## Context

Token bucket separates burst capacity from sustained refill rate.

## Decision

Configure token bucket with:

- `capacity int`
- `refillRate float64`

`refillRate` means tokens per second.

## Why

Capacity represents the maximum number of complete requests that can burst immediately. Refill rate represents sustained throughput over time.

## Tradeoff

This differs from fixed window's limit/window configuration, which confirms that algorithm-specific configuration should live with the algorithm constructor.

## Links

- [[D31 - Algorithm Owns Config Validation]]
- [[D43 - Token Bucket Starts Full]]
- [[G23 - Duration Seconds Conversion]]

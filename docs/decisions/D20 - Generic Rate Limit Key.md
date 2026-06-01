# D20 - Generic Rate Limit Key

Back to [[Rate Limiter Learning Map]].

## Context

The first implementation used `userID` language. A real rate limiter may key by user, IP address, API token, tenant, route, or a composite value.

## Decision

Rename the public input concept from user ID to rate-limit key.

## Why

`key` describes the generic dimension being limited without tying the API to users.

## Tradeoff

`key` is more abstract than `userID`, so examples must make the intended key clear.

## Future direction

HTTP middleware can build keys from request attributes before calling the limiter.

## Links

- [[D12 - Exported API Boundary]]
- [[D22 - Sentinel Error for Empty Key]]

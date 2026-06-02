# D41 - Token Bucket RetryAfter

Back to [[Rate Limiter Learning Map]].

## Context

For token bucket, rejected requests need to know when enough tokens will exist for one request.

## Decision

For a rejected request:

```text
missingTokens = 1 - availableTokens
RetryAfter = missingTokens / refillRate
```

## Why

This preserves the public meaning of `RetryAfter`: minimum wait before retrying.

## Tradeoff

The semantics differ internally from fixed window. Fixed window waits until window reset; token bucket waits until one token is available.

## Links

- [[D21 - Result Contract]]
- [[D39 - Floating Point Token Arithmetic]]

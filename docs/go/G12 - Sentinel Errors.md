# G12 - Sentinel Errors

Back to [[Rate Limiter Learning Map]].

## Go practice

A sentinel error is a package-level error value used for a stable error category.

Callers can compare with `errors.Is`.

## Why it matters

Sentinel errors are useful when callers need to branch on a known public error.

## Tradeoff

Exported sentinel errors are public API. Avoid exporting errors for temporary or overly specific internal failures.

## Rate limiter example

An empty rate-limit key is a stable validation category, so an exported sentinel is reasonable.

## Links

- [[D22 - Sentinel Error for Empty Key]]
- [[G11 - Multiple Return Values]]

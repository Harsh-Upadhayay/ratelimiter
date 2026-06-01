# G13 - Exported Result Structs

Back to [[Rate Limiter Learning Map]].

## Go practice

If a public function returns a struct that callers must inspect, both the type and the inspected fields should be exported.

## Why it matters

Returning an exported type with unexported fields prevents callers outside the package from reading the result details.

## Rate limiter example

`Result` should be exported because `Allow` returns it.

Its fields should also be exported so callers can read `Allowed`, `Remaining`, and `RetryAfter`.

## Links

- [[D21 - Result Contract]]
- [[G08 - Exported Types with Unexported Fields]]

# G11 - Multiple Return Values

Back to [[Rate Limiter Learning Map]].

## Go practice

Go functions commonly return multiple values.

For APIs that can produce a useful result and also fail, the common shape is:

```go
result, err := operation()
```

## Why it matters

The caller checks the error first, then uses the result.

## Rate limiter meaning

For `Allow`, a rate-limit rejection is not an error. It is a successful decision with `Allowed` set to false.

Invalid input is an error.

## Links

- [[D19 - Result and Error Return]]
- [[D21 - Result Contract]]

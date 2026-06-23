# G45 - HTTP ResponseWriter Headers

Back to [[Go Concepts Index]].

## Concept

In Go's `net/http`, response headers are set through:

```go
w.Header().Set("Name", "value")
```

Status code is written with:

```go
w.WriteHeader(statusCode)
```

Headers should be set before writing the status code or response body.

## In this project

The middleware should set rate-limit headers before returning a rejection:

```text
X-RateLimit-Remaining
Retry-After
```

It should not write a response body for v1.

## Common mistake

Once the response body is written, Go implicitly sends status `200 OK` if no status was written.

For middleware that rejects a request, write the intended status explicitly before returning.

## Links

- [[D77 - Rate Limit HTTP Headers]]
- [[G41 - net-http Middleware Pattern]]

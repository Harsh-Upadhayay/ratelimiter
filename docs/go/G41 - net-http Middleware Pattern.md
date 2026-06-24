# G41 - net-http Middleware Pattern

Back to [[Go Concepts Index]].

## Concept

Go's standard HTTP server is built around:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Middleware is usually a function or method that takes one handler and returns another handler.

Conceptually:

```go
func Middleware(next http.Handler) http.Handler
```

The wrapper can inspect the request, decide whether to stop, or call:

```go
next.ServeHTTP(w, r)
```

## HandlerFunc adapter

For simple handlers, Go provides `http.HandlerFunc`, which turns a plain function with this shape:

```go
func(http.ResponseWriter, *http.Request)
```

into an `http.Handler`. See [[G49 - http HandlerFunc Adapter]].

## In this project

The rate-limiting middleware sits before the application handler:

```text
request -> rate-limiting middleware -> app handler
```

If allowed, it calls the next handler.

If rejected, it writes status and headers and does not call the next handler.

## Practice

Use middleware when the behavior is cross-cutting:

- auth,
- logging,
- tracing,
- rate limiting,
- compression.

## Links

- [[D74 - HTTP Middleware Boundary]]
- [[D79 - Behavior Named Middleware]]
- [[G49 - http HandlerFunc Adapter]]

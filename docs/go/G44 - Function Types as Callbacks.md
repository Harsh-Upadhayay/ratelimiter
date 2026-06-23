# G44 - Function Types as Callbacks

Back to [[Go Concepts Index]].

## Concept

Functions are values in Go.

That means a package can define a function type:

```go
type KeyFunc func(*http.Request) string
```

and store it in a struct field, pass it as an argument, or call it later.

## In this project

The middleware does not know the application identity model.

Instead, the caller provides a key extraction function:

```text
request -> KeyFunc -> rate-limit key
```

## Why

This avoids hard-coding assumptions such as:

- remote IP,
- `X-User-ID`,
- `Authorization`,
- tenant ID.

The application owns identity. The middleware owns enforcement behavior.

## Links

- [[D76 - Caller Provided Key Function]]
- [[D74 - HTTP Middleware Boundary]]

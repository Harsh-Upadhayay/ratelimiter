# G49 - http HandlerFunc Adapter

Back to [[Go Concepts Index]].

## Concept

`http.Handler` is an interface:

```go
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```

`http.HandlerFunc` is an adapter type that lets a plain function satisfy that interface.

Conceptually, the standard library gives it a method like:

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

## In this project

`Wrap` returns a handler by converting a function:

```go
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // per-request middleware logic
})
```

`Wrap` runs during setup. The returned function runs once per HTTP request.

## Mental model

```text
Wrap(next)            -> builds protected handler
protected.ServeHTTP   -> runs on each request
next.ServeHTTP        -> continues to the protected endpoint
w.WriteHeader(status) -> stops in middleware with that status
```

## Links

- [[G41 - net-http Middleware Pattern]]
- [[D74 - HTTP Middleware Boundary]]

# D76 - Caller Provided Key Function

Back to [[V9 HTTP Middleware Index]].

## Context

HTTP middleware must turn a request into a rate-limit key.

Possible keys include:

- remote IP,
- user ID,
- API key,
- tenant ID,
- route plus user ID.

The middleware cannot know the application's identity model.

## Decision

Use a caller-provided key function.

Conceptually:

```go
type KeyFunc func(*http.Request) string
```

The middleware calls the function for every request and passes the returned key to the limiter.

If the key function returns an empty key, the middleware returns:

```text
400 Bad Request
```

## Why

Key choice is application-owned. The middleware should not pretend that IP address, headers, or auth
claims are universally correct.

## Tradeoff

- **Correctness:** better, because the application owns identity.
- **Flexibility:** high; one middleware can support many key schemes.
- **Complexity:** caller must provide a function.
- **Security:** the caller must ensure the source of the key is trustworthy.

## Links

- [[D20 - Generic Rate Limit Key]]
- [[D18 - Empty Key Validation]]
- [[G44 - Function Types as Callbacks]]

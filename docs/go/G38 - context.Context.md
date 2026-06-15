# G38 - context.Context

Back to [[Rate Limiter Learning Map]].

## Concept

`context.Context` is a value threaded down a call stack to carry **cancellation signals**,
**deadlines/timeouts**, and (sparingly) request-scoped values. It lets the top of a call chain
tell the bottom "stop waiting, we're done" without every function inventing its own signalling.

The motivating case: a blocking I/O call (a network round-trip to Redis) can hang forever. The
caller — an HTTP handler, a batch job, a test — is the only layer that knows how long the whole
operation may take. `ctx` carries that budget down to the I/O call so it can abort instead of
hanging.

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()                       // always release the timer
result := limiter.Allow(ctx, key)    // ctx flows down into the Redis Eval
```

If Redis does not answer within 2s, the `Eval` inside returns a deadline-exceeded error rather
than blocking.

## Roots

- `context.Background()` — the empty root: never cancels, no deadline. Use at the true top
  (`main`, tests).
- `context.TODO()` — same behavior, but signals to readers "context story not figured out yet."

Everything else is *derived* from a root via `context.WithTimeout` / `context.WithCancel` / etc.

## Two ironclad conventions

1. **`ctx` is the first parameter**, named `ctx`:
   `func (rl *RedisLimiter) Allow(ctx context.Context, key string) Result`.
2. **Never store a `Context` in a struct field.** A context models the lifetime of *one
   operation*, not the lifetime of an *object*. Stashing it in `RedisLimiter` would force every
   call to share one stale context — call #2 could not get a different deadline than call #1, and
   a cancellation meant for one request would poison all of them. Pass it, don't store it.

## Rate limiter use

`RedisLimiter.Allow(ctx, key)` takes the caller's context and passes it straight into go-redis's
`Eval(ctx, script, keys, args...)`. `RedisLimiter` does **not** invent its own
`context.Background()` — that would silently impose "wait forever" on every caller and rob them of
the cancellation decision they alone can make correctly ([[D62 - Redis as Limiter Not Store]]).

This diverges from the in-process `Limiter.Allow(key, now)`, which does no I/O and takes no `ctx` —
a real gap to reconcile if/when a unified `RateLimiter` interface is revisited (see
[[D64 - RedisLimiter Algorithm Interface]]).

## Links

- [[D64 - RedisLimiter Algorithm Interface]]
- [[V8 RedisLimiter Design Index]]

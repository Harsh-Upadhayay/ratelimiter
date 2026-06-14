# R03 - Lua Scripting for Atomicity

Back to [[Redis Concepts Index]].

## Concept

A Lua script sent to Redis via `EVAL` executes atomically — all Redis commands in the script run before any other client's request can proceed.

## Why this matters

Rate limiting logic (check limit, increment, set expiry) must be atomic across processes. Lua scripts provide that guarantee at the Redis server level, without client-side CAS retry loops.

## Shape

```lua
local variable = redis.call('COMMAND', KEYS[1], ARGV[1])
if condition then
    redis.call('ANOTHER_COMMAND', KEYS[1], ARGV[2])
end
return variable
```

- Scripts can store intermediate results in local variables.
- Scripts can use Lua control flow (`if`, loops, etc.).
- Multiple `redis.call()` invocations all happen atomically.
- The final `return` sends a value back to the caller.

## In Go

From Go, you send the script text and its arguments to Redis via the client library (e.g. `go-redis`). Redis executes it atomically and returns the result.

## Related

- [[R04 - redis.call() and Command Execution]]
- [[R05 - KEYS and ARGV Parameters]]
- [[D62 - Redis as Limiter Not Store]]

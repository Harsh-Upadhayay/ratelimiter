# R05 - KEYS and ARGV Parameters

Back to [[Redis Concepts Index]].

## Concept

A Lua script receives data from the caller via two arrays:
- `KEYS[]` — keys the script will access.
- `ARGV[]` — arbitrary arguments (often configuration).

## From redis-cli

```
EVAL "script" numkeys key1 key2 arg1 arg2
       ^^^^^^^          ^^^^^^            ^^^^
       script text      KEYS[]           ARGV[]
```

Inside the script:
- `KEYS[1]` is `key1`, `KEYS[2]` is `key2`.
- `ARGV[1]` is `arg1`, `ARGV[2]` is `arg2`.

## From Go

When calling a Lua script from Go (via `go-redis`), you pass keys and args separately, and the library handles the `numkeys` parameter.

## Convention

Pass only *actual keys* in `KEYS[]` (keys the script modifies/reads). Pass configuration (TTL, limits, thresholds) in `ARGV[]`. This separation helps Redis optimize key tracking for cluster mode.

## Example

For fixed-window rate limiting:

```lua
-- Script called with KEYS[1] = "limit:user123", ARGV[1] = "10" (window seconds)
local count = redis.call('INCR', KEYS[1])
if count == 1 then
    redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return count
```

## Related

- [[R04 - redis.call() and Command Execution]]
- [[R03 - Lua Scripting for Atomicity]]

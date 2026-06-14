# R04 - redis.call() and Command Execution

Back to [[Redis Concepts Index]].

## Concept

Within a Lua script, `redis.call('COMMAND', arg1, arg2, ...)` executes a Redis command synchronously and returns its result.

## Syntax

```lua
local result = redis.call('INCR', mykey)
local exists = redis.call('EXISTS', mykey)
redis.call('EXPIRE', mykey, 60)
```

- The first argument is the command name as a string (`'INCR'`, `'GET'`, etc.).
- Remaining arguments are command arguments.
- `redis.call()` blocks until the command completes and returns the result.

## Return values

Results vary by command:
- `INCR` returns the new integer value.
- `GET` returns the string value or nil.
- `EXPIRE` returns 1 if the key existed and TTL was set, 0 otherwise.

## Error handling

If a command fails, `redis.call()` raises a Lua error. You can use `redis.pcall()` (protected call) to catch errors without aborting the script.

## Related

- [[R03 - Lua Scripting for Atomicity]]
- [[R05 - KEYS and ARGV Parameters]]

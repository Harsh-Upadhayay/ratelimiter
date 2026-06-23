# R10 - Lua Scripts Embedded in Go Strings

Back to [[Redis Concepts Index]].

## Concept

In this project, Redis Lua scripts are stored as Go raw strings:

```go
func (tb *RedisTokenBucket) script() string {
    return `
        -- Lua code here
    `
}
```

The Go compiler only checks that this is a valid Go string. It does not parse Lua.

## Consequence

These mistakes can compile in Go but fail only when Redis executes the script:

- misspelled Lua variables,
- missing `local`,
- wrong `ARGV` index,
- wrong Redis command shape,
- wrong return shape.

Example:

```lua
local refillRate = tonumber(ARGV[2])
local refilled = math.floor((elapsedMs * refilRate) / 1000)
```

`refilRate` and `refillRate` are different Lua names. Go will not catch that because both are inside a
string.

## Practice

Keep Lua variable names explicit and unit-bearing:

```lua
local refillUnitsPerSecond = tonumber(ARGV[2])
local elapsedMs = nowMs - lastRefillMs
```

For Redis-backed algorithms, a passing `go test ./...` only proves the Go package compiles. It does
not prove the Lua script runs correctly unless a test actually executes the script against Redis.

## Related

- [[R03 - Lua Scripting for Atomicity]]
- [[R04 - redis.call() and Command Execution]]
- [[D73 - Redis Token Bucket Scaling Boundary]]

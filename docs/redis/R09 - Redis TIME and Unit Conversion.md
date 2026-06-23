# R09 - Redis TIME and Unit Conversion

Back to [[Redis Concepts Index]].

## Concept

`redis.call("TIME")` returns Redis server time as:

```text
seconds, microseconds
```

In Lua:

```lua
local now = redis.call("TIME")
```

Then:

```text
now[1] = seconds
now[2] = microseconds
```

## Millisecond Conversion

Redis Token Bucket uses milliseconds internally:

```lua
local nowMs = tonumber(now[1]) * 1000 + math.floor(tonumber(now[2]) / 1000)
```

The refill rate is still expressed as scaled token units per second.

So refill math converts elapsed milliseconds back into seconds:

```lua
local refilled = math.floor((elapsedMs * refillUnitsPerSecond) / 1000)
```

The `/ 1000` exists because:

```text
milliseconds / 1000 = seconds
```

## Why Not Microseconds Now

Microseconds remove the `microseconds / 1000` conversion, but the main conversion remains:

```text
elapsed time unit -> seconds
```

Microseconds also create larger numbers. Milliseconds are precise enough for the current limiter and
keep the arithmetic easier to inspect.

## Related

- [[R03 - Lua Scripting for Atomicity]]
- [[D71 - Redis Token Bucket Time and TTL Policy]]
- [[D70 - Redis Token Bucket Scaled Integer State]]

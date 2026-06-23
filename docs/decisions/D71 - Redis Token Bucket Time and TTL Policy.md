# D71 - Redis Token Bucket Time and TTL Policy

Back to [[V8 RedisLimiter Design Index]].

## Context

Redis Token Bucket needs a trustworthy clock and a cleanup policy for idle keys.

The in-memory limiter accepts `now` from the caller for deterministic tests. The Redis limiter runs
inside Redis, where multiple application instances may call the same limiter key.

## Options

- **Caller-provided time** - easy to test, but distributed callers may disagree because of clock
  skew.
- **Application server time inside Go** - simple, but each app instance can still have a different
  clock.
- **Redis `TIME` inside Lua** - one authoritative clock for the atomic Redis decision.

## Decision

Use Redis `TIME` inside the Lua script and convert it to milliseconds:

```lua
local now = redis.call("TIME")
local nowMs = tonumber(now[1]) * 1000 + math.floor(tonumber(now[2]) / 1000)
```

Use TTL as idle cleanup, not as the refill mechanism.

```text
ttlSeconds = ceil(capacity / refillRate)
```

Since Redis Token Bucket passes scaled capacity and scaled refill rate to Lua, the scale cancels:

```lua
local ttlSeconds = math.ceil(capacity / refillRate)
```

Refresh the TTL on every script call, including rejected requests.

## Why

After a full-refill duration of no traffic, even an empty bucket would be full again. Deleting the
Redis key after that much idle time does not change behavior.

Refreshing TTL on every request prevents active keys from expiring while the caller is still sending
traffic. Without the refresh, a constantly rejected caller could let the key expire and then recreate
a full bucket, accidentally receiving a new burst.

## Tradeoffs

- **Memory:** exact full-refill TTL minimizes idle key retention. A buffered TTL would retain more
  keys for observability and lower churn.
- **Latency:** one `EXPIRE` per script call adds Redis work, but it happens inside the same Lua
  execution and avoids another network round trip.
- **Concurrency:** Redis `TIME` plus Lua atomicity makes all callers share the same clock and state
  transition.
- **Precision:** milliseconds are precise enough for this limiter and keep numbers smaller than
  microseconds.

## Links

- [[D70 - Redis Token Bucket Scaled Integer State]]
- [[redis/R01 - Key-Value Store and TTL]]
- [[redis/R03 - Lua Scripting for Atomicity]]
- [[redis/R09 - Redis TIME and Unit Conversion]]

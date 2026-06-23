# R08 - Redis Hashes and HMGET

Back to [[Redis Concepts Index]].

## Concept

Redis has top-level keys, and each key stores a value of some Redis data type.

The top-level key is always located internally by Redis, but the value stored at that key can be a
string, hash, list, set, sorted set, and so on.

## String value

```text
rate:user:123 -> "750"
```

Use:

```lua
redis.call("GET", key)
```

## Hash value

```text
rate:user:123 -> {
  tokens: "750",
  last_refill_ms: "1719000000000"
}
```

Use:

```lua
redis.call("HMGET", key, "tokens", "last_refill_ms")
```

`HMGET` means: get multiple fields from the hash stored at this key.

## Why Token Bucket Uses a Hash

Redis Token Bucket needs two state fields:

- `tokens`
- `last_refill_ms`

A single `GET` value would require packing both into one string and parsing it manually. A Redis hash
keeps the state structured while still using one top-level Redis key per rate-limit key.

## Lua Shape

`HMGET` returns a table. Missing fields appear as missing/false-like values inside the table.

Check the fields, not the table:

```lua
local state = redis.call("HMGET", key, "tokens", "last_refill_ms")
local tokens = state[1]
local lastRefillMs = state[2]

if not tokens or not lastRefillMs then
  -- cold start
end
```

## Related

- [[R04 - redis.call() and Command Execution]]
- [[R05 - KEYS and ARGV Parameters]]
- [[D70 - Redis Token Bucket Scaled Integer State]]

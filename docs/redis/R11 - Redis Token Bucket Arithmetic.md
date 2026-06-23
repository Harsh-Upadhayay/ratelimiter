# R11 - Redis Token Bucket Arithmetic

Back to [[Redis Concepts Index]].

## Why This Exists

Redis Token Bucket has more arithmetic than Fixed Window because refill is continuous.

Fixed Window can count whole requests:

```text
remaining = remaining - 1
```

Token Bucket must preserve partial refill progress:

```text
elapsed = 250ms
refillRate = 2 tokens/sec
refill = 0.5 token
```

That partial `0.5` matters. If it is discarded, frequent small refills never accumulate correctly.

## Choice

Redis Token Bucket uses scaled integer arithmetic.

```text
tokenScale = 1000

1 token   = 1000 units
0.5 token = 500 units
0.001 token = 1 unit
```

The public configuration still uses domain units:

```text
capacity = tokens
refillRate = tokens per second
```

The Redis algorithm converts those into internal units:

```text
capacityScaled = capacity * tokenScale
refillUnitsPerSecond = refillRate * tokenScale
```

## Current Internal Units

Inside Redis Token Bucket:

```text
tokens = scaled token units
time = milliseconds
refillRate = scaled token units per second
request cost = tokenScale
```

So this check:

```lua
tokens >= tokenScale
```

means:

```text
Does this key have at least one full request token?
```

## Refill Formula

Redis `TIME` gives seconds and microseconds. The script converts that to milliseconds:

```lua
local nowMs = seconds * 1000 + math.floor(microseconds / 1000)
```

Elapsed time is:

```lua
local elapsedMs = nowMs - lastRefillMs
```

Refill is:

```lua
local refilled = math.floor((elapsedMs * refillUnitsPerSecond) / 1000)
```

The `/ 1000` converts milliseconds back into seconds:

```text
elapsedMs / 1000 = elapsed seconds
```

Example:

```text
tokenScale = 1000
refillRate = 2 tokens/sec
refillUnitsPerSecond = 2000 units/sec
elapsedMs = 250

refilled = floor((250 * 2000) / 1000)
         = floor(500)
         = 500 units
         = 0.5 token
```

## Consuming a Request

One request costs one token:

```text
one request = tokenScale units
```

Allow branch:

```lua
tokens = tokens - tokenScale
return {1, math.floor(tokens / tokenScale), 0}
```

The returned remaining value is rounded down to whole tokens because `Result.Remaining` is an `int`.

## Retry-After Formula

On rejection, the bucket has less than one token.

```text
deficit = tokenScale - tokens
```

The wait time is:

```lua
local retryAfterSeconds = math.ceil(deficit / refillUnitsPerSecond)
```

Use `ceil`, not `floor`, because returning `0` seconds for a partial wait would tell the caller to
retry immediately before one full token exists.

Example:

```text
tokens = 500 units
tokenScale = 1000
refillUnitsPerSecond = 2000

deficit = 500 units
retryAfterSeconds = ceil(500 / 2000)
                  = ceil(0.25)
                  = 1 second
```

This is conservative because the public Redis contract currently returns seconds.

## Float vs Integer Tradeoff

### Float State

Example:

```text
tokens = 3.75
```

Benefits:

- Directly matches the math.
- Easier to read at first.
- Fewer explicit unit conversions.

Costs:

- Stores/parses float values in Redis.
- Redis state depends on floating-point representation.
- Harder to reason about exact comparisons near boundaries.

### Whole Integer State

Example:

```text
tokens = 3
```

Benefits:

- Very simple storage and comparisons.
- No scaling math.

Costs:

- Loses partial refill progress.
- A rate like `0.5 tokens/sec` cannot accumulate cleanly unless extra remainder state is added.

### Scaled Integer State

Example:

```text
tokens = 3750 units
```

Benefits:

- Preserves partial refill progress.
- Avoids float state in Redis.
- Keeps comparisons simple:

```lua
tokens >= tokenScale
```

Costs:

- More unit conversions.
- More naming discipline required.
- Scale choice is a precision/memory/range tradeoff.

## Scale Choice

Current scale:

```text
tokenScale = 1000
```

This gives millitoken precision.

Alternatives:

- `1` - whole tokens only; simplest, but loses fractional refill.
- `1000` - millitokens; good enough for this project and keeps numbers readable.
- `1_000_000` - microtokens; more precision, but larger numbers and more arithmetic risk.

## Current Rule of Thumb

Keep this mental model:

```text
At the public API boundary:
capacity = tokens
refillRate = tokens/sec
remaining = whole tokens
retryAfter = seconds

Inside Redis Token Bucket:
tokens = scaled units
time = milliseconds
refillRate = scaled units/sec
request cost = tokenScale
```

## Related

- [[D70 - Redis Token Bucket Scaled Integer State]]
- [[D71 - Redis Token Bucket Time and TTL Policy]]
- [[D72 - Redis Token Bucket Result Contract]]
- [[D73 - Redis Token Bucket Scaling Boundary]]
- [[R09 - Redis TIME and Unit Conversion]]

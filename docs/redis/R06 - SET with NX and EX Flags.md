# R06 - SET with NX and EX Flags

Back to [[Redis Concepts Index]].

## Concept

`SET` takes flags that fold a conditional and an expiry into the **single atomic command**:

```
SET key value NX EX 60
        ^^^^^ ^^ ^^^^^
        value NX  EX seconds
```

- `NX` — set **only if the key does not already exist**. (`XX` is the inverse: only if it does.)
- `EX seconds` / `PX milliseconds` — attach a TTL in the same command.

So `SET key value NX EX 60` means "create this key with a 60s TTL, but only if it isn't there
already" — atomically, in one round trip.

## The check-then-act race it removes

The naive instinct is two commands:

```
EXISTS key        # is it there?
SET key value EX 60   # if not, create it
```

Those are **two separate commands**. Between them, another client can create the key — and your
`SET` then stomps it, resetting both the value *and* the TTL mid-window. For a fixed-window
limiter that silently leaks requests at every busy window boundary.

> **Meta-pattern:** any decision expressed as *"check → then act"* across two commands is **not**
> atomic. Before reaching for a transaction or Lua, look for a **single command with a
> conditional flag** (`SET …NX`, `INCR`'s return value, `EXPIRE …NX`) that does it in one shot.
> This is the same class of bug as [[D46 - GetSet Race]], just on the Redis substrate.

## Related

- [[R01 - Key-Value Store and TTL]]
- [[R02 - Atomic Operations with INCR]]
- [[R07 - Levels of Atomicity in Redis]]
- [[D63 - Fixed Window in Redis via Lua]]

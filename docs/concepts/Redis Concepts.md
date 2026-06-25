# Redis Concepts

> Back to [[Rate Limiter Learning Map]] · See also [[Go Concepts]] · Decisions in [[README|ADR index]]

The Redis and Lua ideas behind the distributed `RedisLimiter` (V8), condensed and ordered from primitives to the full token-bucket arithmetic. Each entry names the decision(s) it supports.

## Core Redis

**Key-value store and TTL.** Redis is an in-memory key-value store where keys can carry a TTL (`SET key val EX seconds`); when it expires Redis deletes the key automatically (`TTL key` → seconds left, `-2` missing, `-1` no expiry). A fixed-window limiter can store the count as the value and let the TTL reset the window — no cleanup jobs (ADR-0063).

**Atomic operations with INCR.** `INCR key` atomically increments (creating at 0 first, so a missing key returns 1). Redis is single-threaded, so individual commands can't interleave. But rate limiting needs atomicity across *several* steps (`INCR` then, if first, `EXPIRE`) — and another client can slip between them. That gap is why Lua is needed (ADR-0063).

**SET with NX and EX flags.** `SET key val NX EX 60` folds a conditional and a TTL into one atomic command: `NX` = only if absent (`XX` = only if present), `EX`/`PX` = attach a TTL. It removes the check-then-act race of `EXISTS` then `SET` (between which another client can create the key, and your `SET` would stomp value *and* TTL mid-window). **Meta-pattern:** any "check → then act" across two commands is not atomic — first look for a single command with a conditional flag (`SET …NX`, `INCR`'s return, `EXPIRE …NX`). Same class of bug as the Get/Set race (ADR-0063).

**Redis hashes and HMGET.** A top-level key can hold a string, hash, list, set, etc. A hash stores named fields under one key; `HMGET key f1 f2` reads several. Redis token bucket needs two fields (`tokens`, `last_refill_ms`), so a hash keeps the state structured under one key instead of packing/parsing one string. `HMGET` returns a table; check the fields (missing → false-like) for cold start (ADR-0070).

## Levels of atomicity

**Levels of atomicity in Redis.** Three ways to make several steps atomic, by power: **(1)** a single command with flags — `SET …NX EX`, `INCR`/`DECR` (mutate *and* return) — cheapest, use whenever one command fits. **(2)** `MULTI`/`EXEC` transactions — run queued commands with no interleaving, but **cannot branch**: commands are queued before any run, so you can't "read count, and *if* over limit reject". `WATCH` adds optimistic locking (abort `EXEC` if a watched key changed) = CAS, the Redis twin of the in-process bounded retry loop. **(3)** Lua `EVAL` — the whole script runs atomically server-side *and* can read, branch, and conditionally write. **Rule:** the moment the logic is "read → compute → conditionally write back", the tool is Lua. The systems bridge: the in-memory limiter is *optimistic* (Get→Decide→CAS, retry); `RedisLimiter` is *pessimistic single-shot* (whole read-decide-write atomic in Lua — no conflict, nothing to retry, no `ErrCASConflict`) (ADR-0051, ADR-0062, ADR-0063).

## Lua in Redis

**Lua scripting for atomicity.** A script sent via `EVAL` runs atomically — every `redis.call` completes before any other client proceeds — giving cross-process read-decide-write atomicity without client-side CAS loops. Scripts use local variables, Lua control flow, and a final `return`. From Go you send the script text plus args; Redis executes and returns the result (ADR-0063).

**redis.call() and command execution.** Inside Lua, `redis.call('CMD', arg, …)` runs a command synchronously and returns its result (`INCR` → new int, `GET` → string or nil, `EXPIRE` → 1/0). A failing command raises a Lua error that aborts the script; `redis.pcall` catches it instead (ADR-0063).

**KEYS and ARGV parameters.** A script receives `KEYS[]` (the keys it touches) and `ARGV[]` (other arguments, often config). From `redis-cli`: `EVAL "script" numkeys k1 k2 a1 a2`; from go-redis you pass keys and args separately. Convention: put only real keys in `KEYS[]` and config (TTL, limits) in `ARGV[]` so Redis can track keys for cluster mode (ADR-0064).

**Redis TIME and unit conversion.** `redis.call("TIME")` returns `{seconds, microseconds}` — one authoritative server clock for the atomic decision, immune to per-instance skew. Token bucket converts to ms: `nowMs = sec*1000 + floor(usec/1000)`. Refill then converts ms back to seconds: `floor((elapsedMs * refillUnitsPerSecond) / 1000)`. Milliseconds are precise enough and keep numbers smaller than microseconds (ADR-0070, ADR-0071).

**Lua scripts embedded in Go strings.** Scripts live as Go raw strings, so the Go compiler only checks valid-string-ness, not Lua. Typos (`refilRate` vs `refillRate`), missing `local`, wrong `ARGV` index, or wrong return shape compile fine and fail only at execution. Keep Lua names explicit and unit-bearing (`refillUnitsPerSecond`, `elapsedMs`); a passing `go test ./...` proves the package compiles, not that the script runs — only a test against real Redis does (ADR-0073).

## Token bucket arithmetic

**Redis token bucket arithmetic.** Token bucket has more math than fixed window because refill is continuous and partial progress matters (250ms at 2/s = 0.5 token; discarding it breaks small frequent refills). The Redis path uses **scaled integers** — `tokenScale = 1000`, so 1 token = 1000 units — to preserve sub-token progress without storing floats. The public API stays in domain units (`capacity` tokens, `refillRate` tokens/sec); the algorithm converts to `capacityScaled` and `refillUnitsPerSecond`. Internally: tokens in scaled units, time in ms, one request costs `tokenScale`, so `tokens >= tokenScale` asks "at least one whole request token?". Refill: `refilled = floor((elapsedMs * refillUnitsPerSecond) / 1000)`. Allow: `tokens -= tokenScale`, return `{1, floor(tokens/tokenScale), 0}` (remaining floored to whole tokens for the `int` field). Reject: `deficit = tokenScale - tokens`, `retryAfter = ceil(deficit / refillUnitsPerSecond)` — `ceil`, not `floor`, so a partial wait never returns 0 and invites an immediate premature retry. Scale choice is a precision/range trade: `1` loses fractions, `1000` (millitokens) is readable and enough here, `1_000_000` adds precision but bigger numbers and more arithmetic risk (ADR-0070, ADR-0072, ADR-0073).

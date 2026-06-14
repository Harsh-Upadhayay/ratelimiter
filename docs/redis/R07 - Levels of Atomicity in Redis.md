# R07 - Levels of Atomicity in Redis

Back to [[Redis Concepts Index]].

## Concept

Redis is **single-threaded**: every individual command runs to completion with no interleaving,
so each command is atomic for free. A **sequence** of commands is *not* — other clients' commands
can slip between yours. There are three ways to make several steps atomic, ranked by power:

### 1. A single command with flags

When one command expresses the whole step:
- `SET key val NX EX 60` — conditional create + TTL in one shot ([[R06 - SET with NX and EX Flags]]).
- `INCR` / `DECR` — mutate **and return** the new value atomically ([[R02 - Atomic Operations with INCR]]).

Cheapest and clearest — use it whenever a single command fits.

### 2. MULTI / EXEC (transactions)

Queue several commands; `EXEC` runs them with no interleaving.

**Critical limitation: no branching.** The commands are queued *before* any of them run, so you
**cannot** read a value and decide what to do next based on it ("read count, and *if* over the
limit, reject"). That makes plain MULTI/EXEC useless for *read-decide-write* logic.

`WATCH key` adds optimistic locking: `EXEC` aborts if a watched key changed since `WATCH`. That is
**CAS** — the Redis twin of the in-process [[D51 - Bounded CAS Retry Loop]]. You retry on abort.

### 3. Lua scripting (EVAL)

The **whole script** executes atomically on the server, single-threaded, nothing interleaves — and
unlike transactions you **can** read a value, branch on it, and conditionally write, all inside one
atomic execution. See [[R03 - Lua Scripting for Atomicity]].

> **Rule:** the moment your decision is *"read a value → compute → conditionally write it back,"*
> the tool is **Lua**. Single commands are too narrow; MULTI/EXEC can't branch.

## The systems bridge: optimistic vs pessimistic

Same correctness goal (atomic read-modify-write), opposite strategies depending on substrate:

- **In-memory `Limiter`** — *optimistic*: `Get → Decide → CAS`, retry on conflict ([[D51 - Bounded CAS Retry Loop]]).
  Also what `WATCH`/`MULTI`/`EXEC` would give you in Redis.
- **`RedisLimiter`** — *pessimistic, single-shot*: the entire read-decide-write runs atomically
  server-side in Lua, so conflicts **can't** happen and there is **nothing to retry** — no loop,
  no `ErrCASConflict`. This is the "how" behind [[D62 - Redis as Limiter Not Store]].

## Related

- [[R03 - Lua Scripting for Atomicity]]
- [[R06 - SET with NX and EX Flags]]
- [[D62 - Redis as Limiter Not Store]]
- [[D63 - Fixed Window in Redis via Lua]]

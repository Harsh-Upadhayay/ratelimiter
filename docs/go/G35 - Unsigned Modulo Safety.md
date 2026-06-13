---
Back to [[Go Concepts Index]].

## Concept

Modulo (`%`) on a signed integer can return a negative result in Go. Modulo on an unsigned integer always returns a non-negative result.

## The problem

Hash functions often return `uint32` or `uint64`. Casting to `int` before modulo is dangerous: if the high bit is set, the cast produces a negative signed integer, and `negativeInt % n` returns a negative value in Go — an invalid array index.

## The fix

Keep the modulo in the unsigned domain:

```go
int(h.Sum32() % uint32(len(s.shards)))
```

`h.Sum32()` is already `uint32`. Cast `len(s.shards)` (which is a signed `int`) to `uint32` so both sides match. The final outer `int(...)` cast is safe because the result is guaranteed to be in `[0, len-1]`.

## Why not `int(h.Sum32()) % len(s.shards)`?

`int(h.Sum32())` can be negative on a 32-bit platform if the high bit is set, and the modulo result would be negative — an out-of-bounds index.

## Related

- [[G32 - Key Hashing]]
- [[D58 - Trust Internal Invariants]]

# G48 - Ceiling Duration Conversion

Back to [[Go Concepts Index]].

## Concept

`time.Duration` stores nanoseconds internally.

So this is not a seconds conversion:

```go
strconv.Itoa(int(duration))
```

It converts the raw nanosecond count.

## In this project

`Retry-After` needs whole seconds, rounded up.

The helper should clamp non-positive values and ceiling positive durations:

```text
d <= 0 -> 0
otherwise ceil(d / time.Second)
```

Integer duration math can express the ceiling:

```go
(d + time.Second - 1) / time.Second
```

Only use that after checking `d > 0`.

## Why not Seconds plus int

```go
int(d.Seconds())
```

truncates:

```text
500ms  -> 0
1500ms -> 1
```

For `Retry-After`, that is too aggressive. A client could retry before a full token/window is
available.

## Links

- [[D82 - HTTP Retry After Converts Duration to Seconds]]
- [[G23 - Duration Seconds Conversion]]
- [[D77 - Rate Limit HTTP Headers]]

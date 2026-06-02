# G23 - Duration Seconds Conversion

Back to [[Rate Limiter Learning Map]].

## Go practice

`time.Duration` can be converted to seconds with `Seconds()`.

## Rate limiter use

Token bucket refill uses elapsed seconds:

```text
refilledTokens = elapsed.Seconds() * refillRate
```

## Why it matters

The public configuration uses human-readable tokens per second, while `time.Duration` stores elapsed time internally.

## Links

- [[D36 - Token Bucket Capacity and Refill Rate]]
- [[D39 - Floating Point Token Arithmetic]]

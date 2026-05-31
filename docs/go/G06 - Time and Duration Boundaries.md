# G06 - Time and Duration Boundaries

Back to [[Rate Limiter Learning Map]].

## Go practice

Use `time.Time` for instants and `time.Duration` for elapsed time or window size.

## Boundary model

A half-open window includes the start and excludes the end. That means the exact expiry instant belongs to the next window.

## Why it matters

Rate limiter bugs often live at boundaries. Passing time explicitly makes those cases testable.

## Links

- [[D05 - Explicit Time Input]]
- [[D08 - Half Open Window Interval]]

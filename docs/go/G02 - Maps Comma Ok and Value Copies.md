# G02 - Maps Comma Ok and Value Copies

Back to [[Rate Limiter Learning Map]].

## Go practice

Use the comma-ok form when missing keys need to be distinguished from zero values.

Reading a missing key returns the zero value for the map's value type. The `ok` value tells you whether the key actually existed.

## Value-copy note

If a map stores struct values, reading a value gives you a copy. Mutating the local copy does not update the map entry unless it is assigned back.

## Why it matters

This affects the limiter's state update path: read state, decide, then write updated state back.

## Links

- [[D02 - State Shape Map String UserState]]
- [[D06 - Count Invariant Accepted Requests]]

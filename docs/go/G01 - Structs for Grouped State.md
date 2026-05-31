# G01 - Structs for Grouped State

Back to [[Rate Limiter Learning Map]].

## Go practice

Use a struct when several fields represent one concept.

For this limiter, request count and window start are not independent ideas. Together, they describe a user's current window state.

## Why it matters

A named struct makes the code communicate domain meaning instead of passing loose values around.

## Links

- [[D02 - State Shape Map String UserState]]
- [[D10 - Return Type]]

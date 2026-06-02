# C04 - V4 Algorithm Boundary Checkpoint

Back to [[Rate Limiter Learning Map]].

## Status

Work in progress checkpoint. Fixed-window behavior has been moved behind an internal algorithm interface. Token bucket has not been added yet.

## Implemented direction

- private `algorithm` interface
- marker `algorithmState` interface
- fixed-window-specific state
- fixed-window configuration constructor
- generic limiter orchestration
- explicit algorithm state mismatch error
- nil algorithm validation

## Known follow-up

Two existing constructor validation tests currently overwrite the error returned by `NewFixedWindow` by calling `NewLimiter` with the nil algorithm afterward.

Those tests should assert the algorithm-constructor errors directly.

## Next session

1. Fix the two validation test setups.
2. Add direct nil-algorithm coverage.
3. Polish stale comments.
4. Confirm the fixed-window suite is green.
5. Design token-bucket state and refill arithmetic.

## Links

- [[V4 Algorithm Abstraction Index]]
- [[D31 - Algorithm Owns Config Validation]]
- [[D32 - Manual Assembly with Built In Algorithms]]
- [[D33 - Decide Returns State Errors]]
- [[D34 - Initialize Missing State Before Assertion]]

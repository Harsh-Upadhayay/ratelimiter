# G14 - Pure Helper Functions

Back to [[Rate Limiter Learning Map]].

## Go practice

A helper function can keep logic deterministic by receiving all inputs as parameters and avoiding side effects.

## Rate limiter use

The fixed-window decision helper should not:

- read the map
- lock a mutex
- call `time.Now`
- mutate package-level state

It should only compute a result and next state from its inputs.

## Why it matters

Pure helpers are easier to reason about and make later interfaces easier to discover from concrete code.

## State writeback

The caller may write the returned state back even when a request is rejected. For rejected requests, the helper returns the unchanged state.

## Links

- [[D23 - Private Fixed Window Decision Helper]]
- [[D25 - Allow as Orchestrator]]
- [[D13 - Delay Interfaces]]

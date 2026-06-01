# D24 - Test Through Public API

Back to [[Rate Limiter Learning Map]].

## Context

V3 introduces a private helper, but the external behavior should not change.

## Decision

Keep tests focused on the public `Limiter` API.

## Why

Public API tests preserve refactoring freedom. The helper can change shape later without forcing test rewrites.

## Tradeoff

Some internal edge cases may be slightly harder to target directly, but current behavior is still observable through `Allow`.

## Links

- [[D23 - Private Fixed Window Decision Helper]]
- [[G15 - Public API Tests]]

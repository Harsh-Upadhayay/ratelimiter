# D46 - GetSet Race

Back to [[Rate Limiter Learning Map]].

## Context

A simple store could expose only `Get` and `Set`.

## Problem

If two requests read the same old state before either writes, both can be allowed even when only one should pass.

Example:

```text
state says one request remains
request A reads state
request B reads same state
both decide allowed
both write updated state
```

## Decision pressure

Separate `Get` and `Set` are not enough when decision happens between them without a shared atomic guard.

## Links

- [[D45 - Split Storage from Limiter]]
- [[D47 - StateStore Uses Get and CAS]]
- [[G04 - Mutexes and Critical Sections]]

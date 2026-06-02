# G22 - Nil Interface Guard

Back to [[Rate Limiter Learning Map]].

## Go practice

If a constructor accepts an interface dependency, reject a missing interface before storing it.

## Rate limiter use

The generic limiter constructor checks for a nil algorithm so `Allow` does not panic when calling `Decide`.

## Caveat

Go interfaces have subtle typed-nil behavior. A simple nil check handles a nil interface value, but a non-nil interface containing a typed nil pointer may require additional policy later.

## Links

- [[D31 - Algorithm Owns Config Validation]]
- [[D32 - Manual Assembly with Built In Algorithms]]

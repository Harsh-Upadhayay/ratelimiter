# G30 - Exported Interfaces With Unexported Types

Back to [[Rate Limiter Learning Map]].

## Concept

An exported interface can still be hard or impossible for outside packages to implement if its methods mention unexported types.

## Rate limiter use

`StateStore` is exported, but its methods use `algorithmState`, which is unexported.

That makes the storage boundary visible for learning, while still keeping algorithm state controlled inside the package.

## Tradeoff

This is acceptable in the current learning project, but it is not a polished public-library boundary.

A later public store API may need exported state envelopes, serialization boundaries, or a different constructor shape.

## Links

- [[D52 - Default Memory Store Constructor]]
- [[G08 - Exported Types with Unexported Fields]]
- [[G19 - Marker Interfaces and Opaque State]]

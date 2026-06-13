---
Back to [[Go Concepts Index]].

## Concept

`[]byte(key)` converts a `string` to a `[]byte` (slice of bytes).

## Why it exists

Strings in Go are immutable. Many I/O interfaces (like `io.Writer`, which `hash.Hash` embeds) accept `[]byte`, not `string`, because they deal in raw bytes and need a mutable, sliceable type. The conversion copies the underlying bytes.

## Syntax breakdown

- `[]byte` — the type: "slice of byte"
- `[]byte(key)` — a type conversion expression, not a function call
- The `[]` means slice-of, distinguishing it from `byte` (a single byte, alias for `uint8`)

## Usage in this project

```go
h.Write([]byte(key))  // feed the key's bytes into the fnv hasher
```

## Related

- [[G32 - Key Hashing]]
- [[G02 - Maps Comma Ok and Value Copies]]

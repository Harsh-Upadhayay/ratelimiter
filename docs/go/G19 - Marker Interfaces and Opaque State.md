# G19 - Marker Interfaces and Opaque State

Back to [[Rate Limiter Learning Map]].

## Go practice

A marker interface uses a method only to identify a controlled set of valid types.

## Rate limiter use

An unexported marker method lets the package define which concrete state types qualify as algorithm state.

The limiter stores those values without inspecting their fields.

## Why it matters

This keeps algorithm-specific state opaque to orchestration code while avoiding raw `any`.

## Tradeoff

Marker interfaces add ceremony and do not eliminate runtime type assertions inside each algorithm.

## Links

- [[D28 - Marker Interface for Algorithm State]]
- [[G20 - Type Assertions]]

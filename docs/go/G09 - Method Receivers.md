# G09 - Method Receivers

Back to [[Rate Limiter Learning Map]].

## Go practice

Go does not define methods inside a struct body.

Instead, a function becomes a method by declaring a receiver before the method name.

## Concept

The receiver is the object the method is attached to.

For a stateful type, use a pointer receiver when the method needs to mutate the object or avoid copying it.

## C++ comparison

In C++, methods are declared inside the class.

In Go, the struct holds fields, and methods are declared separately with a receiver.

## Links

- [[G03 - Pointer Receivers]]
- [[D17 - Constructor Returns Pointer]]

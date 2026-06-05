# V6 Storage Boundary Index

Back to [[Rate Limiter Learning Map]].

This hub tracks the move from limiter-owned in-memory state to a storage boundary.

## Core decisions

- [[D45 - Split Storage from Limiter]]
- [[D46 - GetSet Race]]
- [[D47 - StateStore Uses Get and CAS]]

## Prior pressure

- [[D25 - Allow as Orchestrator]]
- [[D26 - Introduce Algorithm Interface]]
- [[D44 - Split Files by Responsibility]]

## Go and systems concepts used

- [[G04 - Mutexes and Critical Sections]]
- [[G26 - Optimistic Concurrency with CAS]]

## Future bridge

- [[plan|Target 6-step plan]]

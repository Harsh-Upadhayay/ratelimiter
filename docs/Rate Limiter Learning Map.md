# Rate Limiter Learning Map

Obsidian entry point for the learning version of the rate limiter.

The notes are organised into two kinds of document:

- **ADRs** (`adr/`) — the architecture decisions (`ADR-0001` … `ADR-0085`) in standard
  Status/Context/Decision/Consequences format, grouped one log per iteration (V1–V9). Each
  log opens with a short narrative of that iteration, so reading a log top-to-bottom tells
  the story of that version.
- **Concepts** (`concepts/`) — the Go and Redis/Lua techniques the decisions rely on,
  condensed into themed prose. Read these to understand the "how"; read the ADRs for the
  "why".

## Start here

- [[README|ADR index]] — every decision by number, with status, in one table.
- [[Go Concepts]] — Go language and style techniques, by theme.
- [[Redis Concepts]] — Redis and Lua techniques for the distributed path.

## Decision logs by iteration

1. [[V1 - Fixed Window]] — first fixed-window counter; Go state, maps, time, locking, API boundary. (ADR-0001–0018)
2. [[V2 - API Evolution]] — boolean → `Result` + sentinel errors. (ADR-0019–0022)
3. [[V3 - Decision Logic]] — pure, separately testable decision helper. (ADR-0023–0025)
4. [[V4 - Algorithm Abstraction]] — one limiter, many algorithms behind a private interface. (ADR-0026–0034)
5. [[V5 - Token Bucket]] — second algorithm: lazy refill, float tokens, retry-after. (ADR-0035–0044)
6. [[V6 - Storage Boundary]] — StateStore interface, CAS, bounded retry loop. (ADR-0045–0052)
7. [[V7 - Sharded MemoryStore]] — lock striping, sharded backend, contract testing, the Redis pivot. (ADR-0053–0062)
8. [[V8 - RedisLimiter]] — parallel Redis limiter via Lua; server-side atomicity. (ADR-0063–0073)
9. [[V9 - HTTP Middleware]] — HTTP request-path integration; MemoryLimiter/Limiter fork. (ADR-0074–0085)

## Other notes

- [[plan|Target 6-step plan]] — the north-star design.
- [[mindmaps/V1 Learning Mindmap|V1 learning mindmap]]
- [[current_assessment|Interview / résumé readiness assessment]]

## Conventions

New decisions append an `ADR-00NN` section to the current iteration's log and a row to the
[[README|ADR index]]; new techniques append to [[Go Concepts]] or [[Redis Concepts]]. One file
per *iteration*, never one file per note. See `CLAUDE.md` → "Documentation conventions".

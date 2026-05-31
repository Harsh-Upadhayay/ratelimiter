# Rate Limiter Learning Map

This is the Obsidian entry point for the learning version of the rate limiter.

## Roadmap

- [[plan|Target 6-step plan]]
- [[mindmaps/V1 Learning Mindmap|V1 learning mindmap]]

## V1 decisions

- [[decisions/D01 - Fixed Window First]]
- [[decisions/D02 - State Shape Map String UserState]]
- [[decisions/D03 - Pointer Receiver for Mutating Limiter]]
- [[decisions/D04 - One Global Mutex]]
- [[decisions/D05 - Explicit Time Input]]
- [[decisions/D06 - Count Invariant Accepted Requests]]
- [[decisions/D07 - Window Reset Starts At Now]]
- [[decisions/D08 - Half Open Window Interval]]
- [[decisions/D09 - Defer Unlock After Lock]]
- [[decisions/D11 - Constructor Validation]]
- [[decisions/D12 - Exported API Boundary]]
- [[decisions/D13 - Delay Interfaces]]
- [[decisions/D14 - Package Name]]
- [[decisions/D15 - One File First]]
- [[decisions/D16 - Limiter Owns Runtime State]]
- [[decisions/D17 - Constructor Returns Pointer]]

## Go practice notes

- [[go/G01 - Structs for Grouped State]]
- [[go/G02 - Maps Comma Ok and Value Copies]]
- [[go/G03 - Pointer Receivers]]
- [[go/G04 - Mutexes and Critical Sections]]
- [[go/G05 - Defer Unlock Pattern]]
- [[go/G06 - Time and Duration Boundaries]]
- [[go/G07 - Do Not Copy Mutexes]]
- [[go/G08 - Exported Types with Unexported Fields]]

## Current open question

- [[decisions/D10 - Return Type]]

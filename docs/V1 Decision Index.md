# V1 Decision Index

Back to [[Rate Limiter Learning Map]].

This hub groups the first implementation decisions so the Obsidian graph is easier to browse.

## Implementation scope

- [[D01 - Fixed Window First]]
- [[D13 - Delay Interfaces]]
- [[D15 - One File First]]
- [[D16 - Limiter Owns Runtime State]]

## State and time

- [[D02 - State Shape Map String UserState]]
- [[D05 - Explicit Time Input]]
- [[D06 - Count Invariant Accepted Requests]]
- [[D07 - Window Reset Starts At Now]]
- [[D08 - Half Open Window Interval]]

## API boundary

- [[D10 - Return Type]]
- [[D11 - Constructor Validation]]
- [[D12 - Exported API Boundary]]
- [[D14 - Package Name]]
- [[D18 - Empty Key Validation]]

## Concurrency

- [[D03 - Pointer Receiver for Mutating Limiter]]
- [[D04 - One Global Mutex]]
- [[D09 - Defer Unlock After Lock]]
- [[D17 - Constructor Returns Pointer]]

## Next bridge

- [[V2 API Evolution Index]]

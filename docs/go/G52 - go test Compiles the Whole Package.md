# G52 - go test Compiles the Whole Package

Back to [[Go Concepts Index]].

## Concept

The unit of compilation for `go test` is the **package**, not the file. Running tests on a package:

1. compiles **every `*_test.go` file in that package together** into one test binary,
2. (only if step 1 succeeds) runs it,
3. executes the test functions matching `-run`.

`-run` filters **execution**, not **compilation**.

## The trap

Migrating one test file and trying to run just its tests:

```sh
go test -race -run 'TestEmptyKey|TestAllowConcurrent' ./...
```

If a *sibling* test file in the same package still has stale code, the output is:

```text
# package [package.test]
./other_test.go:30: assignment mismatch: 2 variables but helper returns 3 values
FAIL  package [build failed]
```

That is **zero tests run** — a build failure, not "all tests ran". `-run` never got a chance,
because the whole package failed to compile. You cannot run one file's tests while a sibling file
in the same package does not compile; they are welded into one binary.

## There is no "run one file"

Test selection is by **function name** (`-run`) or by **package** — never by file.

- `-run` is a **regex matched as a substring** (unanchored): `-run TestFoo` runs anything whose
  name *contains* `TestFoo`. Exact match: anchor it — `-run '^TestFoo$'`.
- `./...` selects all packages recursively; `.` selects the current package. Neither narrows below
  package level.

## Working around it

- Preferred: migrate all test files in the package so it compiles, *then* use `-run '^...$'` to
  focus execution.
- Escape hatch: temporarily rename an unmigrated file off the `_test.go` suffix
  (`mv x_test.go x_test.go.bak`) so the build skips it; restore later. Hacky; only worth it for a
  large file.

## Links

- [[G15 - Public API Tests]]
- [[G28 - Go Benchmarks]]
- [[G25 - Package Scope Across Files]]

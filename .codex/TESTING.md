# Verification, CI, and coverage

## Modules and developer commands

From repository root (Go 1.23 or a compatible later version):

```sh
make verify                # format check, go vet, go test, go test -race
make verify-all            # verify + nested adapter race tests
go test ./...
go test -race ./...
go vet ./...
go test -run '^$' -bench=. -benchmem ./...
go test -fuzz=FuzzParseLevel -fuzztime=30s .   # optional, separate explicit fuzz run
go build ./examples/...
```

Inside the nested module:

```sh
cd adapters/goconfig
go mod tidy
go vet ./...
go test -race ./...
go build ./examples/...
```

`make verify` does **not** exercise the adapter. `make verify-all` does include adapter race tests but omits adapter vet/build and dependency diff checks; CI handles those separately. `fmt-check` in the Makefile traverses all Go files including the adapter subtree.

**Avoid invoking `go mod tidy` without checking for resulting `go.mod`/`go.sum` changes** in reproducibility-sensitive contexts.

## CI jobs

Source: `.github/workflows/ci.yml`, GitHub Actions workflow `CI`.

Triggers:
- Push to `main`;
- All pull requests.

Job `core` (Go `1.23.x` on `ubuntu-latest`):
1. Checkout, go setup and module cache.
2. `go mod tidy` with `git diff --exit-code -- go.mod go.sum`.
3. `gofmt` check, excluding the adapter subtree.
4. `go vet ./...`.
5. `go test -race -coverprofile=coverage.out ./...`.
6. `go build ./examples/...`.

Job `goconfig-adapter` (working directory `adapters/goconfig`, Go `1.23.x`):
1. Tidy + dependency diff check.
2. `gofmt` check.
3. `go vet ./...`.
4. `go test -race ./...`.
5. `go build ./examples/...`.

The workflow records a coverage file for core tests, but this snapshot does not specify a coverage **threshold** or CI-enforced benchmark budgets.

## Evidence observed 2026-10-18

The most recent workflow run for source baseline `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881` was completed successfully:
- Run ID: **32868901673**
- Recorded run date: **2026-08-25**
- [GitHub Actions evidence](https://github.com/majiddarvishan/fluxlog/actions/runs/32868901673)

This is evidence of the remote **workflow result**, not evidence that tests were newly executed in this analysis session. The source inspection covered the complete tracked code tree, but no new `go test`, race run, benchmark or fuzz campaign was performed locally in this session.

## Existing automated tests

| Area | Tracked evidence |
| --- | --- |
| Config validation/aliases | `config_test.go` |
| Console format/color/unicode | `console_test.go` |
| Package functions, caller and lazy formatting | `default_test.go` |
| Structured JSON and per-instance dynamic levels | `logger_test.go` |
| 8 goroutines × 100 events with level flips | `logger_test.go:TestConcurrentLoggingAndLevelChanges` |
| Close idempotence and file parent creation | `logger_test.go` |
| Fuzz seeds/roundtrip level mapping | `level_fuzz_test.go` |
| Enabled/disabled benchmarks | `benchmark_test.go` |
| Adapter parsing, validation and binding setup | `adapters/goconfig/adapter_test.go` |

## Specific verification gaps (not proven defects)

- No tracked targeted test for **actual size threshold rotation**, max-backup pruning, max-age pruning or compression behavior.
- No explicit test that concurrent `Close` and active writes produce the intended bound/return behavior under the race detector.
- No adapter integration test exercising actual `Manager.OnReplace` success/failure **after** registration, including rollback.
- No explicit tests for error propagation from failing custom writers/closers, multiple output sink failures, or file permission failures.
- No quantified performance regression thresholds or sustained stress evidence; source benchmarks are available.
- No explicit test for an externally changed zerolog global minimum level interacting with a fluxlog instance.

Prioritize gap-focused tests before adding new abstractions or changing production APIs. These are recommendations for future verification, not CI failures.

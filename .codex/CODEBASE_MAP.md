# Codebase map

Baseline: `main` at `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881`. No `.codex` directory existed in that tracked tree before this documentation branch.

## Core Go module (repository root)

| Path | Responsibility |
| --- | --- |
| `go.mod`, `go.sum` | Root module declaration and pinned dependency graph |
| `config.go` | Public config types, format/color modes, defaults, validation, defensive normalization |
| `level.go` | Public levels/aliases and mapping to zerolog levels |
| `logger.go` | Constructor, multi-output wiring, public instance API, runtime level, caller context, lifecycle |
| `writer.go` | Mutex-protected output writer/close; atomic severity writer |
| `console.go` | Console output order, service/caller styling, time and color formatting |
| `default.go` | Atomic default instance and package-level convenience functions |
| `README.md` | Primary consumer documentation and samples |
| `CHANGELOG.md` | Historical release notes up to 0.2.1 |
| `LICENSE` | License terms (read the file for legal conditions) |
| `Makefile` | gofmt, go test, go test -race, vet, verification targets |
| `.gitignore` | Ignore *.log, coverage.out, dist/, tmp/ |
| `.github/workflows/ci.yml` | Root and adapter GitHub Actions verification |
| `examples/basic/main.go` | Instance logging and runtime level |
| `examples/file/main.go` | Console + rotating JSON file |
| `examples/package-level/main.go` | SetDefault and printf convenience API |

## Root tests

| File | Coverage |
| --- | --- |
| `config_test.go` | Validation, bad values, ParseLevel aliases |
| `console_test.go` | Console layout, service coloring, NO_COLOR, Unicode-aware caller shortening |
| `default_test.go` | Convenience formats, disabled argument evaluation, default nil rejection, caller site |
| `logger_test.go` | JSON events, instance level isolation, concurrent logging and updates, close idempotence, disabled events, file creation |
| `level_fuzz_test.go` | Fuzz seeds and invariant: successfully parsed level maps to a valid zerolog level |
| `benchmark_test.go` | Enabled/disabled event microbenchmarks, without committed benchmark results |
| `example_test.go` | Go doc-style example of the public package |

## Separate `goconfig` Go module

| Path | Responsibility |
| --- | --- |
| `adapters/goconfig/go.mod`, `go.sum` | Nested module and its own dependencies |
| `adapters/goconfig/adapter.go` | Parse legacy configuration, construct Logger, subscribe to level replacement |
| `adapters/goconfig/adapter_test.go` | Legacy mode parsing, optional fields, invalid mode/types, bind/application validation |
| `adapters/goconfig/README.md` | Provider configuration and install instructions |
| `adapters/goconfig/examples/basic/main.go` | Executable adapter example |

## Recommended change locations

- Public level types/aliases: `level.go` + `config_test.go`/fuzz tests.
- Log routing/close behavior: `writer.go`, `logger.go`, and concurrency tests.
- Console presentation: `console.go`, `console_test.go`.
- Default/package behavior: `default.go`, `default_test.go`.
- Config schema and defaults: `config.go` and adapter mapping/test when needed.
- CI/developer commands: `Makefile`, `.github/workflows/ci.yml`.
- A change to the adapter may require its **independent module's** version and tests, even when the core has not changed.

Source tree enumerated with the GitHub recursive tree endpoint at the baseline SHA (non-truncated).

# Established design invariants

These are **observations grounded in the current code**, not retrospective claims about why maintainers made a particular decision. Changing one requires explicit consideration of compatibility and documentation.

| ID | Established invariant | Evidence |
| --- | --- | --- |
| D-01 | The core module must not require a configuration-provider framework. | Root `go.mod` has no goconfig; optional `adapters/goconfig/go.mod` is separate |
| D-02 | Logger minimum levels are instance-scoped. The library does not call `zerolog.SetGlobalLevel`. | `writer.go:levelWriter`, `logger.go:SetLevel`, `logger_test.go:TestRuntimeLevelIsPerInstance` |
| D-03 | At least one output is required; console and file may coexist. | `config.go:normalized`, `logger.go:New`, tests |
| D-04 | Writes and `Close` are serialized; `Close` is idempotent. | `writer.go:managedWriter`, `logger.go:Close`, tests |
| D-05 | The caller owns the Logger lifecycle, including after `SetDefault`. | `default.go:SetDefault`, README |
| D-06 | Errors from setup/validation/binding are returned instead of calling Fatal internally. | `config.go`, `logger.go:New`, `adapter.go` |
| D-07 | Package convenience funcs optimize disabled events and report application caller location. | `default.go`, `logger.go:packageLevelLogger`, `default_test.go`, CHANGELOG 0.2.1 |
| D-08 | Console-specific service and caller presentation is not applied to JSON fields. | `console.go`, `logger.go`, `console_test.go` |
| D-09 | The adapter supports legacy flat configuration keys and binds **only** `level` to manager replacement. | `adapters/goconfig/adapter.go` and README |
| D-10 | File rotation is delegated to `lumberjack.v2` instead of a home-grown retention engine. | `logger.go:New`, root `go.mod` |

## Compatibility checklist

Before modifying runtime behavior, explicitly answer:

1. Does this change log encoding or field identity (especially timestamp, caller, service or level)?
2. Does it introduce process-global side effects, contention changes, or implicit ownership?
3. Does it alter disabled-level formatting or package-level caller frame accounting?
4. Does it change file layout/rotation/permissions or a file's colorlessness?
5. Does it affect legacy adapter field names, types, defaults or replace/rollback semantics?
6. Are core **and** adapter module tests, dependencies and examples updated?
7. Are docs, `CHANGELOG.md`, and semver implications considered?

## Evidence vs proposals

- The invariants above are **implemented** at baseline `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881`.
- `TASKS.md` lists recommended improvements; they are **not** accepted feature commitments, bug confirmations or historical product decisions.
- Do not assume no other system calls `zerolog.SetGlobalLevel`; fluxlog only controls its own code.

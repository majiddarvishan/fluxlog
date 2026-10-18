# Product overview and scope

## Purpose

Fluxlog supplies a small integration layer for **structured logging in Go services and libraries**: standardized JSON/human console output, optional file rotation, per-logger runtime severity changes, and explicit resource ownership, while retaining advanced `zerolog` APIs.

It is **library code** (`package fluxlog`), with runnable examples rather than its own long-lived server. Minimum declared Go version: **1.23**.

## Shipped capabilities

| Capability | Implementation | Notes |
| --- | --- | --- |
| Structured events | `Logger.Info().Str(...).Msg(...)` and other `*zerolog.Event` methods | Native zerolog fields and fluent API |
| Levels | trace, debug, info, warn, error, fatal, panic, disabled | Alias inputs `warning` → warn; `off` → disabled |
| Runtime severity | `SetLevel`, `SetLevelString` | Atomic, per `Logger`; no call to `zerolog.SetGlobalLevel` |
| Outputs | Console stream and/or rotating file | At least one is mandatory |
| Serialization | JSON or zerolog-style console format | File defaults to JSON |
| Console identity | timestamp, caller, `service` | Caller left-truncated in console; JSON preserves full caller |
| File retention/rotation | `lumberjack.v2` | Size, backup count, age, compression and local time |
| Package functions | `Trace/Debug/Info/Warn/Error/Fatal/Panic` | Atomic default pointer; convenient printf rules |
| Advanced integration | `Zerolog()`, `With()`, `WithContext()` | Use zerolog-specific event/context features |
| Configuration integration | `adapters/goconfig` nested module | Translates legacy flat keys and binds level replacement |

## Non-goals / not implemented in this snapshot

- No network transport, aggregation server, OpenTelemetry exporter, asynchronous buffer, or queue.
- No generic config provider dependency in the core module.
- No automatic hot reload of arbitrary core `Config` fields; the optional adapter dynamically binds **only** the `level` node.
- No built-in sensitive-data redaction, sampling policy, telemetry export, or distributed tracing orchestration.
- No publicly declared hard performance SLO; microbenchmarks exist but no measured results are captured in this documentation.
- No CLI executable; `examples/` contain independent demo programs.

## Consumers and typical flow

The application constructs one or more `Logger` instances, chooses console/file destinations, optionally creates a default logger for package functions, writes structured events, changes runtime severity if needed, and calls `Close` when finished. Configuration-provider users import the separate `goconfig` module.

## Dependencies and releases

- Root module: `github.com/majiddarvishan/fluxlog`; Go `1.23`; direct dependencies `github.com/rs/zerolog v1.35.1` and `gopkg.in/natefinch/lumberjack.v2 v2.2.1`.
- Nested module: `github.com/majiddarvishan/fluxlog/adapters/goconfig`; Go `1.23`; depends on fluxlog `v0.2.1` and `github.com/majiddarvishan/goconfig v1.3.0`.
- `CHANGELOG.md`'s latest entry is **0.2.1 (2026-08-25)**. Treat published tags/releases separately from changelog text; not all tag references were independently verified in this audit.

Evidence: `README.md`, `go.mod`, `adapters/goconfig/go.mod`, `CHANGELOG.md`, `logger.go`, `writer.go`, `default.go`, `adapter.go`.

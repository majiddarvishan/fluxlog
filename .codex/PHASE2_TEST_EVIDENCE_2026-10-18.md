# Fluxlog — Phase 2 executable reliability evidence

**Date:** 2026-10-18  
**Repository:** `github.com/majiddarvishan/fluxlog`  
**Baseline:** `main@7c6e796397da72b05260681cb9ba1b3cf47b4746`  
**Work:** [PR #14 — Lifecycle, rotation and output-error tests](https://github.com/majiddarvishan/fluxlog/pull/14)  
**Verified GitHub Actions:** [run 39158993645](https://github.com/majiddarvishan/fluxlog/actions/runs/39158993645): core and goconfig-adapter both **PASS**.

## Change boundary

Four new root-module Go test files only, plus updates to English `.codex` records. **No production `.go`, Go module, CI workflow, runtime config or public API changed. No release or tag was created.**

New files:

- `writer_lifecycle_integration_test.go` — `LIFE-01`
- `file_rotation_integration_test.go` — size rollover and option forwarding (`FILE-01`)
- `file_retention_integration_test.go` — actual max-backups, max-age and gzip (`FILE-01`)
- `output_error_integration_test.go` — `ERR-01`

## Observed test results

| Task | Test | Outcome | Verified guarantee |
| --- | --- | --- | --- |
| LIFE-01 | `TestLIFE01CloseWaitsForActiveWrite` | PASS | Close waits for an active write; owned resources close once; post-close WriteLevel returns ErrClosedPipe |
| LIFE-01 | `TestLIFE01ConcurrentWritesAndClose` | PASS | Race-detector stress; writes either complete or report ErrClosedPipe |
| LIFE-01 | `TestLIFE01SetDefaultDoesNotCloseReplacedLogger` | PASS | Concurrent default logger swaps do not transfer ownership or close a replaced logger |
| FILE-01 | `TestFILE01SizeRotationPreservesRecords` | PASS | Actual 1MiB size rollover preserves old message in backup and new message in active file |
| FILE-01 | `TestFILE01RotationConfigurationIsForwarded` | PASS | MaxSize, MaxBackups, MaxAge, Compress and LocalTime reach lumberjack unchanged |
| FILE-01 | `TestFILE01MaxBackupsIsEnforced` | PASS | Real asynchronous cleanup retains only the newest backup under MaxBackups=1 |
| FILE-01 | `TestFILE01MaxAgeAndCompression` | PASS | Older-than-MaxAge backup removed; retained backup gzip-compressed and decompressed bytes match |
| ERR-01 | `TestERR01CustomWriterFailureAndShortWrite` | PASS | Custom write error and io.ErrShortWrite propagate through the managed output boundary |
| ERR-01 | `TestERR01MultiSinkPartialFailure` | PASS | Earlier sink failure is surfaced, while later file sink still receives the event |
| ERR-01 | `TestERR01CreateLogDirectoryFailure` | PASS | Invalid parent-directory construction fails deterministically |

Both CI jobs performed `go mod tidy` consistency, formatting checks, `go vet`, `go test -race -v` and example builds. Core coverage: **75.9%** of statements in the CI run (measurement, not a goal). These are **remote CI observations**, not locally executed tests.

## Test design considerations

- Synchronization uses channels and a bounded timeout for blocked Close/Write, never arbitrary sleeps.
- Real async lumberjack cleanup is observed through bounded, condition-driven filesystem polling.
- Each file test runs in an isolated `t.TempDir()`. Permission-failure tests use a regular file as a directory, avoiding privileged assumptions.
- An error from an application's `io.Writer` is visible at Fluxlog's `managedWriter` boundary. The public zerolog `Event.Msg` method instead reports output errors through zerolog's `ErrorHandler` (or stderr); the test does not claim `Msg` returns an error.
- `FileConfig.LocalTime` is checked for propagation to lumberjack, not by changing the process-global timezone or asserting exact wall-clock backup filenames. That avoids global-state manipulation.

## Outstanding unrelated work

`CFG-02` and `CFG-03` remain observed/known-issue **SKIPs**, not fixed. `FMT-02` is still open for additional console-color contract checks; `REL-01` and `DOC-01` remain release/documentation tasks. `PERF-01` and `CI-01` remain optional. No broad logger rewrite or async architecture is recommended merely because these tests passed.

---

## FMT-02 console color contract — 2026-10-18 (PR #15)

New file: `console_contract_integration_test.go`, with three test functions and Auto/Always/Never subcases. No runtime code, dependencies or CI workflow changed.

- `TestFMT02ConsoleColorModes`: **PASS**, with `auto-buffer`, `always`, `never` subcases; buffered output does not accidentally gain ANSI, while explicit AlwaysColor does when `NO_COLOR` is unset.
- `TestFMT02NoColorDisablesServiceLevelAndCallerColors`: **PASS**; `NO_COLOR=1` suppresses ANSI while retaining level, caller, service, request ID and message.
- `TestFMT02DualOutputKeepsFilePlainText`: **PASS**; a colored console stream and a simultaneous human-readable file sink contain the expected fields, but the file has no ANSI escape codes.
- **Evidence:** [GitHub Actions run 39159526896](https://github.com/majiddarvishan/fluxlog/actions/runs/39159526896): `core` and `goconfig-adapter` **PASS**, including Go 1.23 race/vet/format/tidy/build. Core coverage **76.9%** in this run.
- **Result:** FMT-02 is complete with no production changes. Only CFG-02/CFG-03 remain known-issue SKIPs from the earlier Phase 1 behavioral tests.

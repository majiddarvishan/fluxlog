# Fluxlog — Phase 1 executable regression evidence

**Date:** 2026-10-18  
**Baseline (before changes):** `main@78274ee19a1668a71ac563b0e64b69b86b9f0ff9`  
**PR:** [#12 — Phase 1 regression characterization](https://github.com/majiddarvishan/fluxlog/pull/12)  
**Dependency:** `goconfig v1.6.0` / `zerolog v1.35.1`  
**Evidence:** [Actions run 39146082868](https://github.com/majiddarvishan/fluxlog/actions/runs/39146082868), both `core` and `goconfig-adapter` jobs **success**. An earlier, also successful, verbose run [39145924506](https://github.com/majiddarvishan/fluxlog/actions/runs/39145924506) provided the detailed console timestamp observations.

## What was actually run

Both module CI jobs performed dependency `go mod tidy` consistency, `gofmt`, `go vet`, `go test -race -v` (root with coverage profile) and example builds. Verbose mode in CI deliberately exposes `PASS` vs `SKIP` for known defects. This evidence is from **remote GitHub Actions**, not an independent local test run.

## Verified outcome matrix

| Task | Test(s) | GitHub Actions result | Technical meaning |
| --- | --- | --- | --- |
| **CFG-01** | `TestCFG01BindLevelRuntimeReplacementAndInvalidRollback`, `TestCFG01FileSourcePersistsValidLevelAndRejectsInvalid` | **PASS (2 tests)** | A valid replacement updates logger, live node and persisted FileSource; an invalid value is rejected and does not change them. |
| **CFG-02** | `TestCFG02MultipleLoggersOnOneLevelNode` | **SKIP (known issue)** | One goconfig `level` node bound to two fluxlog loggers: first transitions to debug, second stays info. Upstream callback dispatch is first-match, not multicast. |
| **CFG-03** | `TestCFG03FilePersistenceFailureDoesNotLeakRuntimeLevel` | **SKIP (known defect)** | A directory at `file.tmp` forces a deterministic persistence failure; original file and node remain info, but the first successful callback leaves runtime logger debug. |
| **FMT-01** | `TestFMT01ConsoleUnixTimestampFormats`, four cases | **SKIP (4 subtests)** | Numeric timestamp data is dropped from console for seconds, milliseconds, microseconds and nanoseconds. JSON output parity is checked **before** the skip. |
| **FMT-01** | `TestFMT01ConsoleRFC3339Timestamp` | **PASS** | The supported RFC3339 console timestamp is preserved. |
| **LVL-01** | `TestLVL01ExternalZerologGlobalThreshold` | **PASS** | Global zerolog filtering can prevent debug output while fluxlog's internal configured level remains debug. Original global level is restored after the test. |

**Important interpretation:** Exactly **six known-issue SKIPs** appear across these focused tests (CFG-02: 1, CFG-03: 1, FMT-01: 4). The overall CI jobs succeed because these cases use explicit `t.Skipf` on the known deviation. **A green CI does not mean these six defects are fixed.** If the behavior is subsequently corrected, these tests are written to pass without a skip. A different unexpected result produces a failure.

## Code-change boundaries

- Added test-only files: `adapters/goconfig/bind_level_integration_test.go`, `console_timestamp_integration_test.go`, `logger_global_level_test.go`.
- CI only: changed existing race-test commands to `go test -race -v` to retain visible PASS/SKIP evidence, without adding duplicate test runs.
- README/documents: clarify the effective zerolog global threshold, track genuine defects and next decisions.
- **Not changed:** `adapter.go`, `logger.go`, `writer.go`, `console.go`, root `go.mod`, adapter `go.mod` or dependency versions; no releases or tags.

## Required decisions before runtime fixes

1. **CFG-02:** Is multiple-listener delivery on one `goconfig` path required? If yes, coordinate an upstream handler dispatch/registration change; otherwise document and enforce one binding per node.
2. **CFG-03:** Must a failed configuration persistence attempt restore **external logger state**? If yes, design an explicit compensating/transactional API contract with goconfig. Avoid partial fixes that can leave persisted configuration and runtime inconsistent.
3. **FMT-01:** Confirm Unix timestamp formats are part of Fluxlog's console compatibility contract. If yes, accept `json.Number` as the upstream console formatter already does.

Until these decisions are made, keep their tasks open. Next: `LIFE-01`, `FILE-01`, `ERR-01` only if the user asks for further test expansion. No behavior-changing fix is authorized in this test-only iteration.

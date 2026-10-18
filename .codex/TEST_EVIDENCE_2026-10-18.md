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

---

## FMT-01 fix verification — 2026-10-18 (PR #13)

This addendum **supersedes only FMT-01's historical SKIP outcomes above**; earlier results remain traceable to their original CI runs.

- **Authorized behavior:** show timestamps as human-readable local time in console output, using `ConsoleConfig.TimeFormat`. Default is RFC3339, while `"2006-01-02 15:04:05"` is a documented simpler 24-hour layout.
- **Core code change:** `console.go` now parses zerolog `json.Number` timestamps for Unix seconds, milliseconds, microseconds and nanoseconds; fractional precision and negative epochs are handled. JSON output is unchanged.
- **Passing checks:** `TestFMT01ConsoleUnixTimestampFormats` (all four cases PASS), `TestFMT01ConsoleRFC3339Timestamp`, `TestFMT01CustomHumanReadableTimeLayout`, `TestFMT01UnixTimestampPrecisionAndNegativeEpoch` (all four cases PASS), `TestFMT01LoggerRendersUnixTimestampAsHumanReadable`.
- **CI evidence:** [Actions run 39147165603](https://github.com/majiddarvishan/fluxlog/actions/runs/39147165603), both core and adapter jobs PASS, including fmt, vet, race tests and examples.
- **Remaining known-issue SKIPs:** **two**, CFG-02 and CFG-03. They are not fixed by this change. No tag or release created.


---

## CFG-02 / CFG-03 implementation verification — 2026-10-19

This section supersedes the **historical** known-issue SKIP statuses for
CFG-02 and CFG-03 earlier in this record, but only for development source
containing the upstream post-commit integration, **not** the published
adapter `v0.1.0` or goconfig `v1.6.0`.

- **Dependency:** immutable goconfig commit `20a58eb9a426d0ed7d522b1e5848a0971ce58198`, encoded as Go pseudo-version `v1.6.1-0.20261018202246-20a58eb9a426`. Generated sums: [bootstrap run 39165338586](https://github.com/majiddarvishan/fluxlog/actions/runs/39165338586).
- **CFG-01:** both normal and FileSource validation/rollback tests **PASS**.
- **CFG-02:** `TestCFG02MultipleLoggersOnOneLevelNode` and `TestCFG02BatchMultipleLoggersAndInvalidRollback` **PASS** — shared binding receives the committed level for every Logger.
- **CFG-03:** `TestCFG03FilePersistenceFailureDoesNotLeakRuntimeLevel` and `TestCFG03FileSourceBatchFailureKeepsBothLoggersUnchanged` **PASS** — no external Logger state leak after a failed persist.
- **Overall:** root and adapter CI jobs passed `go mod tidy`, gofmt, go vet, `go test -race` and example builds in [Actions run 39165664802](https://github.com/majiddarvishan/fluxlog/actions/runs/39165664802). **No CFG SKIPs remain.**
- **Scope caveat:** legacy goconfig OnReplace uses the first pre-commit handler registered on a path. Fluxlog validation must own that first handler; mixing an unrelated validator registered earlier is unsupported. Post-commit listeners are retained for the Manager lifetime.
- **No new tag or release** was created for goconfig or either Fluxlog module.

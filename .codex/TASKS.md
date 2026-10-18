# Fluxlog — prioritized engineering backlog

**Planning date:** 2026-10-18  
**Scope:** `github.com/majiddarvishan/fluxlog` (core plus independently versioned `adapters/goconfig`).  
**Maintainer decision:** Use **`github.com/majiddarvishan/goconfig v1.6.0`**. All other items are **planned**, not authorized as implemented changes.  
**Tracking:** task IDs below are stable; update checkboxes only when acceptance criteria have actually been met.

## Current execution boundary

- **In scope now:** dependency upgrade to goconfig v1.6.0, canonical module/checksum files, compatibility verification, this English plan and handoff documentation.
- **Deferred:** every runtime implementation, new test, public API modification, release/tag, and optimization described below. Those remain separate follow-up tasks.
- **Not requested:** asynchronous logging, new exporters, queueing, broad architecture rewrites or new providers.

## P0 — dependency migration and baseline verification

- [x] **DEP-01 — Pin goconfig v1.6.0 in the adapter module.** The tag `v1.6.0` exists in the upstream goconfig repository. Update **`adapters/goconfig/go.mod`**, regenerate `go.sum` with real Go tooling, verify `go mod tidy` is clean; do not guess checksums. Leave root `go.mod` untouched. **Acceptance:** adapter module graph is reproducible, the referenced upstream version is exactly 1.6.0, and no new runtime source code has changed. **Status:** complete in PR #10; tidy and exact checksums verified by CI run 39144064913.
- [x] **DEP-02 — Validate adapter and core with Go 1.23 CI.** Run adapter `go mod tidy`, `go test -race ./...`, `go vet ./...`, `go build ./examples/...`; root `make verify` and `go build ./examples/...`. Confirm CI is green on the final PR commit. **Acceptance:** both modules pass; record exact runs. **Depends on:** DEP-01. **Status:** both root and adapter CI jobs passed on commit `27cfb35b20158e5fa59050f309c41aeccd7307b0` (run 39144064913).

## P1 — correctness and existing API contracts (tests first)

- [x] **INV-01 — Source-level behavior characterization.** Reviewed pinned goconfig v1.6.0 mutation/OnReplace/FileSource code and zerolog v1.35.1 formatting/global-level code. Recorded evidence and reproducible test scenarios in [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md). **Important:** this is **not** execution of the regression tests below.

- [ ] **CFG-01 — Test runtime goconfig replacement on v1.6.0.** Exercise `Manager.Replace` after `BindLevel`: successful `info -> debug`, invalid level rejection and rollback; check new/current level and persisted node consistency. **Acceptance:** focused integration tests reproduce expected states (including invalid input) under `-race`. **Depends on:** DEP-02.
- [ ] **CFG-02 — Characterize multi-logger binding and callback lifetime.** Bind two loggers to one manager level node, replace once and observe both; examine `OnReplace` registration and whether manager retains callbacks after Logger.Close. v1.6.0 still appends handler registrations and selects a matching path, so **do not assume** fan-out or deregistration. **Acceptance:** tests + explicit documented support policy; only then propose minimal fix if a supported use case fails. **Depends on:** CFG-01. **Source review:** OBS-01 confirms the first-match handler selection in goconfig v1.6.0; still requires test and contract decision.
- [ ] **CFG-03 — Persistence-failure rollback of logger side effects.** After a successful `BindLevel` callback and a failed `FileSource.setConfig`, the goconfig Node is restored but the fluxlog runtime level may remain changed. Introduce deterministic failure-injection verification, specify transactional consistency requirements and choose a minimal fix only if approved. **Acceptance:** reproducible error-path test with node, file and runtime level assertions. **Source review:** OBS-02. **Depends on:** CFG-01.
- [ ] **LVL-01 — Define effective per-instance threshold with zerolog global level.** Test a fluxlog logger after external `zerolog.SetGlobalLevel` calls (restore original global state; avoid parallel pollution). Decide whether the contract promises independence from **other fluxlog instances** or also from **external zerolog global filtering**. **Acceptance:** reproducible test and explicit contract; defer code changes until the expected policy is approved.
- [ ] **FMT-01 — Characterize Unix-format timestamps.** Test console formatting when `zerolog.TimeFieldFormat` uses numeric Unix variants, versus RFC3339 strings; also confirm JSON output is unaffected. **Acceptance:** expected output documented and regression tests; avoid changing behavior until reproduced.

## P2 — lifecycle and sink reliability

- [ ] **LIFE-01 — Concurrent Close/Write and default ownership tests.** Stress active writer vs Close, multiple Close calls, post-close events and `SetDefault` pointer swaps; include the race detector. **Acceptance:** no races or deadlocks, deterministic ownership expectations and controlled error outcomes.
- [ ] **FILE-01 — Real rotation and retention tests.** Use temp dirs to verify size rollover, MaxBackups, MaxAgeDays, Compress and LocalTime where applicable. **Acceptance:** deterministic tests with bounded input and no arbitrary sleeps or dependency on external log directories.
- [ ] **ERR-01 — Output-failure tests.** Cover failing custom `io.Writer`, parent directory/permission failure and partial multi-output failure. **Acceptance:** tests document observed error propagation and partial-write semantics; improvements considered only if an actual defect is identified.
- [ ] **FMT-02 — Complete ANSI/color contract tests.** Test Auto/Always/Never modes, `NO_COLOR`, file console-formatted output, and service/caller formatting. **Acceptance:** clear compatibility contract with tests on a CI-compatible platform.

## P3 — release hygiene and optional optimization (after correctness)

- [ ] **REL-01 — Reconcile adapter release/version instructions.** Verified Git tag list includes `adapters/goconfig/v0.1.0` but **not** `adapters/goconfig/v0.2.1` as of 2026-10-18. Root README advertises adapter 0.1.0; adapter README advertises 0.2.1. Choose a truthful installation example referencing an actually available **adapter module tag**, or explicitly document branch-only installation; issue a new adapter tag **only on explicit approval** after successful validation. **Acceptance:** documentation and published tag state agree.
- [ ] **DOC-01 — Update public docs after accepted behavior decisions.** Bring READMEs, configuration recipes, release notes and `.codex` snapshots into agreement, including dependency version and any supported multi-binding semantics. **Acceptance:** examples build against published versions or clearly describe development-only paths.
- [ ] **PERF-01 — Optional reproducible microbenchmark baseline.** Compare enabled/disabled-event benchmarks, multi-sink overhead and sink contention with pinned Go toolchain/hardware and a common benchmark harness. **Acceptance:** reproducible raw results; no speculative performance rewrites.
- [ ] **CI-01 — Revisit fuzz/coverage policy only if maintenance needs justify it.** Avoid introducing more CI gates without a concrete requirement.

## Suggested execution order

```text
DEP-01 -> DEP-02
               |
               +-> CFG-01 -> CFG-02
               |          +-> CFG-03
               +-> LVL-01
               +-> FMT-01
               |
               +-> LIFE-01 -> ERR-01
               +-> FILE-01
               +-> FMT-02
               |
               +-> REL-01 -> DOC-01
               +-> PERF-01 / CI-01 (optional)
```

## Work-item discipline

Each follow-up should be narrowly scoped: first reproduce (or document expected behavior), then decide whether code change is warranted, then add regression tests and verify both modules. Prefer one focused PR per coherent correction. **Do not merge speculative fixes** and do not tag a release implicitly. After PR merge, verify `main`, remove temporary remote branches if possible, and update `CURRENT_STATUS.md` / `SESSION.md`.

See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for phases, change boundaries and final acceptance gates. Read [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md) for source-backed risks and precise test scenarios.

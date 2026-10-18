# Fluxlog — prioritized engineering backlog

**Planning date:** 2026-10-18  
**Scope:** `github.com/majiddarvishan/fluxlog` (core plus independently versioned `adapters/goconfig`).  
**Maintainer decision:** Use **`github.com/majiddarvishan/goconfig v1.6.0`**. Phase 1 characterization is complete. The maintainer explicitly approved human-readable console time, and FMT-01 was implemented without modifying unrelated behavior.  
**Tracking:** task IDs below are stable; update checkboxes only when acceptance criteria have actually been met.

## Current execution boundary

- **Completed:** dependency upgrade to v1.6.0, Phase 1 characterization, FMT-01 human-readable Unix timestamp fix, Phase 2 LIFE-01/FILE-01/ERR-01/FMT-02 verification, and English CI evidence.
- **Deferred:** CFG-02/CFG-03 behavior changes, unrelated API modifications, release/tag or speculative optimization. See [TEST_EVIDENCE_2026-10-18.md](TEST_EVIDENCE_2026-10-18.md).
- **Not requested:** asynchronous logging, new exporters, queueing, broad architecture rewrites or new providers.

## P0 — dependency migration and baseline verification

- [x] **DEP-01 — Pin goconfig v1.6.0 in the adapter module.** The tag `v1.6.0` exists in the upstream goconfig repository. Update **`adapters/goconfig/go.mod`**, regenerate `go.sum` with real Go tooling, verify `go mod tidy` is clean; do not guess checksums. Leave root `go.mod` untouched. **Acceptance:** adapter module graph is reproducible, the referenced upstream version is exactly 1.6.0, and no new runtime source code has changed. **Status:** complete in PR #10; tidy and exact checksums verified by CI run 39144064913.
- [x] **DEP-02 — Validate adapter and core with Go 1.23 CI.** Run adapter `go mod tidy`, `go test -race ./...`, `go vet ./...`, `go build ./examples/...`; root `make verify` and `go build ./examples/...`. Confirm CI is green on the final PR commit. **Acceptance:** both modules pass; record exact runs. **Depends on:** DEP-01. **Status:** both root and adapter CI jobs passed on commit `27cfb35b20158e5fa59050f309c41aeccd7307b0` (run 39144064913).

## P1 — correctness and existing API contracts (tests first)

- [x] **INV-01 — Source-level behavior characterization.** Reviewed pinned goconfig v1.6.0 mutation/OnReplace/FileSource code and zerolog v1.35.1 formatting/global-level code. Recorded evidence and reproducible test scenarios in [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md). **Important:** this is **not** execution of the regression tests below.

- [x] **CFG-01 — Test runtime goconfig replacement on v1.6.0.** Exercise `Manager.Replace` after `BindLevel`: successful `info -> debug`, invalid level rejection and rollback; check new/current level and persisted node consistency. **Acceptance:** focused integration tests reproduce expected states (including invalid input) under `-race`. **Depends on:** DEP-02. **Evidence:** 2 PASS in Actions run 39146082868 (including real FileSource persistence).
- [ ] **CFG-02 — Characterize multi-logger binding and callback lifetime.** Bind two loggers to one manager level node, replace once and observe both; examine `OnReplace` registration and whether manager retains callbacks after Logger.Close. v1.6.0 still appends handler registrations and selects a matching path, so **do not assume** fan-out or deregistration. **Acceptance:** tests + explicit documented support policy; only then propose minimal fix if a supported use case fails. **Depends on:** CFG-01. **Source review:** OBS-01 confirms the first-match handler selection in goconfig v1.6.0; Observed in CI: SKIP with second logger stuck at info; requires explicit supported contract and fix decision.
- [ ] **CFG-03 — Persistence-failure rollback of logger side effects.** After a successful `BindLevel` callback and a failed `FileSource.setConfig`, the goconfig Node is restored but the fluxlog runtime level may remain changed. Introduce deterministic failure-injection verification, specify transactional consistency requirements and choose a minimal fix only if approved. **Acceptance:** reproducible error-path test with node, file and runtime level assertions. **Executable evidence:** SKIP in Actions run 39146082868 (runtime debug, node/file info). **Depends on:** CFG-01. Remains open pending expected consistency policy/fix.
- [x] **LVL-01 — Define effective per-instance threshold with zerolog global level.** Test a fluxlog logger after external `zerolog.SetGlobalLevel` calls (restore original global state; avoid parallel pollution). Decide whether the contract promises independence from **other fluxlog instances** or also from **external zerolog global filtering**. **Acceptance:** reproducible test and explicit contract. **Evidence:** PASS in Actions run 39146082868; README explicitly documents that external zerolog global filtering remains effective. No logger implementation change.
- [x] **FMT-01 — Render Unix-format timestamps as human-readable local time.** The console formatter now converts zerolog's integer Unix seconds/milliseconds/microseconds/nanoseconds using `ConsoleConfig.TimeFormat`, including negative and fractional values, while preserving the RFC3339 default and unchanged JSON output. **Evidence:** all FMT-01 cases, custom 24-hour layout and full Logger pipeline PASS in [Actions run 39147165603](https://github.com/majiddarvishan/fluxlog/actions/runs/39147165603); PR #13. Earlier four SKIPs in run 39145924506 are historical regression evidence.

## P2 — lifecycle and sink reliability

- [x] **LIFE-01 — Concurrent Close/Write and default ownership tests.** Stress active writer vs Close, multiple Close calls, post-close events and `SetDefault` pointer swaps; include the race detector. **Acceptance:** no races or deadlocks, deterministic ownership expectations and controlled error outcomes. **Evidence:** 3 PASS in Actions [run 39158993645](https://github.com/majiddarvishan/fluxlog/actions/runs/39158993645); includes blocked active writer, concurrent Close/Write, default logger swap and post-close error.
- [x] **FILE-01 — Real rotation and retention tests.** Use temp dirs to verify size rollover, MaxBackups, MaxAgeDays, Compress and LocalTime where applicable. **Acceptance:** temporary-directory tests and bounded observable-completion polling for asynchronous retention. **Evidence:** 4 PASS in Actions run 39158993645: size rotation, option propagation, actual MaxBackups cleanup, MaxAge deletion and gzip content. **Caveat:** LocalTime filename timezone only checked via configuration forwarding.
- [x] **ERR-01 — Output-failure tests.** Cover failing custom `io.Writer`, parent directory/permission failure and partial multi-output failure. **Acceptance:** 3 PASS (with two subcases) in Actions run 39158993645: failing/short stream, partial multi-sink failure and invalid directory. Event.Msg error reporting remains zerolog-owned; no new Fluxlog API.
- [x] **FMT-02 — Complete ANSI/color contract tests.** Auto/Always/Never for non-terminal output, `NO_COLOR` overriding AlwaysColor, service/caller fields, and dual console+file sinks are verified. **Acceptance:** all FMT-02 subtests PASS in [CI 39159526896](https://github.com/majiddarvishan/fluxlog/actions/runs/39159526896), both modules passing race, vet, formatting and example builds; file console format contains no ANSI codes even when stream output is colored.

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

See [PHASE2_TEST_EVIDENCE_2026-10-18.md](PHASE2_TEST_EVIDENCE_2026-10-18.md) for the verified Phase 2 results, [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for phases and gates, [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md) for source paths, and [TEST_EVIDENCE_2026-10-18.md](TEST_EVIDENCE_2026-10-18.md) for actual PASS/SKIP outcomes.

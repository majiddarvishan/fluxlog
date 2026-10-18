# Fluxlog — prioritized engineering backlog

**Planning date:** 2026-10-18  
**Scope:** `github.com/majiddarvishan/fluxlog` (core plus independently versioned `adapters/goconfig`).  
**Lockstep policy (2026-10-19):** Starting with core v0.3.0, all official adapters share the same version number; the next nested tag is `adapters/goconfig/v0.3.0`, not v0.2.0. Modules stay separate and may point to different release commits. Historical notes are retained.
**Historical release decision:** migrate the already-tested post-commit adapter integration to the official tagged upstream `goconfig v1.7.0`. The prior pseudo-version is recorded as historical development evidence. Phase 1 characterization is complete. The maintainer explicitly approved human-readable console time, and FMT-01 was implemented without modifying unrelated behavior.  
**Tracking:** task IDs below are stable; update checkboxes only when acceptance criteria have actually been met.

## Current execution boundary

- **Completed:** dependency upgrade to v1.6.0, Phase 1 characterization, FMT-01 human-readable Unix timestamp fix, Phase 2 LIFE-01/FILE-01/ERR-01/FMT-02 verification, and English CI evidence.
- **Current change:** replace the previously tested goconfig source-commit pseudo-version with published upstream `v1.7.0` in the nested adapter module, verify both modules and prepare independent release gates. **Deferred:** unrelated API changes or speculative optimization. **REL-01** is complete for the existing release. **DOC-01** is complete for the current supported source contract; future release notes require separate tag authorization. See [TEST_EVIDENCE_2026-10-18.md](TEST_EVIDENCE_2026-10-18.md).
- **Not requested:** asynchronous logging, new exporters, queueing, broad architecture rewrites or new providers.

## Release gates — after CFG-02 and CFG-03 are fixed

- [x] **REL-02 — Publish and verify upstream goconfig v1.7.0.** Annotated Git tag resolves to `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`; upstream [CI 39167180503](https://github.com/majiddarvishan/goconfig/actions/runs/39167180503) passed. Go module origin and sums verified in [run 39167874434](https://github.com/majiddarvishan/fluxlog/actions/runs/39167874434). A separate GitHub Release object has not been published.
- [x] **REL-03 — Adopt the tagged upstream dependency in the nested adapter.** Nested `go.mod` now uses official `goconfig v1.7.0` and the sums verified by Go. Strict CFG-01/CFG-02/CFG-03 cases PASS without SKIP. Both modules passed dependency consistency, format, vet, race tests and example builds in [PR CI run 39168203987](https://github.com/majiddarvishan/fluxlog/actions/runs/39168203987). Merge and post-merge CI must still be verified as a separate process gate.
- [x] **REL-04 — Publish Fluxlog core.** Existing annotated `v0.3.0` tag resolves to `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd` and both jobs passed post-merge CI 39168420174. Go proxy resolution remains an independent check. Never recreate or move the tag.
- [ ] **REL-05 — Publish the lockstep-versioned adapter.** Update the separate module's core dependency to `v0.3.0`, regenerate sums with Go, pass both CI jobs and strict CFG-02/CFG-03, merge, check post-merge CI, and tag `adapters/goconfig/v0.3.0` on the tested main commit. Do not publish adapter v0.2.0.

## P0 — dependency migration and baseline verification

- [x] **DEP-01 — Pin goconfig v1.6.0 in the adapter module.** The tag `v1.6.0` exists in the upstream goconfig repository. Update **`adapters/goconfig/go.mod`**, regenerate `go.sum` with real Go tooling, verify `go mod tidy` is clean; do not guess checksums. Leave root `go.mod` untouched. **Acceptance:** adapter module graph is reproducible, the referenced upstream version is exactly 1.6.0, and no new runtime source code has changed. **Status:** complete in PR #10; tidy and exact checksums verified by CI run 39144064913.
- [x] **DEP-02 — Validate adapter and core with Go 1.23 CI.** Run adapter `go mod tidy`, `go test -race ./...`, `go vet ./...`, `go build ./examples/...`; root `make verify` and `go build ./examples/...`. Confirm CI is green on the final PR commit. **Acceptance:** both modules pass; record exact runs. **Depends on:** DEP-01. **Status:** both root and adapter CI jobs passed on commit `27cfb35b20158e5fa59050f309c41aeccd7307b0` (run 39144064913).

## P1 — correctness and existing API contracts (tests first)

- [x] **UP-01 — Add goconfig post-commit replacement notifications upstream.** Merged [goconfig PR #1](https://github.com/majiddarvishan/goconfig/pull/1) at `main@20a58eb9a426d0ed7d522b1e5848a0971ce58198`. An additive `OnReplaceCommitted` hook broadcasts after successful single or batch persistence. Verified Go race/vet/build/format in [Actions run 39162857364](https://github.com/majiddarvishan/goconfig/actions/runs/39162857364). **No new goconfig tag and no Fluxlog integration yet.**
- [x] **UP-02 — Integrate and verify upstream post-commit API in the Fluxlog adapter.** Implementation is staged in the feature branch with an immutable Go pseudo-version pinned to upstream commit `20a58eb9a426` (checksum generated by Go in [dependency resolution run 39165338586](https://github.com/majiddarvishan/fluxlog/actions/runs/39165338586)); pure validation + committed level application and strict no-SKIP CFG-02/CFG-03 tests. **Acceptance:** both modules passed format, tidy, vet, race and example builds in [CI run 39165664802](https://github.com/majiddarvishan/fluxlog/actions/runs/39165664802). The dependency is deliberately an immutable development pseudo-version; replace it with a real upstream tag only after a separately authorized release. See [CFG02_CFG03_INTEGRATION.md](CFG02_CFG03_INTEGRATION.md).


- [x] **INV-01 — Source-level behavior characterization.** Reviewed pinned goconfig v1.6.0 mutation/OnReplace/FileSource code and zerolog v1.35.1 formatting/global-level code. Recorded evidence and reproducible test scenarios in [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md). **Important:** this is **not** execution of the regression tests below.

- [x] **CFG-01 — Test runtime goconfig replacement on v1.6.0.** Exercise `Manager.Replace` after `BindLevel`: successful `info -> debug`, invalid level rejection and rollback; check new/current level and persisted node consistency. **Acceptance:** focused integration tests reproduce expected states (including invalid input) under `-race`. **Depends on:** DEP-02. **Evidence:** 2 PASS in Actions run 39146082868 (including real FileSource persistence).
- [x] **CFG-02 — Verify multi-logger post-commit updates and callback lifetime.** Bind two loggers to one manager level node, replace once and observe both; examine `OnReplace` registration and whether manager retains callbacks after Logger.Close. v1.6.0 still appends handler registrations and selects a matching path, so **do not assume** fan-out or deregistration. **Acceptance:** tests + explicit documented support policy; only then propose minimal fix if a supported use case fails. **Depends on:** CFG-01. **Source review:** OBS-01 confirms the first-match handler selection in goconfig v1.6.0; Historical v1.6.0 CI SKIP is addressed by a new committed-observer implementation in this feature branch; the test now fails instead of SKIP on a mismatch. **Verified:** strict multi-logger Replace and final-value Batch tests PASS in CI run 39165664802. Legacy first-OnReplace registration caveat and observer lifetime are documented.
- [x] **CFG-03 — Persistence-failure rollback of logger side effects.** After a successful `BindLevel` callback and a failed `FileSource.setConfig`, the goconfig Node is restored but the fluxlog runtime level may remain changed. Introduce deterministic failure-injection verification, specify transactional consistency requirements and choose a minimal fix only if approved. **Acceptance:** reproducible error-path test with node, file and runtime level assertions. **Executable evidence:** SKIP in Actions run 39146082868 (runtime debug, node/file info). **Depends on:** CFG-01. The feature branch now applies levels only after successful commit; original and batch failure-injection tests contain no SKIP. **Verified:** both FileSource failed Replace and failed Batch tests PASS in CI run 39165664802 with no SKIP.
- [x] **LVL-01 — Define effective per-instance threshold with zerolog global level.** Test a fluxlog logger after external `zerolog.SetGlobalLevel` calls (restore original global state; avoid parallel pollution). Decide whether the contract promises independence from **other fluxlog instances** or also from **external zerolog global filtering**. **Acceptance:** reproducible test and explicit contract. **Evidence:** PASS in Actions run 39146082868; README explicitly documents that external zerolog global filtering remains effective. No logger implementation change.
- [x] **FMT-01 — Render Unix-format timestamps as human-readable local time.** The console formatter now converts zerolog's integer Unix seconds/milliseconds/microseconds/nanoseconds using `ConsoleConfig.TimeFormat`, including negative and fractional values, while preserving the RFC3339 default and unchanged JSON output. **Evidence:** all FMT-01 cases, custom 24-hour layout and full Logger pipeline PASS in [Actions run 39147165603](https://github.com/majiddarvishan/fluxlog/actions/runs/39147165603); PR #13. Earlier four SKIPs in run 39145924506 are historical regression evidence.

## P2 — lifecycle and sink reliability

- [x] **LIFE-01 — Concurrent Close/Write and default ownership tests.** Stress active writer vs Close, multiple Close calls, post-close events and `SetDefault` pointer swaps; include the race detector. **Acceptance:** no races or deadlocks, deterministic ownership expectations and controlled error outcomes. **Evidence:** 3 PASS in Actions [run 39158993645](https://github.com/majiddarvishan/fluxlog/actions/runs/39158993645); includes blocked active writer, concurrent Close/Write, default logger swap and post-close error.
- [x] **FILE-01 — Real rotation and retention tests.** Use temp dirs to verify size rollover, MaxBackups, MaxAgeDays, Compress and LocalTime where applicable. **Acceptance:** temporary-directory tests and bounded observable-completion polling for asynchronous retention. **Evidence:** 4 PASS in Actions run 39158993645: size rotation, option propagation, actual MaxBackups cleanup, MaxAge deletion and gzip content. **Caveat:** LocalTime filename timezone only checked via configuration forwarding.
- [x] **ERR-01 — Output-failure tests.** Cover failing custom `io.Writer`, parent directory/permission failure and partial multi-output failure. **Acceptance:** 3 PASS (with two subcases) in Actions run 39158993645: failing/short stream, partial multi-sink failure and invalid directory. Event.Msg error reporting remains zerolog-owned; no new Fluxlog API.
- [x] **FMT-02 — Complete ANSI/color contract tests.** Auto/Always/Never for non-terminal output, `NO_COLOR` overriding AlwaysColor, service/caller fields, and dual console+file sinks are verified. **Acceptance:** all FMT-02 subtests PASS in [CI 39159526896](https://github.com/majiddarvishan/fluxlog/actions/runs/39159526896), both modules passing race, vet, formatting and example builds; file console format contains no ANSI codes even when stream output is colored.

## P3 — release hygiene and optional optimization (after correctness)

- [x] **REL-01 — Reconcile adapter release/version instructions.** Verified repository tags: root `v0.1.0`/`v0.2.0`/`v0.2.1`; adapter **only** `adapters/goconfig/v0.1.0` as of 2026-10-18. Both READMEs now recommend `adapters/goconfig@v0.1.0` when installing a tagged adapter, and distinguish the published dependency graph (core `v0.1.0`, goconfig `v1.3.0`) from current `main` (core `v0.2.1`, goconfig `v1.6.0`). The development-only optional adapter configuration fields are clearly labeled. **Acceptance:** install references match verified tags; no new tag created. Documentation PR #16.
- [x] **DOC-01 — Document the accepted development-source behavior.** Both README files and `.codex` now distinguish the published adapter from the immutable pre-release dependency, describe the post-commit notification semantics, successful/failed Replace and Batch behavior, first-registered pre-commit handler limitation, and observer lifetime. **Evidence:** strict CFG-02/CFG-03 tests and both module CI jobs PASS in [run 39165664802](https://github.com/majiddarvishan/fluxlog/actions/runs/39165664802). New published-version release notes remain separately gated on explicit tag authorization.
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

See [CFG02_CFG03_INTEGRATION.md](CFG02_CFG03_INTEGRATION.md) for the staged, unshipped upstream fix and integration gates. See [PHASE2_TEST_EVIDENCE_2026-10-18.md](PHASE2_TEST_EVIDENCE_2026-10-18.md) for the verified Phase 2 results, [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for phases and gates, [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md) for source paths, and [TEST_EVIDENCE_2026-10-18.md](TEST_EVIDENCE_2026-10-18.md) for actual PASS/SKIP outcomes.

# Prioritized follow-ups / engineering backlog

This is a **proposed backlog**, derived from source review, not a record of accepted product requirements. Do not modify runtime code merely to clear this list without prioritization.

## P1 — Documentation correctness and release hygiene

- [ ] **DOC-01 — Verify adapter module tag/version instructions.** Root `README.md` says `adapters/goconfig@v0.1.0`; adapter README says `@v0.2.1`, while nested `go.mod` requires root `v0.2.1`. Verify published module tags and reconcile instructions, with a documentation-only change.
- [ ] **DOC-02 — Validate the new `.codex` snapshot after merge.** Confirm `main` contains these files and update `CURRENT_STATUS.md` with merged SHA and CI results.
- [ ] **DOC-03 — Keep changelog/release language precise.** Check real Git tags before claiming adapter versions are installable; GitHub releases endpoint returned no release objects during the review.

## P2 — Contract-critical test gaps

- [ ] **TEST-01 — Rotation/retention integration tests.** Verify actual rollover, max backups, age and compression across file-sink lifecycle; use temporary directories, bounded test data and OS-compatible assertions.
- [ ] **TEST-02 — Close/write race tests.** Assert concurrency guarantees for active writes, repeated Close, writes after Close, and returned errors.
- [ ] **TEST-03 — Adapter live replacement.** Exercise actual manager updates after `BindLevel`, both success and invalid-value rollback, and document listener ownership/lifecycle.
- [ ] **TEST-04 — Error paths.** Cover failing console writers, filesystem permissions and closers without terminating applications.
- [ ] **TEST-05 — Level/format boundary tests.** Validate external zerolog global-level interaction and complete ANSI behavior under `NO_COLOR`/color modes before changing color or filter code.

## P3 — Optional improvements after requirements exist

- [ ] **PERF-01 — Baseline benchmarks.** Run `BenchmarkEnabledLog` and `BenchmarkDisabledLog` under pinned Go version/hardware and preserve comparable data. Do not assert a speedup without matched evidence.
- [ ] **CI-01 — Consider automated fuzz smoke or an explicit coverage gate** if maintaining this library warrants the additional CI overhead.
- [ ] **DOC-04 — Optional public API recipes** for `WithContext`, multi-output service integration and graceful shutdown, using tests as the reference.

## Constraints

- Preserve source-backed invariants in [DECISIONS.md](DECISIONS.md).
- Core must remain independent of `goconfig`.
- Always test **both** modules. Never confuse remote CI success for local execution.
- No implicit new feature commitments, breaking changes, tags, or releases.

## Task status as of this review

No runtime tasks were implemented by this documentation effort. These unchecked items are suggestions; prioritization is still a maintainer decision.

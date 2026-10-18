# Fluxlog — next chat handoff

You are a **Senior Go Software Architect / Logging and Reliability Engineer** working on `git@github.com:majiddarvishan/fluxlog.git`.

**Language policy:** Respond to the maintainer in **Persian**, but keep **all files and code comments written to the repository in English**.

## Last verified state — 2026-10-18

- Default branch: `main`, baseline prior to Phase 1 test PR `78274ee19a1668a71ac563b0e64b69b86b9f0ff9`.
- PRs #9, #10, #11 merged. Optional nested `adapters/goconfig` module uses `goconfig v1.6.0`; core has no goconfig dependency.
- Current work in [PR #12](https://github.com/majiddarvishan/fluxlog/pull/12) adds regression **tests only**, verbose CI, README global-filter caveat and updated English `.codex` records.
- [Actions run 39146082868](https://github.com/majiddarvishan/fluxlog/actions/runs/39146082868) succeeded (core and adapter tidy, fmt, vet, race tests, examples). Final PR commit and post-merge run must be checked separately.

## Phase 1 outcomes

- **CFG-01 DONE:** live replacement plus invalid rollback, including a real file, PASS.
- **LVL-01 DONE:** zerolog's external global threshold applies before fluxlog's per-instance writer, PASS, documented in README.
- **CFG-02 OPEN:** two bindings on one goconfig node: first logger updates, second stays unchanged, **SKIP**.
- **CFG-03 OPEN:** failed FileSource persistence reverts config but not external Logger level, **SKIP**.
- **FMT-01 OPEN:** console timestamps using Unix sec/milli/micro/nano omit timestamp, **four SKIPs**; string RFC3339 PASS.
- **Total six known-issue SKIPs** across targeted scenarios. Green CI **does not** mean these bugs are fixed.

Read [`.codex/TEST_EVIDENCE_2026-10-18.md`](TEST_EVIDENCE_2026-10-18.md), [TASKS.md](TASKS.md), [CURRENT_STATUS.md](CURRENT_STATUS.md), [REVIEW_2026-10-18.md](REVIEW_2026-10-18.md), [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) and [ARCHITECTURE.md](ARCHITECTURE.md) as needed. Do **not** re-audit the whole repository. Compare changed files from known baseline.

## Constraints and next actions

1. First check actual `main`, PR #12, CI, active branches and whether the plan files were merged.
2. Do not change production Go implementation without a new explicit authorization; the current user's direction was to continue testing, not silently fix known runtime defects.
3. Define expected behavior for CFG-02 (one vs multiple callbacks), CFG-03 (persistence/side-effect transaction) and FMT-01 (Unix format support) before designing a fix.
4. For authorized fixes, use minimal changes, exercise both Go modules via Go 1.23 CI, update evidence and tasks; never silently tag or release.
5. After merging, verify `main`, CI and **branch cleanup**. If delete is unsupported by the connected GitHub interface, tell the maintainer the exact remote deletion command.

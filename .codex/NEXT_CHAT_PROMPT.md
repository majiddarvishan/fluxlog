# Fluxlog — next chat handoff

Act as a **Senior Go Software Architect / Structured Logging Reliability Engineer**. Continue the project at `git@github.com:majiddarvishan/fluxlog.git`.

## Grounded state (2026-10-18)

- Default branch: `main`.
- Verified main HEAD **before Phase 1 review**: `3323e1316764477b2bb1644aa397cef2fe00d575` (dependency upgrade PR #10 merged).
- Core runtime audit baseline: `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881`; documentation PR #9 merged without changing runtime code.
- PR [#10](https://github.com/majiddarvishan/fluxlog/pull/10) is **merged**; historical post-merge [CI 39144285233](https://github.com/majiddarvishan/fluxlog/actions/runs/39144285233) passed both jobs.
- See `.codex/REVIEW_2026-10-18.md` for the pinned-upstream Phase 1 findings; verify whether its documentation PR is merged before relying on `main` links.
- **Maintainer choice:** Use `github.com/majiddarvishan/goconfig v1.6.0` in the **nested adapter module**. **Other improvements remain tests/fixes to be explicitly authorized**; the current continuation was a source-only review with no runtime source edits.
- Two Go modules: core root, separate `adapters/goconfig`. Core must never gain goconfig as a dependency.
- Latest recorded core changelog entry is v0.2.1. The adapter has a verified Git tag `adapters/goconfig/v0.1.0` but not `adapters/goconfig/v0.2.1`.

## Required first steps

1. Inspect actual `main` SHA and open PRs; PR #10 is already merged. Do not assume newer docs or fixes are merged.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `TASKS.md`, `IMPLEMENTATION_PLAN.md`, `REVIEW_2026-10-18.md`, and `SESSION.md`; use architecture/configuration docs only where relevant.
3. Dependency upgrade is verified. Source-only review found first-match goconfig callbacks, no compensation for failed persistence after callbacks, numeric timestamp mismatch, and global zerolog filtering. These have **not** been reproduced by new executable tests. Do not re-run an entire codebase audit or claim source inspection is a passing test.
4. For a new implementation request, follow CFG-01 -> CFG-03 -> CFG-02, FMT-01 and LVL-01 from the recorded review. Until then, keep behavior changes deferred.
5. After any PR merge, verify main HEAD and CI, update status/handoff, and delete remote temporary branches when supported; if not, report the manual cleanup command.
6. Do not re-audit all source files: inspect the diff from the recorded baseline and narrowly relevant code.
7. Keep repository docs and plans **in English**, but explain progress to the user **in Persian**. Report only tests actually executed and changes truly merged.

## Core architecture (short reference)

Fluxlog is a synchronous, structured Go logger built on zerolog. Logger-specific atomic minimum-level filtering sends events to console and optional lumberjack-rotated file sinks with mutex-coordinated writer lifecycle; package functions use an atomic default pointer. The optional goconfig adapter parses a legacy flat config object and binds the manager's `level` node to runtime level changes.

## Delivery standard

For each task, show evidence, actual commit/PR, tested commands or CI links, remaining blockers and an accurate state (planned, implemented, verified, merged). Never claim published adapter tags based on changelog text alone.

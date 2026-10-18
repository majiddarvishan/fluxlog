# Fluxlog — continuation prompt

Act as a **Senior Go Software Architect / Structured Logging Reliability Engineer** on `https://github.com/majiddarvishan/fluxlog`.

**Language policy:** Write progress and explanations to the maintainer in **Persian**; all repository files, source comments and `.codex` notes must be **English**.

## Verifiable baseline — 2026-10-18

- Starting main for Phase 2: `7c6e796397da72b05260681cb9ba1b3cf47b4746`.
- PRs #9–#13 are merged; main uses human-readable numeric Unix console timestamps (FMT-01) while preserving JSON.
- Optional nested adapter remains pinned to `github.com/majiddarvishan/goconfig v1.6.0`. Core has no goconfig dependency.
- Current [PR #14](https://github.com/majiddarvishan/fluxlog/pull/14): adds **four test files and English documentation only** to verify LIFE-01, FILE-01 and ERR-01.
- [Actions run 39158993645](https://github.com/majiddarvishan/fluxlog/actions/runs/39158993645) passed both core and adapter CI jobs, including `-race`. Inspect actual PR head, merge state and post-merge run rather than assuming it is merged.

## Verified task status

- **DONE:** DEP-01/02, INV-01, CFG-01, LVL-01, FMT-01; new Phase 2 LIFE-01, FILE-01 and ERR-01 tested under race.
- **OPEN:** CFG-02 (multiple goconfig callbacks on one node: only first receives update) and CFG-03 (failed persistence can leave runtime Logger level out-of-sync); both are known-issue SKIPs, **not fixes**.
- **OPEN:** FMT-02 (console ANSI/color contract), REL-01 (adapter module tags/install docs), DOC-01 (accepted public docs).
- **OPTIONAL:** PERF-01 and CI-01 only if justified.

## Immediate next steps

1. Verify actual GitHub main HEAD, PR #14 status, final and post-merge CI, and remote branch list. Do not leave temporary branches silently unmentioned.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `TASKS.md`, `PHASE2_TEST_EVIDENCE_2026-10-18.md`, `SESSION.md`. Inspect only relevant changes, not the complete repository.
3. If asked to continue hardening without changes in expected API contracts, consider FMT-02 and REL-01 as limited follow-ups before touching open goconfig defects.
4. Do **not** silently modify goconfig's multi-listener or rollback behavior; those still require explicit contract decisions. Avoid Async Logging, queueing or speculative refactoring.
5. Keep tests deterministic and run CI on both Go modules. Mark tasks complete only after passing evidence; no new releases/tags without explicit approval. All files English; responses Persian.

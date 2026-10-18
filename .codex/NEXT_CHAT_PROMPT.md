# Fluxlog — next session handoff

Act as a **Senior Go Software Architect / Logging Reliability Engineer** on `https://github.com/majiddarvishan/fluxlog`.

**Communication policy:** Explain work to the maintainer in **Persian**. Keep every repository file, `.codex` document and source comment **in English**.

## Verified context (2026-10-18)

- Main before FMT-02: `5c972345043b5e315c8125a444f4c68754ffbaf9`.
- PRs #9–#14 are **merged**. Last post-merge CI #14: [39159303947](https://github.com/majiddarvishan/fluxlog/actions/runs/39159303947), PASS.
- Current [PR #15](https://github.com/majiddarvishan/fluxlog/pull/15) adds only FMT-02 console ANSI/color contract tests and docs. [CI 39159526896](https://github.com/majiddarvishan/fluxlog/actions/runs/39159526896) passed both modules; check final commit/merge/post-merge status.
- Nested optional adapter remains pinned to `goconfig v1.6.0`, core isolated.
- Human-readable console timestamps support numeric Unix units and configured layout, RFC3339 by default; JSON output unchanged.

## Task status

**Completed:** DEP-01/02, INV-01, CFG-01, LVL-01, FMT-01; Phase 2 **LIFE-01, FILE-01, ERR-01, FMT-02**. Evidence in `.codex/PHASE2_TEST_EVIDENCE_2026-10-18.md` and `TASKS.md`.

**Open known defects:** CFG-02 (only first OnReplace callback gets invoked for multiple logger bindings at one goconfig path), CFG-03 (FileSource persistence failure rolls back config but leaves changed Logger level). Their regression tests are **SKIP**, not successful fixes. No contract decision yet.

**Open administrative tasks:** REL-01 (nested module version/install instruction consistency), DOC-01 (accepted contract documentation). PERF-01 and CI-01 optional; no performance-driven rewrite requested.

## Working method

1. Confirm current GitHub main HEAD, PR #15 merge state and CI, remote branch cleanup before changing code.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `TASKS.md`, `PHASE2_TEST_EVIDENCE_2026-10-18.md`, and `SESSION.md`; inspect only changed areas since this baseline.
3. Prioritize REL-01 documentation/tag accuracy. Do not create a new release or tag without authorization.
4. Do **not** change CFG-02 / CFG-03 behavior until a support/transaction policy is settled; reproduce only if necessary using existing tests.
5. Keep docs English, user explanations Persian, and record actual verification precisely. No unrequested architecture growth.

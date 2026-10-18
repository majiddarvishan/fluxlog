# Current status — Fluxlog

**Date:** 2026-10-18  
**Repository:** https://github.com/majiddarvishan/fluxlog  
**Verified main HEAD before the timestamp fix:** `d4866c7a62789f31531e2017fd43bcadcdaa2757`  
**Current change:** [PR #13 — Human-readable Unix console timestamps](https://github.com/majiddarvishan/fluxlog/pull/13)  
**Adapter version:** `github.com/majiddarvishan/goconfig v1.6.0` (nested module only).  
**Go version in CI:** 1.23.x.

## Completed baseline work

- [PR #9](https://github.com/majiddarvishan/fluxlog/pull/9): English architecture/handoff documentation, merged.
- [PR #10](https://github.com/majiddarvishan/fluxlog/pull/10): upgrade optional goconfig adapter dependency to v1.6.0, merged.
- [PR #11](https://github.com/majiddarvishan/fluxlog/pull/11): behavior review and evidence-backed issues, merged.
- [PR #12](https://github.com/majiddarvishan/fluxlog/pull/12): executable Phase 1 regression scenarios, merged. Historical results: CFG-01 and LVL-01 PASS; CFG-02, CFG-03 and four FMT-01 numeric cases SKIP.

## Current authorized change: human-readable console timestamps

The maintainer requested time to remain **human-readable when printed**. This authorizes FMT-01 only. The implementation:

- Supports numeric Unix seconds, milliseconds, microseconds and nanoseconds emitted by zerolog in the human-readable console formatter.
- Converts to the process's **local timezone**, using `ConsoleConfig.TimeFormat`; default remains **RFC3339** for compatibility. A simpler 24-hour layout is `2006-01-02 15:04:05`.
- Preserves string timestamps, negative/fractional timestamp precision and **JSON output encoding unchanged**.
- Changes only `console.go`, tests, configuration documentation and `.codex`; **not** goconfig, logger lifecycle, serialization format or the dependency graph.

**Verification:** [GitHub Actions run 39147165603](https://github.com/majiddarvishan/fluxlog/actions/runs/39147165603) passed both `core` and `goconfig-adapter` jobs, including dependency consistency, format, vet, race tests and examples. All previously SKIPPED FMT-01 numeric cases now PASS; custom-layout, negative/fractional and full Logger pipeline tests also PASS. This is remote CI, not a local Go run.

## Current backlog

- **Complete:** DEP-01, DEP-02, INV-01, CFG-01, LVL-01, **FMT-01**.
- **Open and reproducible:** CFG-02 (multiple loggers bound to one goconfig path; second callback not invoked) and CFG-03 (failed FileSource persistence rolls back config but not runtime Logger level). Two known-issue test cases remain SKIPPED, **not fixed**.
- **Open, not started:** LIFE-01, FILE-01, ERR-01, FMT-02, REL-01, DOC-01, PERF-01, CI-01.

Do not add unrelated fixes or tags. Before future work, check actual `main` SHA, PR #13 merge state and CI, review the small diff, and update this status with confirmed merge SHA. Verify/deal with branch cleanup after merging.

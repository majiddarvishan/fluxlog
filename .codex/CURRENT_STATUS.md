# Current status — Fluxlog

**Updated:** 2026-10-18  
**Repository:** https://github.com/majiddarvishan/fluxlog  
**Verified main SHA before Phase 2 tests:** `7c6e796397da72b05260681cb9ba1b3cf47b4746`  
**Active work:** [PR #14 — Lifecycle, file rotation and sink errors](https://github.com/majiddarvishan/fluxlog/pull/14)  
**Pinned optional adapter:** `github.com/majiddarvishan/goconfig v1.6.0`. Core remains independent of goconfig.

## Completed work and evidence

- PR #9: initial English engineering documentation, merged.
- PR #10: goconfig v1.6.0 dependency update, merged.
- PR #11: source-backed behavior review, merged.
- PR #12: Phase 1 regression tests and global level contract, merged.
- PR #13: human-readable Unix console timestamp conversion (FMT-01), merged. Post-merge CI [run 39147439513](https://github.com/majiddarvishan/fluxlog/actions/runs/39147439513) PASS.
- PR #14: test-only Phase 2 hardening evidence. [Actions run 39158993645](https://github.com/majiddarvishan/fluxlog/actions/runs/39158993645) PASS for both core and goconfig adapter, race detector included. Phase 2 evidence: [PHASE2_TEST_EVIDENCE_2026-10-18.md](PHASE2_TEST_EVIDENCE_2026-10-18.md).

## Phase 2 status

- **LIFE-01: complete** — blocked active writes vs Close, concurrent Close/Write, post-close ErrClosedPipe, resource ownership and SetDefault swaps PASS.
- **FILE-01: complete within recorded scope** — real size rollover, on-disk backup preservation, retention MaxBackups, MaxAge deletion and gzip decompression PASS. LocalTime is tested as forwarded configuration; exact local-time filename string is not asserted.
- **ERR-01: complete** — failing and short writers, partial multi-sink propagation, and invalid log parent path PASS.

**No production logging source, dependency, CI workflow, public API or behavior changed in PR #14.** Four `*_test.go` files and documentation only. Coverage in remote Core run: 75.9% (not a performance claim).

## Outstanding tasks

- **CFG-02: OPEN:** two loggers binding the same goconfig v1.6.0 path do not both receive a Replace callback; recorded SKIP still present.
- **CFG-03: OPEN:** a failed FileSource persistence after a successful callback can leave runtime Logger level out of sync; recorded SKIP still present.
- **FMT-02: OPEN:** additional ANSI/color contract tests.
- **REL-01 / DOC-01: OPEN:** adapter tag and README release instructions, plus accepted contract documentation.
- **PERF-01 / CI-01: OPTIONAL:** only on a concrete need.

## Verification and handoff

Before considering PR #14 complete, verify final commit CI, merge result, post-merge CI and branch cleanup. If branch removal cannot be performed with the connected GitHub action, report the exact deletion command instead of claiming the branch was deleted. No new release tag has been created.

Respond to the maintainer in **Persian**; keep every committed repository file and code comment **in English**.

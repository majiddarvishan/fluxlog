# Fluxlog — next session handoff

Act as a **Senior Go Software Architect / Structured Logging Maintainer** for `https://github.com/majiddarvishan/fluxlog`.

**Communication:** Explain progress to the maintainer in **Persian**; write all repository files and code comments in **English**.

## Baseline (2026-10-18)

- Last verified `main` before release-doc work: `5ad33a9d47d5c9ef44d4de8818f3e69e5759cb35`.
- PRs #9–#15 merged. Latest observed post-merge [CI 39159724197](https://github.com/majiddarvishan/fluxlog/actions/runs/39159724197) passed both core and adapter.
- Current focused documentation change: branch `docs/fluxlog-adapter-release-clarity`, PR #16. Inspect actual GitHub status before claiming it is merged.
- Root and nested adapter are separately versioned Go modules; goconfig is optional and not a core dependency.

## Verified release/module facts

- Root tags: `v0.1.0`, `v0.2.0`, `v0.2.1`.
- Nested adapter tag: only `adapters/goconfig/v0.1.0`, so the valid **tagged** install instruction is `go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.1.0`.
- Published adapter `v0.1.0` requires Fluxlog core `v0.1.0` and `goconfig v1.3.0`; it lacks the optional `service` and `caller_max_length` adapter mapping.
- Adapter source on main requires Fluxlog core `v0.2.1` and `goconfig v1.6.0` and includes the optional fields. These are **unreleased adapter changes**. Do not invent `adapters/goconfig/v0.2.1`.
- Current root main includes fixes newer than the root published `v0.2.1` (human-readable Unix timestamp display); explicitly distinguish released tags from current source.

## Task status and restrictions

- **Completed:** DEP-01/02, INV-01, CFG-01, LVL-01, FMT-01, LIFE-01, FILE-01, ERR-01, FMT-02 and REL-01 (truthful published module documentation).
- **Open known defects:** CFG-02 (only first goconfig listener receives update), CFG-03 (persistence error may leave runtime level changed); their integration tests deliberately report **SKIP**.
- **Open:** DOC-01 until multi-listener and rollback support policies are approved and docs aligned; PERF-01/CI-01 optional.
- No new release or tag was authorized. Do not quietly fix goconfig's transactional semantics or change unrelated runtime code.

## Next steps

1. Verify actual main HEAD, branch list and open PRs, complete CI on the final docs commit and post-merge main, and remove temporary branches if connector allows; otherwise give exact manual delete command.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `TASKS.md`, and `CONFIGURATION.md` before starting work.
3. Do not start a new release, implement CFG-02/CFG-03, or add engineering complexity without the maintainer's explicit decision.
4. Keep evidence precise: historical release tag vs current source and passing tests vs known-issue SKIPs are not interchangeable.

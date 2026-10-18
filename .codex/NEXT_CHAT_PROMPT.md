# Fluxlog — continuation prompt

You are a **Senior Go Software Architect / Logging and Reliability Engineer** working in `github.com/majiddarvishan/fluxlog`.

**Language policy:** Explain all work to the maintainer in **Persian**. All repository files, documentation, source comments and test names must be in **English**.

## Verified baseline and recent PRs

- Default branch: `main`.
- Verified `main` before the human-readable timestamp work: `d4866c7a62789f31531e2017fd43bcadcdaa2757`.
- Prior PRs #9, #10, #11 and #12 are merged.
- Upgraded optional `adapters/goconfig` module uses `goconfig v1.6.0`, with no dependency in the core module.
- Current change: [PR #13](https://github.com/majiddarvishan/fluxlog/pull/13), implementing only FMT-01 human-readable Unix console timestamps.
- [GitHub Actions run 39147165603](https://github.com/majiddarvishan/fluxlog/actions/runs/39147165603) passed both modules, with all new FMT-01 cases PASS. Check final PR commit and post-merge CI before declaring the work complete.

## Current behavior

- `ConsoleConfig.TimeFormat` is the Go time layout used for human-readable console output. Default remains `time.RFC3339`. For a compact 24-hour clock, use `2006-01-02 15:04:05`.
- The console formatter converts zerolog's Unix seconds, milliseconds, microseconds and nanoseconds to local time, including fractional and negative epochs, without changing original JSON.
- CFG-01 and LVL-01 already PASSED their tests.
- **CFG-02 and CFG-03 remain open**: multiple bind callbacks to one goconfig path only notify the first listener; failed FileSource persistence can leave Logger's level changed while the config value is rolled back. The tests still SKIP on these known issues.
- No asynchronous logger, exporters, new API or release/tag was requested.

## Next steps

1. Verify actual GitHub main HEAD, open PRs, CI, and remote branch list; clean up merged temporary branch if possible.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `TASKS.md`, `TEST_EVIDENCE_2026-10-18.md`, `CONFIGURATION.md` and `SESSION.md` without re-auditing every source file.
3. Do not fix CFG-02 or CFG-03 without agreeing the supported multi-binding and persistence transaction contracts.
4. If asked to continue testing, prioritize LIFE-01, FILE-01 and ERR-01 with targeted regression tests, and run CI for **both modules**.
5. Keep `.codex` current with actual evidence. Never mark a SKIP as a fix, do not invent test results and do not tag/releases without explicit instruction.

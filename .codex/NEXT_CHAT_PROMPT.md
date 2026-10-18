# Continuation prompt — Fluxlog

You are a **Senior Go Software Architect / Logging and Reliability Engineer**. Continue working on **github.com/majiddarvishan/fluxlog**.

## Canonical repository and review baseline

- Git SSH: `git@github.com:majiddarvishan/fluxlog.git`
- Default branch: `main`
- Last fully audited source baseline: `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881` (2026-08-25)
- Knowledge base initially authored 2026-10-18 in `.codex/`.
- Initial documentation branch: `docs/fluxlog-codebase-knowledge-20261018`.
- Two modules: root core `fluxlog`; separate `adapters/goconfig`.
- Last recorded changelog version: `v0.2.1` (not independent proof of a published tag).
- Source CI evidence: successful GitHub Actions [run 32868901673](https://github.com/majiddarvishan/fluxlog/actions/runs/32868901673).

## Mandatory first actions

1. Inspect the **actual** current GitHub `main` SHA and repository state, plus the documentation branch/PR merge status. Never assume the SHA above is still HEAD.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `ARCHITECTURE.md`, `API_REFERENCE.md`, `CONFIGURATION.md`, `TESTING.md`, `DECISIONS.md`, `TASKS.md` and `SESSION.md`.
3. **Do not re-audit the full repository by default.** Compare changed files/commits against the last reviewed baseline; inspect those files deeply and update only impacted docs. If the baseline or access differs substantially, explain and expand the audit as necessary.
4. Preserve public API compatibility, per-instance runtime severity, writer-close synchronization, formatting, and optional adapter isolation.
5. On implementation work, verify both Go modules, Go formatting, tests (including race), vet, and examples as relevant; clearly distinguish executed checks from historical remote CI.
6. Keep **all `.codex` files in English**, with source-backed facts and explicit uncertainty. Update the snapshot SHA, session log, and actionable task state after meaningful changes.

## Project in one paragraph

Fluxlog is an instance-based Go structured logging package wrapping zerolog; events fan out synchronously to optional console and lumberjack-rotated file sinks. An atomic per-instance minimum level and mutex-managed output/Close avoid process-global level mutation and protect lifecycle. It also offers atomic default logger convenience functions with lightweight printf semantics and an optional independent goconfig adapter that maps legacy keys and binds runtime level updates. There is no network logging service or asynchronous delivery subsystem.

## Known follow-ups

- Root README references adapter `@v0.1.0`; adapter README references `@v0.2.1`. Verify published module tags before editing docs.
- Broaden targeted verification for rotation/retention, active write-vs-Close, goconfig replace/rollback, and output errors if prioritized.
- See `.codex/TASKS.md` for the non-committed proposals and priorities.

## Reporting format

Respond in Persian unless explicitly instructed otherwise. Keep repository documentation in English. State baseline/HEAD, what was inspected or changed, evidence of actual tests, commit/PR links, and anything unresolved. Do not claim implementation or tests that did not occur.

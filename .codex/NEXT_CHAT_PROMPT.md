# Fluxlog — next chat handoff

Act as a **Senior Go Software Architect / Structured Logging Reliability Engineer**. Continue the project at `git@github.com:majiddarvishan/fluxlog.git`.

## Grounded state (2026-10-18)

- Default branch: `main`.
- Verified main HEAD **before** dependency work: `93e0945e14a34a72ab93d8d048cd1131c523c6b6`.
- Core runtime audit baseline: `f0b79fcd7afeb0406e75a1c2ca5fe1e15243e881`; documentation PR #9 merged without changing runtime code.
- Active change branch: `chore/goconfig-1.6.0-and-hardening-plan`, PR [#10](https://github.com/majiddarvishan/fluxlog/pull/10).
- **Maintainer choice:** Use `github.com/majiddarvishan/goconfig v1.6.0` in the **nested adapter module**. **Other improvements are to be planned only** until a new implementation instruction.
- Two Go modules: core root, separate `adapters/goconfig`. Core must never gain goconfig as a dependency.
- Latest recorded core changelog entry is v0.2.1. The adapter has a verified Git tag `adapters/goconfig/v0.1.0` but not `adapters/goconfig/v0.2.1`.

## Required first steps

1. Inspect actual `main` SHA, PR #10 commit/CI state and remote branch status. Do not assume the PR has or has not merged.
2. Read `.codex/START_HERE.md`, `CURRENT_STATUS.md`, `TASKS.md`, `IMPLEMENTATION_PLAN.md`, and `SESSION.md`; use architecture/configuration docs only where relevant.
3. If dependency upgrade is unverified, prioritize exact Go-generated `adapters/goconfig/go.mod` and `go.sum`, then complete adapter/core CI. Never invent checksum lines.
4. Other tasks remain **planning only**: live replace and rollback, multi-bind callback behavior, external zerolog global level, numeric timestamps, Close/Write, rotation, output errors, release docs and optional benchmarks.
5. After any PR merge, verify main HEAD and CI, update status/handoff, and delete remote temporary branches when supported; if not, report the manual cleanup command.
6. Do not re-audit all source files: inspect the diff from the recorded baseline and narrowly relevant code.
7. Keep repository docs and plans **in English**. Report only tests actually executed and changes truly merged.

## Core architecture (short reference)

Fluxlog is a synchronous, structured Go logger built on zerolog. Logger-specific atomic minimum-level filtering sends events to console and optional lumberjack-rotated file sinks with mutex-coordinated writer lifecycle; package functions use an atomic default pointer. The optional goconfig adapter parses a legacy flat config object and binds the manager's `level` node to runtime level changes.

## Delivery standard

For each task, show evidence, actual commit/PR, tested commands or CI links, remaining blockers and an accurate state (planned, implemented, verified, merged). Never claim published adapter tags based on changelog text alone.

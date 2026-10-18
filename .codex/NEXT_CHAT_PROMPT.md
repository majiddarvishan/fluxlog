# Fluxlog — next chat handoff

Use Persian to explain progress; commit all code and repository documents in English.

## Baseline and current work

- Fluxlog source: `github.com/majiddarvishan/fluxlog`; baseline `main@3c3b38352e44eb311f7841f670c8071ae7cf41a7`.
- Working branch: `feature/fluxlog-postcommit-level-binding`. Verify actual PR and main state before acting.
- Upstream `goconfig` merged `OnReplaceCommitted` at `20a58eb9a426d0ed7d522b1e5848a0971ce58198`. No upstream tag newer than `v1.6.0` was verified.
- Adapter development module pins that upstream commit via Go canonical pseudo-version `v1.6.1-0.20261018202246-20a58eb9a426`; checksum generated from real Go tooling in checksum-only Actions run `39165338586`. This is an unreleased dependency: do not describe it as tagged `v1.7.0`.
- `BindLevel` now validates before and applies Logger level only after successful commit, enabling multiple subscribed loggers and avoiding FileSource side-effect drift. Strict no-SKIP CFG-02/CFG-03 tests cover Replace, Batch and failure paths.
- Root Go module and unrelated production code unchanged. No release or tag requested.

## Verification and release gate

1. Verify the current branch's final consolidated CI: `gofmt`, dependency tidy/sum, `go vet`, `go test -race`, example builds for both modules. Do **not** assume it succeeded.
2. If the CI passes, mark UP-02/CFG-02/CFG-03 done and record the exact CI evidence; merge when appropriate, then verify post-merge CI and branch cleanup.
3. If it fails, inspect exact failing step and correct only that issue.
4. A final upstream goconfig `v1.7.0` and a new nested Fluxlog adapter tag are **not** authorized. Ask explicitly before creating tags. A future tagged migration can replace the pinned pseudo-version without adding a new feature.
5. Read `.codex/CFG02_CFG03_INTEGRATION.md`, `TASKS.md`, `CURRENT_STATUS.md`, and `SESSION.md` as needed. Do not re-audit unrelated core code.

## Known limitation

`Manager.OnReplace` executes the first pre-commit callback registered on a level path. The new Fluxlog adapter's validator must own that first handler; shared validation with externally registered handlers is not generally guaranteed. The post-commit API broadcasts to all subscribed observers. There is no unsubscribe API; callbacks remain retained until the Manager is disposed.

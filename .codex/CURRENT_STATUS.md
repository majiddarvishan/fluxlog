# Current status — Fluxlog

**Updated:** 2026-10-19  
**Verified main baseline:** `3c3b38352e44eb311f7841f670c8071ae7cf41a7`  
**Working branch:** `feature/fluxlog-postcommit-level-binding`  
**Scope:** Implement and test CFG-02/CFG-03 and UP-02 without creating any release/tag.

## Implemented, pending final CI

- The nested optional goconfig adapter now targets the exact upstream commit `20a58eb9a426d0ed7d522b1e5848a0971ce58198` through the canonical `v1.6.1-0.20261018202246-20a58eb9a426` Go pseudo-version.
- Go-generated checksums were obtained from GitHub Actions [module resolution run 39165338586](https://github.com/majiddarvishan/fluxlog/actions/runs/39165338586); no checksum was guessed. This is **not a tagged upstream release**.
- `BindLevel` now separates initial application, pure pre-commit validation and post-persistence observer updates. All loggers bound to one node receive its committed level.
- Existing CFG-02 and CFG-03 tests no longer use `t.Skipf` to conceal known deviations. Additional tests exercise successful final-value Batch updates, dry runs, invalid levels and FileSource batch write failure.
- Root Fluxlog `go.mod` remains independent of goconfig. The published adapter `v0.1.0` remains untouched. No tags/releases authorized or created.

## Final verification gate

**Not yet executed on this full integration:** the consolidated final GitHub Actions CI (both modules: format, `go mod tidy`, `go vet`, `go test -race`, example builds). Do not claim CFG-02/CFG-03 or UP-02 complete until strict tests genuinely pass. If the final CI is green, update the task state, record run ID, merge with reviewable scope and verify the post-merge main. If not, preserve the PR unmerged and record the exact failure.

## Supported contract limits

Legacy goconfig `Manager.OnReplace` selects the first registered pre-commit callback per path. Fluxlog must own the first validator registration for a managed level path or mixed external handlers may bypass validation. The new `OnReplaceCommitted` callbacks cannot return errors and remain registered for the Manager lifetime. Refer to [CFG02_CFG03_INTEGRATION.md](CFG02_CFG03_INTEGRATION.md).

**Open beyond this work:** final named goconfig release/tag, future optional nested adapter tag and behavior docs after passing tests. Optional PERF-01/CI-01 remain deferred. All repository files English, conversational output Persian.

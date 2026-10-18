# Current status — Fluxlog

**Updated:** 2026-10-18  
**Repository:** https://github.com/majiddarvishan/fluxlog  
**Verified main baseline before release-doc corrections:** `5ad33a9d47d5c9ef44d4de8818f3e69e5759cb35`  
**Current work:** Release/module install documentation, PR #16 (no Go implementation changes).  
**Go modules:** root `github.com/majiddarvishan/fluxlog`; independent optional `github.com/majiddarvishan/fluxlog/adapters/goconfig`.

## Completed development and testing

- PRs #9–#15 were merged. Phase 1/2 work includes human-readable Unix timestamp rendering, lifetime, rotation, output error and ANSI/color tests.
- [Post-merge CI run 39159724197](https://github.com/majiddarvishan/fluxlog/actions/runs/39159724197) passed core and adapter jobs on the baseline above.
- **REL-01:** tag and module dependency state verified, both README install references aligned to the published adapter tag, and published-vs-current-source differences documented. Verify this PR's final merge and post-merge CI before treating the correction as live on `main`.

## Verified published versions

| Module | Existing published tag(s) | Tag/module dependency facts |
| --- | --- | --- |
| Root Fluxlog | `v0.1.0`, `v0.2.0`, `v0.2.1` | Latest verified root release `v0.2.1`, older than newer main-branch fixes |
| Adapter | `adapters/goconfig/v0.1.0` only | Published `v0.1.0` requires core Fluxlog `v0.1.0` and goconfig `v1.3.0`; optional `service` and `caller_max_length` mappings not yet included |
| Adapter on `main` | No newer tag | Source `go.mod` requires core `v0.2.1` and goconfig `v1.6.0`; source includes extra optional fields |

**Important:** a root `v0.2.1` tag is not `adapters/goconfig/v0.2.1`. The documented install command `go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.1.0` selects the tagged adapter, **not** the latest unreleased source changes. No release/tag was created in this iteration.

## Open technical issues

- **CFG-02: OPEN**, observed SKIP: goconfig v1.6.0 invokes the first registered OnReplace handler for a path; multiple loggers do not all update.
- **CFG-03: OPEN**, observed SKIP: after a successful runtime level callback, failed FileSource persistence can leave Logger state out of sync with the restored node.
- **DOC-01: PARTIAL/OPEN:** README/schema/tag clarity corrected; final supported multi-bind and rollback contracts need maintainer decisions and potential code fixes.
- **PERF-01 / CI-01:** optional and not planned without a clear need.

## Workflow notes

Only READMEs and English `.codex` files were changed. Do not edit runtime Go code or create tags for release documentation work. Actual Go CI evidence must be checked on the final PR commit and post-merge main commit. After merge, verify and clean up the temporary branch.

Explain progress to the user **in Persian**; all repository files and comments must remain **English**.

# Current status — Fluxlog

**Verified:** 2026-10-19
**Main release baseline:** `0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e`.
**Lockstep v0.3.0 publication:** complete for core and goconfig adapter.

## Published Go modules

| Module | Version | Annotated repository tag | Peeled commit |
| --- | --- | --- | --- |
| Core | `v0.3.0` | `v0.3.0` | `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd` |
| goconfig adapter | `v0.3.0` | `adapters/goconfig/v0.3.0` | `0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e` |

- Core tag object: `6ad7581ac7fe611689b6087a183d10f6f4ba06b6`. Core source passed [CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174).
- Adapter tag object: `0cf414a7adf067db8c870f6a30f32030d7240622`. [PR #21](https://github.com/majiddarvishan/fluxlog/pull/21) passed [CI 39170725761](https://github.com/majiddarvishan/fluxlog/actions/runs/39170725761), merged, and both jobs passed [post-merge main CI 39170859081](https://github.com/majiddarvishan/fluxlog/actions/runs/39170859081), including strict CFG-02/CFG-03 tests without SKIP.
- The maintainer verified the published nested module resolves using `GOPROXY=https://proxy.golang.org go list -m github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0`.
- Upstream `goconfig v1.7.0` is published on `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`; [CI 39167180503](https://github.com/majiddarvishan/goconfig/actions/runs/39167180503) passed.

## Architectural and release guarantees

- Core and all official adapters share release version numbers from v0.3.0 onward; Go modules and tags remain separate. Release commits can differ.
- The adapter requires core v0.3.0 and upstream goconfig v1.7.0. The core module has **no** goconfig dependency.
- Historical releases `v0.2.1` and `adapters/goconfig/v0.1.0` remain unchanged. No adapter v0.2.0 tag was published.
- The first legacy goconfig `OnReplace` handler validates the initial candidate; committed observers live for the Manager lifetime. This limitation is documented and accepted for v0.3.0.
- No standalone GitHub Release object was found during this verification. Annotated Git tags and Go module availability constitute the publication evidence.

**Next:** normal maintenance; optional PERF-01 and CI-01 remain backlog only.
See [RELEASING.md](../RELEASING.md), [TASKS.md](TASKS.md) and
[RELEASE_READINESS_2026-10-19.md](RELEASE_READINESS_2026-10-19.md).

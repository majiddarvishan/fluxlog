# Continue Fluxlog lockstep v0.3.0 release

Act as Senior Go Software Architect / Go Module and Release Engineer. Explain in Persian; all repository code, comments and documentation must be English.

## Verified baseline on 2026-10-19

- Fluxlog `main@fe57c3e1271d5ada9339b255ac4ca0d727aef5bd` passed both jobs in [post-merge CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174).
- Core `v0.3.0` is an **existing annotated tag** on that precise tested main SHA. Never recreate or move it.
- Upstream `goconfig v1.7.0` is an existing annotated tag on `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`. Upstream CI passed.
- Official upstream dependency [Fluxlog PR #20](https://github.com/majiddarvishan/fluxlog/pull/20) merged. Strict CFG-02/CFG-03 passed without SKIP.
- Adapter `adapters/goconfig/v0.1.0` was the latest published adapter tag at this checkpoint.
- Working branch: `release/fluxlog-lockstep-v0.3.0`. Recheck its latest HEAD, PR and CI before resuming.

## Approved lockstep policy

From `v0.3.0` onwards core and all official adapters use matching release version numbers, but distinct Go modules and tag prefixes. Adapter `v0.3.0` must depend on released core `v0.3.0`. Release tag commits need not match. The root module never depends on goconfig.

## Next steps

1. Inspect the current release PR and adapter `go.mod/go.sum`; resolve Go-generated checksums, not invented values. Local networking was unavailable when preparing the release; use verified Go CI artifacts if necessary.
2. Ensure both modules pass tidy, formatting, vet, race tests, example builds, and strict CFG-02/CFG-03 on the final PR commit.
3. Merge only after green CI; confirm post-merge `main` CI also passes both jobs.
4. Publish annotated `adapters/goconfig/v0.3.0` on the CI-tested post-merge commit using an authorized Git client if connector tagging is unavailable; do not use v0.2.0 and never move tags.
5. Update final release status, README install commands and changelogs only according to actual published versions and Go module resolution.

See `RELEASING.md`, `.codex/CURRENT_STATUS.md`, `.codex/TASKS.md` and `.codex/RELEASE_READINESS_2026-10-19.md`.

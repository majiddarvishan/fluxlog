# Current status — Fluxlog

**Verified on 2026-10-19:** `main@fe57c3e1271d5ada9339b255ac4ca0d727aef5bd`; both core and goconfig-adapter jobs passed [post-merge CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174).

Core `v0.3.0` is an existing annotated release tag on this exact SHA. The official upstream goconfig `v1.7.0` tag points to `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc` and passed its own CI. [PR #20](https://github.com/majiddarvishan/fluxlog/pull/20) is merged; strict CFG-02/CFG-03 integration tests are no longer skipped.

## Lockstep policy

Beginning at `v0.3.0`, core and official adapters use the same version number while retaining independent Go modules. The nested tag is `adapters/goconfig/v0.3.0`. The adapter must use core `v0.3.0`, and core must not depend on goconfig. Tag commits can differ.

## Remaining release gates

- At the checkpoint, the latest adapter tag is still `adapters/goconfig/v0.1.0`.
- On `release/fluxlog-lockstep-v0.3.0`, regenerate the adapter core dependency and sums using Go tooling; run complete checks on both modules and confirm CFG-02/CFG-03.
- Merge the tested PR, verify post-merge CI, then publish `adapters/goconfig/v0.3.0` with an authorized Git client; never rewrite any tag.
- Core proxy resolution and nested adapter proxy resolution must be verified rather than presumed.
- goconfig's first legacy `OnReplace` handler validates the initial candidate; committed listeners remain for the lifetime of the Manager. This documented limitation is not a v0.3.0 blocker.

See [RELEASING.md](../RELEASING.md), [TASKS.md](TASKS.md), and [RELEASE_READINESS_2026-10-19.md](RELEASE_READINESS_2026-10-19.md).

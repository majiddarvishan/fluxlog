# Current status — Fluxlog

**Updated:** 2026-10-19
**Main before dependency migration:** `70b17302e5cc5b2974fff502227552ed8cf4d694`.
**Migration branch:** `release/fluxlog-goconfig-v1.7.0-migration`.
**Scope:** official upstream goconfig module upgrade and release documentation.
No root Go implementation changes.

## Verified baseline

- Previous fixes CFG-02 and CFG-03 (shared Logger updates and persistence
  rollback consistency) passed without SKIP in PR #18 and post-merge
  CI [39166047372](https://github.com/majiddarvishan/fluxlog/actions/runs/39166047372).
- Previous main release documentation passed both jobs in
  [CI 39167270820](https://github.com/majiddarvishan/fluxlog/actions/runs/39167270820).
- Upstream goconfig v1.7.0 is an **annotated tag** on commit
  `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`.
- Go resolved the official upstream tag and generated hashes in
  [CI 39167874434](https://github.com/majiddarvishan/fluxlog/actions/runs/39167874434).
  The development pseudo-version is replaced by `v1.7.0` in the **nested
  adapter module only**.

## Remaining release gates

- Verify final migration PR and post-merge CI for both modules, then finish
  documentation with the exact main SHA and CI run.
- Root Fluxlog `v0.3.0` and nested adapter `adapters/goconfig/v0.2.0`
  are **planned, not published**. The root tag requires approval and
  access to an authorized Git client; after it exists, update the nested
  adapter's root Fluxlog dependency and verify before tagging the adapter.
- Old published tags `v0.2.1` and `adapters/goconfig/v0.1.0`
  remain available. Do not misrepresent them as including current source.
- Legacy goconfig `OnReplace` uses the first registered validation
  handler; users must register Fluxlog's validator first on a level path.
  Post-commit listeners live until their Manager is discarded.

See [RELEASE_READINESS_2026-10-19.md](RELEASE_READINESS_2026-10-19.md)
and [RELEASING.md](../RELEASING.md) for exact gates and commands.

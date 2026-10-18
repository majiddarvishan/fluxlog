# Fluxlog — independent release readiness (2026-10-19)

**Release preparation, not publication.** The technical fixes CFG-02 and
CFG-03 passed on Fluxlog `main@5de1e20644c56d0db3163ffcfd59a1d0b7937539`;
[post-merge CI 39166047372](https://github.com/majiddarvishan/fluxlog/actions/runs/39166047372)
passed both root and adapter modules.

The upstream `goconfig` project has the new post-commit API merged on main
but lacks a new tagged release. An upstream v1.7.0 release-readiness PR
adds permanent CI, an English release-note draft and a safe tagging procedure.

## Release sequence and current blockers

1. **goconfig v1.7.0:** wait for successful release-readiness PR + main CI,
   then create the authorized annotated tag. The GitHub connector does not
   expose direct tag/release creation, so an authorized Git client/operator
   must publish it. Verify the actual remote ref and Go module resolution.
2. **Fluxlog adapter dependency:** only after the upstream tag exists,
   replace the source-commit pseudo-version with `v1.7.0` using
   `go get`/`go mod tidy`; no manual hash fabrication. Run both module
   checks, and merge the verified change.
3. **Core and adapter independent releases:** propose `v0.3.0` for the root
   and `adapters/goconfig/v0.2.0` for the nested module, with separate
   tagging gates and a root dependency update for the adapter after root
   publication. These are **proposed versions, not created tags**.

Precise commands and gates: [`RELEASING.md`](../RELEASING.md).
Root change history: [`CHANGELOG.md`](../CHANGELOG.md).
Adapter change history: [`adapters/goconfig/CHANGELOG.md`](../adapters/goconfig/CHANGELOG.md).

Avoid claiming the pseudo-version is an official release. Do not change root
`go.mod` to depend on goconfig. Maintain the legacy OnReplace first-handler
registration caveat and retained listener lifetime in release notes.

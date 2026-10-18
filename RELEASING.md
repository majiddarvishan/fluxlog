# Fluxlog independent release procedure

Fluxlog contains two independently versioned Go modules:

| Module | Module path | Suggested next tag (not published) |
| --- | --- | --- |
| Core | `github.com/majiddarvishan/fluxlog` | `v0.3.0` |
| Optional goconfig adapter | `github.com/majiddarvishan/fluxlog/adapters/goconfig` | `adapters/goconfig/v0.2.0` |

**Verified 2026-10-19:** upstream goconfig `v1.7.0` exists as an
annotated tag on commit `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`.
The tagged Go module was resolved with verified checksums in
[Actions run 39167874434](https://github.com/majiddarvishan/fluxlog/actions/runs/39167874434).
The current core tag is `v0.2.1`; the existing nested adapter tag is
`adapters/goconfig/v0.1.0`. Proposed new Fluxlog tags have **not** been created.

## Release stages and gates

1. **DONE — Publish upstream goconfig.** The official upstream `v1.7.0`
   annotated tag exists at the CI-verified main commit. Its Go module
   resolves against GitHub and its checksums were generated with Go.
2. **DONE in the migration PR — Pin the upstream tag.** The nested
   `adapters/goconfig/go.mod` now requires `goconfig v1.7.0`, and its
   `go.sum` contains the Go-generated checksums. Verify the PR and
   post-merge CI for both modules before treating this step as merged.
   The root `go.mod` must not gain a goconfig dependency.
3. **NEXT — Tag Fluxlog core.** After both root and adapter CI jobs succeed
   on the exact `main` commit, create and push an annotated root
   `v0.3.0` tag. Confirm the tag's peeled commit equals the tested main SHA
   and the Go module resolves:

   ```bash
   GOPROXY=https://proxy.golang.org go list -m github.com/majiddarvishan/fluxlog@v0.3.0
   ```

4. **NEXT — Update the nested adapter core dependency.** After the root tag
   is visible to Go, change the **nested** `adapters/goconfig/go.mod` from
   core Fluxlog `v0.2.1` to `v0.3.0` using:

   ```bash
   cd adapters/goconfig
   go get github.com/majiddarvishan/fluxlog@v0.3.0
   go mod tidy
   cd ../..
   ```

   Regenerate `go.sum` using Go, verify both modules including strict
   CFG-02/CFG-03 tests, and merge the exact tested change.
5. **NEXT — Tag the nested adapter.** Create and push a separate annotated
   `adapters/goconfig/v0.2.0` tag **at the post-migration tested main
   commit**, never at the root `v0.3.0` commit unless they are identical
   and the dependency gate is satisfied. Verify that the Go module
   `github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.2.0`
   resolves against the pushed tag.
6. **Publish notes.** The root release uses `CHANGELOG.md`; the nested
   adapter release uses `adapters/goconfig/CHANGELOG.md`. Do not claim a
   GitHub Release has been published unless a release object exists.

A source commit that has not passed both CI jobs cannot be tagged by this
procedure. Never overwrite an existing release tag.

## Annotated tag commands

Execute only after the applicable gate is verified. Replace `TAG` with
the exact approved release tag, one release at a time.

```bash
git fetch origin --tags
git switch main
git pull --ff-only origin main
git status --short
git log -1 --format='%H %s'
git ls-remote --tags origin "refs/tags/$TAG"
# Confirm the exact HEAD passed CI and the tag does not exist.
git tag -a "$TAG" -m "Release $TAG"
git push origin "refs/tags/$TAG"
git ls-remote --tags origin "refs/tags/$TAG" "refs/tags/$TAG^{}"
```

The GitHub connector may lack tag/release creation actions. If it does,
run this through an authorized Git client or an approved release UI and
verify the resulting remote tag. **Never force-push or silently move tags.**

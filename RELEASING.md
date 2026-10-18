# Fluxlog independent release procedure

Fluxlog contains two independently versioned Go modules:

| Module | Module path | Suggested next tag (not published) |
| --- | --- | --- |
| Core | `github.com/majiddarvishan/fluxlog` | `v0.3.0` |
| Optional goconfig adapter | `github.com/majiddarvishan/fluxlog/adapters/goconfig` | `adapters/goconfig/v0.2.0` |

**As of 2026-10-19**, the existing root tag is `v0.2.1`, the adapter tag is
`adapters/goconfig/v0.1.0`, and upstream goconfig's latest verified tag is
`v1.6.0`. No suggested next tag is currently published.

## Order and gates

1. **Upstream goconfig first.** Merge and verify its v1.7.0 readiness PR,
   then have an authorized maintainer create the annotated `v1.7.0` tag
   on the exact successful main commit. Verify the tag exists and
   `go list -m github.com/majiddarvishan/goconfig@v1.7.0` resolves.
2. **Adapter dependency migration, before any Fluxlog release.** Change only
   `adapters/goconfig/go.mod` and the corresponding Go-generated
   `go.sum` from the pinned upstream pseudo-version to `v1.7.0`:

   ```bash
   cd adapters/goconfig
   go get github.com/majiddarvishan/goconfig@v1.7.0
   go mod tidy
   cd ../..
   ```

   Do **not** commit fabricated checksums or alter the root `go.mod`.
   Check that the module metadata diff contains only intended changes.
3. **Verify both modules** on the dependency-migration PR and after merging
   to main. Required checks: formatting, tidy consistency, `go vet ./...`,
   `go test -race -count=1 ./...`, and example builds in each Go module.
   The CFG-02/CFG-03 strict tests must PASS with no SKIP.
4. **Core release.** Once the final main commit is clean and both module
   jobs succeed, create and push the root `v0.3.0` annotated tag on that
   exact commit. Confirm its GitHub ref and Go-module resolution.
5. **Adapter release, separately.** Update the nested adapter module to
   depend on `github.com/majiddarvishan/fluxlog@v0.3.0` through Go,
   regenerate `go.sum`, and again verify both modules in CI. Merge first.
   Then create and push the **nested tag**
   `adapters/goconfig/v0.2.0` at the tested main commit.
6. Publish separate notes using the root `CHANGELOG.md` and the nested
   `adapters/goconfig/CHANGELOG.md`. Verify each module can be resolved.

No tag is created by this document. Tags must not be created before their
respective CI and module-version gates pass.

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

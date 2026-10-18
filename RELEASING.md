# Fluxlog lockstep release procedure

Starting with `v0.3.0`, Fluxlog core and every official adapter release use the **same version number** even if a module has no code changes. They remain separate Go modules; core never depends on goconfig. Distinct tag commits are allowed.

| Component | Module path | Git tag |
| --- | --- | --- |
| Core | `github.com/majiddarvishan/fluxlog` | `v0.3.0` |
| goconfig adapter | `github.com/majiddarvishan/fluxlog/adapters/goconfig` | `adapters/goconfig/v0.3.0` |

## Verified publication state — 2026-10-19

- Annotated core `v0.3.0` already exists and points to `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd`; both core and adapter CI jobs passed in [39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174). **Do not recreate or move this tag.**
- Upstream `goconfig v1.7.0` is published at `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`.
- Previous adapter tag: `adapters/goconfig/v0.1.0`. The new `adapters/goconfig/v0.3.0` tag has not yet been published at this checkpoint. Skip adapter v0.2.0.

## Release gates

1. Verify the core tag and module resolution with `GOPROXY=https://proxy.golang.org go list -m github.com/majiddarvishan/fluxlog@v0.3.0`.
2. In `adapters/goconfig` run `go get github.com/majiddarvishan/fluxlog@v0.3.0` and `go mod tidy`. Only commit Go-generated checksums. Keep `goconfig v1.7.0` and leave the root `go.mod` untouched.
3. Verify **both modules**: clean `go mod tidy`, gofmt, `go vet ./...`, `go test -race ./...`, and `go build ./examples/...`. CFG-02/CFG-03 must PASS without SKIP.
4. Merge the PR only after CI passes on its final HEAD, then verify both jobs pass on the merged `main` commit.
5. Using an authorized Git client, create the annotated nested tag on that *exact tested post-merge commit*. Verify its peeled target and `GOPROXY=https://proxy.golang.org go list -m github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0`.
6. Reconcile READMEs, changelogs and `.codex` with *actual* tags. A GitHub Release object is optional and is not the same as a Git tag.

## Safe adapter tagging (only after passing post-merge CI)

```bash
git fetch origin --tags
git switch main
git pull --ff-only origin main
test -z "$(git status --porcelain)" || { echo "Dirty worktree"; exit 1; }
git rev-parse HEAD
TAG=adapters/goconfig/v0.3.0
EXPECTED_SHA=<CI_VERIFIED_POST_MERGE_MAIN_SHA>
test "$(git rev-parse HEAD)" = "$EXPECTED_SHA" || { echo "Unexpected HEAD"; exit 1; }
test -z "$(git ls-remote --tags origin "refs/tags/$TAG")" || { echo "Tag exists"; exit 1; }
git tag -a "$TAG" "$EXPECTED_SHA" -m "Release Fluxlog goconfig adapter v0.3.0"
git push origin "refs/tags/$TAG"
git ls-remote --tags origin "refs/tags/$TAG" "refs/tags/$TAG^{}"
GOPROXY=https://proxy.golang.org go list -m github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0
```

Never force-push a release tag, even if core and adapter tags point to different commits.

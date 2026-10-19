# Fluxlog lockstep release procedure

Starting with `v0.3.0`, Fluxlog core and every official adapter use **matching
release version numbers**, even if one module has no source changes. Core and
adapters remain independent Go modules. Core never imports or requires goconfig.
Their release tags need not target the same Git commit.

| Component | Module | Repository tag |
| --- | --- | --- |
| Core | `github.com/majiddarvishan/fluxlog` | `v0.3.2` |
| goconfig adapter | `github.com/majiddarvishan/fluxlog/adapters/goconfig` | `adapters/goconfig/v0.3.1` (v0.3.2 pending) |

These are the latest Git tags verified on 2026-10-20; query current Git
tags and Go module versions for all future releases and integrations.

## v0.3.2 release preparation — 2026-10-20

- **Core release candidate:** documentation-only patch since v0.3.1.
  The dynamic latest-tag integration prompt was merged in
  [PR #27](https://github.com/majiddarvishan/fluxlog/pull/27)
  and passed [post-merge CI 39326544866](https://github.com/majiddarvishan/fluxlog/actions/runs/39326544866).
  This preparation commit adds the v0.3.2 changelog; **do not tag until
  the final preparation merge commit passes both main CI jobs**.
- **Core tag published:** annotated `v0.3.2` tag object
  `163e1774aa4635454d0f52186318eec85155ea67` peels to
  `9050def72d3ad4195f3c5044d1ba4dd0a077f166`.
  [Post-merge CI 39327982376](https://github.com/majiddarvishan/fluxlog/actions/runs/39327982376)
  passed both module jobs. Historical `v0.3.1` is immutable.
- **Adapter tag pending:** the nested module is being updated to require
  Core `v0.3.2`, while retaining upstream goconfig `v1.7.0`.
  Regenerate `go.sum` using Go tooling and verify both CI jobs. Only
  after the final adapter commit passes post-merge CI, publish an annotated
  `adapters/goconfig/v0.3.2` tag on exactly that commit.
- **Publication status:** root tag verified in GitHub; nested tag is not
  yet published. Confirm module proxy resolution separately.
  Never move existing tags.

## v0.3.1 publication status — 2026-10-19

- **Core published:** `v0.3.1` is an annotated tag (object
  `4ccbe38d44cb97930caf80eb593edd0d16dd6070`) which peels to
  `fdb55d6d48ed9a8fc24b8c1e00994278b7fd1785`.
  [Core source CI 39176344867](https://github.com/majiddarvishan/fluxlog/actions/runs/39176344867)
  passed both Go module jobs. No core code changes beyond v0.3.0; the
  patch release captures examples and public agent integration docs.
- **Adapter published:** annotated `adapters/goconfig/v0.3.1` tag
  (object `e23ecdd7c0e00e6b2d0555f205e627da5fbde11a`) peels to
  `8aa21dd639b4e40afdeb750f4be90a8558e88995`.
  [Post-merge CI 39177136303](https://github.com/majiddarvishan/fluxlog/actions/runs/39177136303)
  passed both module jobs. It requires core `v0.3.1` and upstream
  goconfig `v1.7.0` in the published commit.
- **Go module proxy:** GitHub tags are verified independently. Public
  proxy resolution of the adapter `v0.3.1` has not been confirmed in
  this session; consumers should validate with `go list -m`.
- **Release integrity:** historical root and adapter tags remain immutable.

## Verified v0.3.0 publication — 2026-10-19

- **Core published:** annotated tag `v0.3.0` (tag object `6ad7581ac7fe611689b6087a183d10f6f4ba06b6`) peels to `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd`. [Core release source CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174) passed both jobs.
- **Adapter published:** annotated tag `adapters/goconfig/v0.3.0` (tag object `0cf414a7adf067db8c870f6a30f32030d7240622`) peels to `0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e`. [PR #21](https://github.com/majiddarvishan/fluxlog/pull/21) merged; [PR CI 39170725761](https://github.com/majiddarvishan/fluxlog/actions/runs/39170725761) and [post-merge CI 39170859081](https://github.com/majiddarvishan/fluxlog/actions/runs/39170859081) passed both jobs.
- **Upstream goconfig:** annotated `v1.7.0` tag at `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`. The v0.3.0 adapter depends on upstream goconfig v1.7.0 and core v0.3.0. The published core module has no goconfig dependency.
- **Go module resolution:** the maintainer verified `go list -m github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0` using `proxy.golang.org`. Do not infer that a separate GitHub Release object exists; none was found during the publication verification.
- **Historical tags:** `v0.2.1` and `adapters/goconfig/v0.1.0` remain unchanged. No adapter v0.2.0 tag was created.

Install either published module:

```bash
go get github.com/majiddarvishan/fluxlog@v0.3.0
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0
```

## Procedure for future versions

1. Select one SemVer number for core **and all official adapters**, and document the planned tags (`vX.Y.Z` for the root module; `adapters/goconfig/vX.Y.Z` for the nested module).
2. Verify the exact core release commit and successful core and adapter CI (tidy, gofmt, vet, race tests and example builds). Tag core with a new annotated root tag **only if absent**. Confirm the peeled commit and Go module resolution.
3. In each nested module run `go get github.com/majiddarvishan/fluxlog@vX.Y.Z`, then `go mod tidy`. Never fabricate checksums. Retain module independence and ensure upstream requirements remain deliberate.
4. Verify clean tidy diffs and both CI jobs, including strict CFG-02/CFG-03 for the goconfig adapter. Merge only after final PR-head CI succeeds.
5. Verify post-merge main CI on the exact commit to tag. Create an annotated nested tag at that commit, checking local and remote refs for prior existence. Confirm its peeled target and Go module resolution.
6. Update changelogs, README install instructions and `.codex` from actual published refs and evidence. GitHub Release objects are optional and distinct from Git tags.

**Never recreate, force-push, or move a published tag.** The published v0.3.0
tags point to different CI-tested commits by design.

# Fluxlog v0.3.0 lockstep release readiness — 2026-10-19

## Confirmed evidence

- Fluxlog `main`: `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd`.
- Core annotated `v0.3.0` tag object `6ad7581ac7fe611689b6087a183d10f6f4ba06b6` peels to this commit.
- Both jobs passed [post-merge CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174), including tidy, formatting, vet, race and examples.
- [PR #20](https://github.com/majiddarvishan/fluxlog/pull/20) merged. CFG-02/CFG-03 passed without SKIP.
- Upstream `goconfig v1.7.0` annotated tag object `f3263dd3160cf226ec0f3371d7753615862d0c32` points to `213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`. [Upstream CI 39167180503](https://github.com/majiddarvishan/goconfig/actions/runs/39167180503) passed.

## Remaining work

Core and official adapters follow **lockstep version numbers** from `v0.3.0`, remaining separate modules with distinct Go-compliant tag names and potentially different release commits.

1. Verify public Go module resolution for already-published core `v0.3.0`.
2. Upgrade the nested adapter's core dependency from `v0.2.1` to `v0.3.0`, regenerate sums with Go, and retain upstream `goconfig v1.7.0`.
3. Pass full final PR CI for both modules including CFG-02/CFG-03; merge and verify post-merge `main` CI.
4. Publish annotated `adapters/goconfig/v0.3.0` at that tested post-merge commit. This tag was absent at the checkpoint. Existing adapter `v0.1.0` remains immutable.
5. Update published-version documentation only after actual tag verification.

The first registered legacy `OnReplace` validation handler caveat and Manager-lifetime post-commit observer retention remain supported documented limitations.

# Fluxlog v0.3.0 lockstep release — final evidence (2026-10-19)

## Publication completed

| Module | Repository tag | Annotated tag object | Peeled tested commit |
| --- | --- | --- | --- |
| Fluxlog core | `v0.3.0` | `6ad7581ac7fe611689b6087a183d10f6f4ba06b6` | `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd` |
| goconfig adapter | `adapters/goconfig/v0.3.0` | `0cf414a7adf067db8c870f6a30f32030d7240622` | `0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e` |

Both tags are annotated; their peeled commits intentionally differ.
No existing tag was moved or force-pushed.

## Verification

- Core release commit: [GitHub CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174) **PASS**, both jobs.
- Adapter [PR #21](https://github.com/majiddarvishan/fluxlog/pull/21): Go-generated dependency checksums, official core `v0.3.0`, and unchanged upstream `goconfig v1.7.0`.
- Final PR: [CI 39170725761](https://github.com/majiddarvishan/fluxlog/actions/runs/39170725761) **PASS**, core and adapter jobs. All four strict CFG-02/CFG-03 tests **PASS**, no SKIP.
- Post-merge: [CI 39170859081](https://github.com/majiddarvishan/fluxlog/actions/runs/39170859081) **PASS** for both jobs on `main@0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e` (tidy, gofmt, vet, race, example builds).
- Maintainer-reported and GitHub-verified adapter publication: `GOPROXY=https://proxy.golang.org go list -m github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0` resolved to version `v0.3.0`.
- Upstream [goconfig v1.7.0 CI 39167180503](https://github.com/majiddarvishan/goconfig/actions/runs/39167180503) **PASS**.

## Outcome

**REL-02 through REL-05 are complete** for v0.3.0. Core and adapter
share the version number but retain separate Go modules and Git tag prefixes.
Core has no dependency on goconfig; the nested adapter pins matching core
v0.3.0. Historical tags are unchanged, and adapter v0.2.0 was intentionally
skipped. No standalone GitHub Release object was found during verification.

The first legacy `OnReplace` validator limitation and Manager-lifetime
post-commit observer retention remain documented and accepted.

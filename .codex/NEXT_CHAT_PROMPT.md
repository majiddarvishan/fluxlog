# Continue Fluxlog maintenance after the completed v0.3.0 release

Act as a Senior Go Software Architect / Go Modules and Release Engineer.
Explain in Persian. Keep repository code, comments, documentation and PR text in English.

## Published baseline — 2026-10-19

- Fluxlog core: annotated `v0.3.0` tag at `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd`, with [CI 39168420174](https://github.com/majiddarvishan/fluxlog/actions/runs/39168420174) PASS.
- goconfig adapter: annotated `adapters/goconfig/v0.3.0` tag at `0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e`; [PR #21](https://github.com/majiddarvishan/fluxlog/pull/21) merged and [post-merge CI 39170859081](https://github.com/majiddarvishan/fluxlog/actions/runs/39170859081) PASS for both modules.
- The maintainer verified Go proxy resolution of `github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0`. Upstream `goconfig v1.7.0` is released.
- Core and all official adapters use matching release numbers starting with v0.3.0, remain separate Go modules and use Go-compliant nested tag names. The adapter requires core v0.3.0; core never depends on goconfig.
- Old tags remain unchanged; adapter v0.2.0 was intentionally skipped. Do not repeat v0.3.0 tagging.

## Next session

1. Check actual GitHub HEAD, branches, tags, open PRs and CI before changing anything; do not treat the above SHAs as current main if docs commits or other changes have landed.
2. Read `.codex/CURRENT_STATUS.md`, `.codex/TASKS.md`, `RELEASING.md` and the release evidence document first.
3. Preserve accepted CFG-02/CFG-03 behavior, initial legacy `OnReplace` validator ordering and Manager-lifetime observers.
4. The only nonrelease backlog items are optional `PERF-01` and `CI-01`. Do not start speculative optimization, new features or another release without a concrete requirement.
5. For future releases, publish matching core and adapter version numbers with independently tested commits; never move an existing Git tag.

Do not claim a GitHub Release object exists unless verified. The v0.3.0
Go module release itself is complete.

# goconfig Adapter Changelog

This separate Go module follows **lockstep version numbers** with Fluxlog core
starting at `v0.3.0`, while retaining its own nested Git tags and dependencies:
`github.com/majiddarvishan/fluxlog/adapters/goconfig`.

## [Unreleased]

No changes recorded yet.

## [0.3.0] - 2026-10-19

### Added

- Optional `service` and `caller_max_length` configuration mapping.
- Post-commit level updates for multiple loggers bound to the same level node.
- Regression tests for invalid values, dry-run and committed Batch updates,
  and persistence failures.

### Fixed

- Do not change runtime Logger levels when FileSource persistence fails after
  candidate validation.

### Compatibility

- Requires Fluxlog core `v0.3.0` and official upstream `goconfig v1.7.0`;
  module checksums were generated using Go tooling.
- goconfig's original `OnReplace` selects the first registered pre-commit
  handler. Fluxlog must own the first validation registration on a level path.
  Post-commit observers remain registered for the Manager's lifetime.

**Published:** annotated tag `adapters/goconfig/v0.3.0` points to the verified
post-merge commit `0b5c2f46dfaf1b6ca2ae08a832d855dffe51582e`.
Both modules passed [post-merge CI 39170859081](https://github.com/majiddarvishan/fluxlog/actions/runs/39170859081).
The maintainer verified Go module resolution through `proxy.golang.org`.

## [0.1.0]

- Initial nested-module release with compatibility for the original flat
  goconfig logger configuration and runtime level binding.

Previous published adapter tag: `adapters/goconfig/v0.1.0`.

[0.3.0]: https://github.com/majiddarvishan/fluxlog/tree/adapters/goconfig/v0.3.0
[0.1.0]: https://github.com/majiddarvishan/fluxlog/tree/adapters/goconfig/v0.1.0

# goconfig Adapter Changelog

This separate Go module follows **lockstep version numbers** with the core
starting at `v0.3.0`, while retaining its own nested Git tag and dependencies:
`github.com/majiddarvishan/fluxlog/adapters/goconfig`.

## [Unreleased]

**v0.3.0 release candidate:** pin core `v0.3.0` and upstream
`goconfig v1.7.0`. Publish only after green CI, merge and post-merge CI.

### Added

- Optional `service` and `caller_max_length` configuration mapping.
- Post-commit level updates for multiple loggers bound to the same level node.
- Additional tests for invalid values, dry-run and committed Batch updates,
  and persistence failures.

### Fixed

- Do not change runtime Logger levels if FileSource persistence fails after
  candidate validation.

### Compatibility

- The current source depends on the **published `goconfig v1.7.0`** tag;
  `go.sum` was regenerated from the official module using Go tooling.
  The adapter itself still requires a separate nested-module release.
- goconfig's original `OnReplace` selects the first registered pre-commit
  handler. Fluxlog must own the first validation registration on a level path.
  Post-commit observers remain registered for the Manager's lifetime.

**Status:** These changes are on the current source branch but have not been
released as a new adapter tag. The approved candidate tag is
`adapters/goconfig/v0.3.0`; it is **not published** until created by an
authorized release operator.

## [0.1.0]

- Initial nested-module release with compatibility for the original flat
  goconfig logger configuration and runtime level binding.

Existing published adapter tag: `adapters/goconfig/v0.1.0`.

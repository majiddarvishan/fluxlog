# goconfig Adapter Changelog

This is the independently versioned Go module
`github.com/majiddarvishan/fluxlog/adapters/goconfig`.

## [Unreleased]

### Added

- Optional `service` and `caller_max_length` configuration mapping.
- Post-commit level updates for multiple loggers bound to the same level node.
- Additional tests for invalid values, dry-run and committed Batch updates,
  and persistence failures.

### Fixed

- Do not change runtime Logger levels if FileSource persistence fails after
  candidate validation.

### Compatibility

- The current source depends on an **immutable, unreleased upstream goconfig
  commit** via canonical Go pseudo-version. Before tagging a new adapter
  version, update this to an actual, published goconfig tag and regenerate
  `go.sum` using Go tooling.
- goconfig's original `OnReplace` selects the first registered pre-commit
  handler. Fluxlog must own the first validation registration on a level path.
  Post-commit observers remain registered for the Manager's lifetime.

**Status:** These changes are on the current source branch but have not been
released as a new adapter tag. A candidate tag is
`adapters/goconfig/v0.2.0`; it is **not published** until created by an
authorized release operator.

## [0.1.0]

- Initial nested-module release with compatibility for the original flat
  goconfig logger configuration and runtime level binding.

Existing published adapter tag: `adapters/goconfig/v0.1.0`.

# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

No changes recorded yet.

## [0.3.2] - 2026-10-20

### Documentation

- Update the reusable coding-agent prompt to discover the latest stable,
  Go-module-resolvable Fluxlog core and goconfig adapter tags independently
  instead of assuming a fixed v0.3.0 baseline.
- Document compatible release selection, independent adapter tag prefixes,
  module dependency checks, and explicit version pins in consuming projects.
- Align public integration examples and release status with the verified
  v0.3.1 tags as of 2026-10-20.

### Compatibility

- Documentation-only patch: no changes to Fluxlog core runtime behavior,
  public API, dependency graph, or Go toolchain requirement.
- Both module CI jobs (including vet, race tests and example builds) passed
  on [PR #27](https://github.com/majiddarvishan/fluxlog/pull/27)
  and [post-merge main CI 39326544866](https://github.com/majiddarvishan/fluxlog/actions/runs/39326544866).
- The independently versioned goconfig adapter now has its matching
  `adapters/goconfig/v0.3.2` annotated tag on the verified post-merge
  commit, with core `v0.3.2` in its dependency graph.
  [Post-merge CI 39329602416](https://github.com/majiddarvishan/fluxlog/actions/runs/39329602416)
  passed both module jobs.

## [0.3.1] - 2026-10-19

### Documentation and examples

- Add runnable package-level Fluxlog and goconfig usage examples with
  `SetDefault` ownership, runtime level changes and structured logging.
- Publish the developer/AI coding-agent integration decision guide, a
  configuration reference, practical recipes and reusable Codex prompt.
- No core API, implementation, or dependency changes since v0.3.0.

**Core published:** annotated `v0.3.1` tag points to
`fdb55d6d48ed9a8fc24b8c1e00994278b7fd1785`.
[Core source CI 39176344867](https://github.com/majiddarvishan/fluxlog/actions/runs/39176344867)
passed both Go module jobs. The matching nested goconfig adapter v0.3.1
publication is prepared separately; see its changelog and `RELEASING.md`.

## [0.3.0] - 2026-10-19

### Fixed

- Render zerolog numeric Unix timestamps as human-readable local timestamps in
  console output, honoring the configured Go time layout and preserving the
  original JSON encoding.

### Tested

- Add regression coverage for concurrent writes and logger shutdown, default
  logger ownership, real file rotation/retention/compression, output errors,
  and ANSI/NO_COLOR console formatting.
- Verify root and nested adapter modules with Go race-detector CI.

Core `v0.3.0` is published as an annotated tag on the tested source commit.
Starting at this release, official adapter versions match the core version,
while retaining separate Go modules and nested tag names. The matching
adapter `adapters/goconfig/v0.3.0` tag is also published, and both modules
passed [post-merge CI 39170859081](https://github.com/majiddarvishan/fluxlog/actions/runs/39170859081). Adapter changes
are tracked in [`adapters/goconfig/CHANGELOG.md`](adapters/goconfig/CHANGELOG.md).

## [0.2.1] - 2026-08-25

### Fixed

- Package-level logging reports the application call site instead of the
  internal convenience-function wrapper.

## [0.2.0] - 2026-08-25

### Added

- Concurrency-safe package-level logging functions with a configurable default
  Logger and a single printf-style function per level.
- Runnable package-level API example.
- Console pattern with caller and service before the message.
- Distinct cyan service name in color-enabled console output.
- Configurable caller shortening with a default maximum of 15 characters.

## [0.1.0] - 2026-08-25

### Added

- Independent, instance-based structured logger built on zerolog.
- Console and size-rotated file outputs.
- JSON and human-readable console formats.
- Runtime-safe, per-instance log-level changes.
- Optional goconfig adapter in a separate Go module.
- Compatibility with the original flat goconfig logger settings.
- Core and adapter examples, unit tests, race tests, fuzz seeds, and benchmarks.
- GitHub Actions verification for both Go modules.

[0.3.2]: https://github.com/majiddarvishan/fluxlog/tree/v0.3.2
[0.3.1]: https://github.com/majiddarvishan/fluxlog/tree/v0.3.1
[0.3.0]: https://github.com/majiddarvishan/fluxlog/tree/v0.3.0
[0.2.1]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.2.1
[0.2.0]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.2.0
[0.1.0]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.1.0

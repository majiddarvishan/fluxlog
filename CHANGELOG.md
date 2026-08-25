# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.1.0]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.1.0

# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.2.1]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.2.1
[0.2.0]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.2.0
[0.1.0]: https://github.com/majiddarvishan/fluxlog/releases/tag/v0.1.0

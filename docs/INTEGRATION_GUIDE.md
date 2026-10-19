# Fluxlog Integration Guide for Developers and AI Coding Agents

This guide describes how to **select**, **implement**, and **verify** Fluxlog
usage in a consuming Go application. It applies to human developers and to
Codex, Claude Code, and other coding agents. Make decisions from the
application's requirements, not merely from an example's defaults.

**Verified release snapshot (2026-10-20):** Fluxlog core `v0.3.1` and
goconfig adapter `v0.3.1` (requiring goconfig `v1.7.0`). These are
**two Go modules** with distinct tag prefixes, not a permanent integration
baseline. Look up the latest stable releases for each module when using
this guide; inspect the actual selected release's source and go.mod.

- [Configuration reference](CONFIGURATION_REFERENCE.md): accepted fields,
  defaults, and **adapter-specific limitations**.
- [Recipes](RECIPES.md): patterns for different application needs.
- [Core runnable examples](../examples/basic/main.go) and
  [package-level example](../examples/package-level/main.go).
- [Runnable goconfig example](../adapters/goconfig/examples/basic/main.go).

## 1. Inspect the consuming project first

Identify the Go version and existing logging package(s), configuration owner,
deployment/runtime environment, expected output destination, JSON requirements,
log-level change requirements, ownership/shutdown boundaries, and whether
multiple components need independent thresholds. Check dependency policy and
the actual public APIs before changing imports.

If deployment or operational requirements are not specified, use the simplest
compatible implementation, **state assumptions**, and avoid adding optional
providers or file rotation unnecessarily. Do not silently remove existing
logging/tracing or overwrite application-specific fields.

## 2. Decision matrix

| Decision | Choose | When to use it | Important limit |
| --- | --- | --- | --- |
| Config provider | **Core only** | Application does not already use goconfig, or needs direct core settings | Core needs no goconfig dependency |
| Config provider | **goconfig adapter** | Application already uses goconfig and its mapping meets the requirements | Adapter dynamically binds only `level` |
| Config provider | **Parse + custom core config + BindLevel** | Application uses goconfig but needs core-only constructor options (e.g. `TimeFormat`) | Construct and bind explicitly; no general hot reload |
| Logging API | **`fluxlog.Info("text")`** | One application-wide default Logger and concise messages | No fluent `.Str().Msg()` on package functions |
| Logging API | **`logger.Info().Str(...).Msg(...)`** | Structured events, injected Logger, or per-component isolation | Caller owns instance lifecycle |
| Output | **Console** | Interactive apps or stdout-based collection | Choose `ConsoleFormat` for humans or `JSONFormat` for collectors |
| Output | **File** | Local, size-rotated on-disk logs required | Requires writable directory, rotation and retention policy |
| Output | **Both** | A real requirement for concurrent console and file output | Avoid needless duplicate ingestion |
| Change level | **`logger.SetLevelString("debug")`** | Core-managed threshold at runtime | Only changes that Logger instance |
| Change level | **`manager.Replace("/logger/level", "debug")`** | goconfig-bound level from an existing Manager | Works after successful validated/committed replacement |

These choices **compose**: a default Logger may be goconfig-managed and still
expose structured events through its instance API.

### Practical selection path

1. Is goconfig **already used** and is its legacy logger mapping enough?
   If yes, use the adapter. Otherwise use core. Do **not** introduce goconfig
   merely to configure a basic logger.
2. Is there one application-wide Logger? If yes, the package-level convenience
   functions are optional; install its default once and own it for the
   application lifetime. For independently configured modules or tests,
   prefer explicit instance injection.
3. Are fields such as `request_id` or `duration_ms` needed? Use the
   **instance** structured API, including when a default was installed.
4. For stdout collectors prefer core `JSONFormat`; for human local console
   use `ConsoleFormat`. The stock goconfig adapter always chooses the
   human-readable console format and JSON files.
5. Is a non-level setting expected to hot-reload? The stock adapter does not
   implement that. Design a controlled reinitialization/lifecycle or
   application-owned configuration mechanism instead.

## 3. Resolve versions, then install

For **each required module**, check GitHub tags and the Go module proxy:

~~~bash
git ls-remote --tags https://github.com/majiddarvishan/fluxlog.git
go list -m -versions github.com/majiddarvishan/fluxlog
go list -m -versions github.com/majiddarvishan/fluxlog/adapters/goconfig
~~~

Core tags are `vX.Y.Z`; adapter tags use
`adapters/goconfig/vX.Y.Z`. Choose the highest stable SemVer tag that
is actually Go-resolvable, excluding prerelease and untagged commits.
Verify the adapter's published go.mod dependency on core and goconfig.
If the newest tags differ or are incompatible, identify a verified
compatible pair and explain why. Do not invent tags, silently fall back
to stale versions, or claim verification when the network is unavailable.
Honor explicit user-provided version constraints.

**As of 2026-10-20** (illustrative, not a permanent default):

~~~bash
go get github.com/majiddarvishan/fluxlog@v0.3.1
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.1
~~~

The v0.3.1 adapter requires core v0.3.1 and upstream goconfig v1.7.0.
Do not add the nested adapter to applications that do not need it.

## 4. API boundaries and ownership

- `fluxlog.Info("started")` logs a message using the default Logger;
  `fluxlog.Info("port %d", 8080)` supports printf-style formatting when
  the first argument is a string and **additional** arguments are present.
  `fluxlog.Info("100% ready")` treats that one string literally.
- `logger.Info().Str("request_id", "req-1").Msg("accepted")` is an instance
  event; package-level `Info` returns nothing. Both APIs may coexist.
- `fluxlog.New(config)` returns an independently owned `*fluxlog.Logger`;
  `fluxlog.SetDefault(logger)` atomically replaces the package default and
  returns the old pointer. It does **not** take ownership or close anything.
- Restore the previous default **before** closing the installed Logger.
  Keep the Manager and Logger alive for their intended usage lifetime.
  A defer in a short-lived setup helper closes the Logger too early. Do not
  concurrently replace the default in another owner and then unconditionally
  restore a stale previous Logger.
- Stop dependent goroutines, requests, and child loggers **before** closing
  the owning Logger. Handle `Close()` errors appropriately in production.
- Runtime level methods are concurrency-safe per Logger. This is not a promise
  that closing and swapping defaults while in-use log calls race is safe.
- `Fatal` and `Panic` have deliberate termination/panic behavior when
  enabled; do not use them in reusable library code or shutdown defers.
- Neither this library nor its adapter makes unrelated zerolog global-level
  changes. An external `zerolog.SetGlobalLevel` can additionally filter
  events regardless of the instance level.

## 5. goconfig binding contract

With a node at `manager.Config().At("logger")`:

~~~go
logger, err := fluxgoconfig.New(manager, loggerNode)
if err != nil {
    return err
}
// Retain logger and manager until the application shuts down.
// To update only its level after validated successful commit:
err = manager.Replace("/logger/level", "debug")
~~~

`New` parses settings, creates an instance, and calls `BindLevel`.
`BindLevel` validates candidate levels before commit, then applies
successfully committed levels. Invalid replacements and failed file
persistence do not update bound Logger levels. Multiple bound Loggers receive
the committed change.

**Registration constraint:** goconfig's legacy `OnReplace` runs the first
handler for a path. Register the adapter validation for that path *before*
unrelated pre-commit handlers; otherwise invalid level rejection is not
guaranteed. goconfig retains post-commit observers for the lifetime of its
Manager; there is no unsubscribe operation. Bind loggers for a matching
lifetime and do not assume calling `Close` unregisters a handler.

A `StrSource` example demonstrates the runtime `Replace` call; it does
**not** imply automatic file watching or external config reload. For
alternative `FileSource` behavior, use the actual goconfig source API and
test its persistence/failure path. If core settings unavailable in the
adapter are needed, use `fluxgoconfig.Parse`, customize the returned core
`Config`, call `fluxlog.New`, then `fluxgoconfig.BindLevel` (and manage
cleanup on bind failures).

## 6. Validation checklist

- Confirm imports and pinned versions, including independent module roots.
- Verify console format, timestamp layout and caller presence **from actual
  output**, not imagined sample lines.
- Verify filtered debug messages; test successful and invalid level changes.
- Test JSON fields when structured logging is required.
- For file output, test directory permissions, rotation/retention and
  output errors for the deployment target.
- For goconfig, check initial validation registration order and commit-only
  update behavior; never assume all config fields are hot-reloadable.
- Run `gofmt`, `go vet ./...`, `go test -race ./...`, and builds/tests
  appropriate to the **consuming** project. The Fluxlog repository has
  separate core and adapter CI jobs.

## 7. Copy/paste request for a coding agent

For a reusable standalone prompt maintained alongside this documentation,
see [INTEGRATE_FLUXLOG.md](../prompts/INTEGRATE_FLUXLOG.md). It is designed for
Codex and other coding agents working in a **consuming application**.

The following is a **consumer-project** prompt. The remote library's
`AGENTS.md` is **not automatically loaded** into another repository.

> Integrate `github.com/majiddarvishan/fluxlog` into this Go application.
> First inspect the application's current logging, config provider, module
> versions, deployment output and shutdown lifecycle. Read the library's
> `docs/INTEGRATION_GUIDE.md`, `docs/CONFIGURATION_REFERENCE.md`,
> `docs/RECIPES.md` and runnable examples at the version we will use.
> Query current published Git tags and Go module versions separately for
> Core and the goconfig Adapter. Select and pin the latest stable,
> Go-resolvable, compatible releases; do not assume an old documentation
> example contains the latest. Add the adapter only if the project uses
> goconfig and its mapping fits. Choose instance vs package-level logging,
> structured fields,
> console/file/JSON and runtime-level changes based on the project.
> Do not invent config fields, imply hot reload for settings other than level,
> break Logger ownership, rewrite other logging systems without cause, or
> create tags. Implement the smallest compatible change, validate with
> formatting/vet/tests/race/build, and report chosen options and limitations.

See [Recipes](RECIPES.md) for concrete follow-up tasks.

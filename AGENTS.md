# Agent guidance — Fluxlog repository

These instructions apply to contributors and coding agents **working inside
this Fluxlog repository**. They are not automatically loaded by agents working
inside a different application. For consumer-project integration, give the
agent an explicit URL/link to
[docs/INTEGRATION_GUIDE.md](docs/INTEGRATION_GUIDE.md) and the other public
documentation.

## Product boundaries

Fluxlog core (`github.com/majiddarvishan/fluxlog`) is a small Go structured
logging library built on zerolog and lumberjack. The optional goconfig adapter
(`github.com/majiddarvishan/fluxlog/adapters/goconfig`) is a **separate**
Go module. The core must not import goconfig.

As of the published lockstep `v0.3.0` release, the adapter requires core
`v0.3.0` and upstream `goconfig v1.7.0`. Modules use matching release
version numbers but distinct tags: root `vX.Y.Z` and nested
`adapters/goconfig/vX.Y.Z`. Do not create or move release tags without
explicit maintainer authorization.

## Reading order

1. [README.md](README.md) and
   [docs/INTEGRATION_GUIDE.md](docs/INTEGRATION_GUIDE.md) for when/how to
   integrate the library.
2. [docs/CONFIGURATION_REFERENCE.md](docs/CONFIGURATION_REFERENCE.md)
   for exact current Go fields versus goconfig JSON keys.
3. [docs/RECIPES.md](docs/RECIPES.md) and runnable
   [core](examples/package-level/main.go) and
   [adapter](adapters/goconfig/examples/basic/main.go) examples.
4. [RELEASING.md](RELEASING.md) for version/tag policy, plus
   [.codex/CURRENT_STATUS.md](.codex/CURRENT_STATUS.md) and
   [.codex/TASKS.md](.codex/TASKS.md) for maintenance context.
5. For API changes, inspect `config.go`, `logger.go`, `default.go` and
   `adapters/goconfig/adapter.go` before claiming support.

## Implementation rules

- Prefer small, justified changes. Preserve instance independence,
  concurrency/lifecycle guarantees, structured output and existing contracts.
- Do not silently add goconfig to core or assume the adapter supports all
  core settings; the stock adapter binds only `level` dynamically.
- Preserve the adapter's pre-commit validation / post-commit application
  behavior and the first-registered legacy `OnReplace` handler caveat.
- `SetDefault` is not ownership transfer. Restore it before closing an
  installed Logger and do not leave callers using a closed instance.
- Distinguish code at `main` from behavior actually published in immutable
  release tags; avoid describing future/unreleased work as published.
- Repository source code, documentation, commit messages and PR descriptions
  are written in English.

## Verification

The root and nested adapter require independent checks:

~~~bash
gofmt -l .
go mod tidy
git diff --exit-code -- go.mod go.sum
go vet ./...
go test -race ./...
go build ./examples/...

cd adapters/goconfig
go mod tidy
git diff --exit-code -- go.mod go.sum
go vet ./...
go test -race ./...
go build ./examples/...
~~~

These are separate commands in their respective Go module directories. The
GitHub Actions workflow verifies both modules. Do not claim a check passed
unless actually run and observed. Do not modify release tags as part of
documentation updates.

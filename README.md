# fluxlog

`fluxlog` is an instance-based structured logging library for Go. It is built
on top of [zerolog](https://github.com/rs/zerolog), supports console and
size-rotated file outputs, and allows changing the level safely at runtime.

The core package does not depend on a configuration framework and does not
modify zerolog's process-global logger or global level.

## Install

```bash
go get github.com/majiddarvishan/fluxlog@v0.3.0
```

This installs the published `v0.3.0` core module, including human-readable
console handling for numeric Unix timestamps. Its annotated tag points to
the CI-tested source commit `fe57c3e1271d5ada9339b255ac4ca0d727aef5bd`.

## Choose the right integration

Fluxlog supports several **composable** choices: direct core versus the
optional goconfig adapter; package-level `fluxlog.Info(...)` versus
instance-based structured events; human-readable or JSON output; console,
file, or both; and manually managed versus goconfig-committed level updates.

**Start here:** [Integration Guide for Developers and AI Coding Agents](docs/INTEGRATION_GUIDE.md).
The guide includes a decision matrix, installation workflow, lifecycle
constraints, and a copy/paste prompt for Codex or another coding agent.

- [Configuration Reference](docs/CONFIGURATION_REFERENCE.md) — exact core
  fields, adapter JSON mapping, defaults and unsupported adapter options.
- [Integration Recipes](docs/RECIPES.md) — six common patterns with runnable
  example links and guidance for package-level, structured, JSON, dual-output,
  goconfig and custom timestamp integrations.

**Ready-to-copy prompt:** [Integrate Fluxlog with Codex or other Coding Agents](prompts/INTEGRATE_FLUXLOG.md).
It includes the short integration request, decision criteria, release version
constraints, validation requirements, and optional scenario-specific inputs.

For a coding agent integrating this library into **another repository**:

> Use Fluxlog for this Go project's logging. First inspect the existing
> configuration provider, deployment, output requirements and Logger
> ownership. Read the [Integration Guide](docs/INTEGRATION_GUIDE.md),
> [Configuration Reference](docs/CONFIGURATION_REFERENCE.md), and
> [Recipes](docs/RECIPES.md). Use stable core `v0.3.0`, adding adapter
> `v0.3.0` only if needed. Select the smallest compatible integration,
> verify the real APIs, run tests and explain any limitations.

The `AGENTS.md` file in **this library repository** guides agents editing
Fluxlog itself. It is not automatically loaded by an agent working in the
consuming application's repository.

## Basic usage

```go
package main

import (
	"log"

	"github.com/majiddarvishan/fluxlog"
)

func main() {
	logger, err := fluxlog.New(fluxlog.Config{
		Level:     fluxlog.InfoLevel,
		Service:   "gateway",
		Timestamp: true,
		Caller:    true,
		Console: &fluxlog.ConsoleConfig{
			Format:          fluxlog.ConsoleFormat,
			Color:           fluxlog.AutoColor,
			CallerMaxLength: 15,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.Info().
		Str("request_id", "req-1").
		Msg("request received")
}
```

## Package-level convenience API

For applications that prefer direct package functions, configure a Logger as
the package default:

```go
logger, err := fluxlog.New(fluxlog.DefaultConfig())
if err != nil {
	return err
}

previous := fluxlog.SetDefault(logger)
defer func() {
	fluxlog.SetDefault(previous) // Restore the old default first.
	_ = logger.Close()           // Then close this logger's outputs.
}()

fluxlog.Info("server started")
fluxlog.Info("listening on port %d", 8080)
fluxlog.Warn("retry %d of %d", 1, 3)
```

Each level has a single package function: `Trace`, `Debug`, `Info`, `Warn`,
`Error`, `Fatal`, and `Panic`. A single value is rendered directly, including an
`error`. When the first value is a string and more arguments follow, it is used
as a `fmt.Printf` format string. Separate `Debugf`, `Infof`, and similar variants
are therefore unnecessary. Formatting is skipped when the selected level is
disabled.

```go
fluxlog.Error(err)
fluxlog.Info("processed %d requests", count)
```

`SetDefault` is concurrency-safe and returns the previous Logger. It does not
take ownership of either Logger, so callers remain responsible for `Close`.
Restore the previous default **before** closing the configured Logger, and keep
this lifecycle in the scope that owns the Logger (usually `main`, not a
short-lived initialization helper). Avoid replacing the package default again
while this owner is active: unconditional restoration would overwrite a newer
default. Package-level `fluxlog.Info(...)` also works before `SetDefault`,
but uses the library's built-in default configuration.

To change the installed Logger's level at runtime, call
`logger.SetLevelString("debug")`; subsequent `fluxlog.Debug(...)` calls use the
new level. Use the instance API when structured fields are needed:

```go
logger.Info().Str("request_id", "req-1").Msg("request received")
```

See [`examples/package-level`](./examples/package-level) for a runnable example.

Run it from the repository root with `go run ./examples/package-level`.
The example demonstrates `Info` formatting, level changes, structured events,
and default-Logger ownership.

## Runtime level changes

Runtime configuration is part of the core API and is not tied to a config
library:

```go
if err := logger.SetLevelString("debug"); err != nil {
	return err
}

logger.Debug().Msg("debug logging is now enabled")
```

Level changes affect only the target `Logger` instance. They never call
`zerolog.SetGlobalLevel`.

**External zerolog global filtering still applies.** If other application code
calls `zerolog.SetGlobalLevel(zerolog.ErrorLevel)`, a Fluxlog Logger configured
at `DebugLevel` will still suppress debug events. The effective output threshold
is constrained by both zerolog's process-global level and the Fluxlog instance
level; setting a lower instance threshold cannot override global filtering.

## Console and rotating JSON file

```go
logger, err := fluxlog.New(fluxlog.Config{
	Level:     fluxlog.InfoLevel,
	Service:   "gateway",
	Timestamp: true,
	Console: &fluxlog.ConsoleConfig{
		Format: fluxlog.ConsoleFormat,
		Color:  fluxlog.AutoColor,
	},
	File: &fluxlog.FileConfig{
		Path:       "./logs/application.log",
		Format:     fluxlog.JSONFormat,
		MaxSizeMB:  100,
		MaxBackups: 10,
		MaxAgeDays: 28,
		Compress:   true,
	},
})
```

Console-formatted logs place caller and service before the message. Caller is
shortened from the left so the filename and line number remain visible:

```text
2026-08-25T18:22:43+03:30 INF [ …_server.go:272 ] (Gateway) Web server starting on 0.0.0.0:4567
```

Console timestamps use a human-readable time layout (RFC3339 by default),
including when zerolog encodes its internal timestamp as Unix seconds,
milliseconds, microseconds, or nanoseconds. To display a simpler 24-hour
local-time layout, configure:

```go
TimeFormat: "2006-01-02 15:04:05",
```

For example, this produces `2026-10-18 20:30:45 INF (Gateway) started`
in the process's local timezone. Set `Timestamp: true` to include it.
JSON output keeps zerolog's original timestamp encoding unchanged.

`CallerMaxLength` defaults to `15`. JSON output is unchanged and retains the
full `caller` and `service` fields. The parenthesized service name is cyan when
color output is enabled; `NeverColor` and the `NO_COLOR` environment variable
disable it along with the other console colors.

File output defaults to JSON and never contains ANSI color codes.

## Lifecycle and concurrency

- A `Logger` is safe for concurrent use.
- `SetLevel` and `SetLevelString` are safe during concurrent writes.
- Call `Close` when file output is configured.
- `Close` is idempotent and waits for an in-progress write.
- Library internals return errors and never terminate the application.

## Advanced zerolog API

`Zerolog` returns a copy of the underlying logger when a zerolog-specific API
is needed. It continues to use fluxlog's managed output and runtime level:

```go
zerologLogger := logger.Zerolog()
child := zerologLogger.With().Str("component", "worker").Logger()
child.Info().Msg("worker started")
```

## Configuration adapters

Configuration adapters, including the optional `goconfig` adapter, are kept
outside the core module. This ensures applications that do not use a provider
do not inherit its dependencies.

The goconfig adapter preserves the original configuration keys and binds the
`level` node to runtime-safe per-instance level changes. Install the
published lockstep `v0.3.0` adapter with:

```bash
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0
```

The published nested module tag is `adapters/goconfig/v0.3.0`, independent
from the core `v0.3.0` Git tag. The adapter requires Fluxlog core `v0.3.0`
and upstream `goconfig v1.7.0`; it supports post-commit level updates for
multiple loggers and optional `service` / `caller_max_length` fields.
The historical adapter `v0.1.0` remains available with its earlier
dependency graph and behavior. Official modules share release **numbers**
but retain separate Go module paths and tag prefixes.

```go
logger, err := fluxgoconfig.New(manager, loggerConfigNode)
```

See [`adapters/goconfig`](./adapters/goconfig) for the full configuration and
[`adapters/goconfig/examples/basic`](./adapters/goconfig/examples/basic) for a
runnable example.

## Development

```bash
make verify
(cd adapters/goconfig && go test ./...)
```

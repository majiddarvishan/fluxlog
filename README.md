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
defer logger.Close()

previous := fluxlog.SetDefault(logger)
defer fluxlog.SetDefault(previous)

fluxlog.Info("server started")
fluxlog.Debug("listening on port %d", 8080)
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
Use the instance API when structured fields are needed:

```go
logger.Info().Str("request_id", "req-1").Msg("request received")
```

See [`examples/package-level`](./examples/package-level) for a runnable example.

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
`level` node to runtime-safe per-instance level changes. Install the latest
verified **published** adapter version with:

```bash
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.1.0
```

**Published release vs source:** `v0.1.0` is the latest verified adapter
module tag. Its `go.mod` requires Fluxlog core `v0.1.0` and `goconfig v1.3.0`.
The current adapter release-candidate source requires core `v0.3.0` and the
published `goconfig v1.7.0`, supporting post-commit level changes for multiple loggers
and optional `service` / `caller_max_length` fields. These adapter changes
are not part of the published adapter `v0.1.0` release. From `v0.3.0` onward official adapters follow core's release version number
while retaining separate Go modules. The next adapter tag is
`adapters/goconfig/v0.3.0` (not yet published); publishing core does not
publish a nested adapter.

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

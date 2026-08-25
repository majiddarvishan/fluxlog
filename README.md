# fluxlog

`fluxlog` is an instance-based structured logging library for Go. It is built
on top of [zerolog](https://github.com/rs/zerolog), supports console and
size-rotated file outputs, and allows changing the level safely at runtime.

The core package does not depend on a configuration framework and does not
modify zerolog's process-global logger or global level.

## Install

```bash
go get github.com/majiddarvishan/fluxlog@v0.2.1
```

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
`level` node to runtime-safe per-instance level changes:

```bash
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.1.0
```

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

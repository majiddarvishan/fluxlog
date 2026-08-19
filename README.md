# fluxlog

`fluxlog` is an instance-based structured logging library for Go. It is built
on top of [zerolog](https://github.com/rs/zerolog), supports console and
size-rotated file outputs, and allows changing the level safely at runtime.

The core package does not depend on a configuration framework and does not
modify zerolog's process-global logger or global level.

## Install

```bash
go get github.com/majiddarvishan/fluxlog
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
			Format: fluxlog.ConsoleFormat,
			Color:  fluxlog.AutoColor,
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
do not inherit its dependencies. The adapter is planned as a separate Go
submodule in a later implementation phase.

## Development

```bash
make verify
```

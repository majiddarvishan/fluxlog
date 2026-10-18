# Architecture and runtime behavior

## Component diagram

```text
Application
  |-- core Config --> fluxlog.New --> normalized configuration
  |                                    |
  |                                    +-- ConsoleConfig --> os.Stdout or io.Writer
  |                                    +-- FileConfig    --> parent mkdir + lumberjack.Logger
  |                                    |
  |                                    +-- outputWriter: JSON pass-through
  |                                         or zerolog.ConsoleWriter
  |                                    |
  |                                    +-- zerolog.MultiLevelWriter
  |                                              |
  |                                        managedWriter (mutex, Close)
  |                                              |
  |                                        levelWriter (atomic min level)
  |                                              |
  |                                        zerolog.Logger
  |                                            |   |
  |                                     instance  package caller-adjusted copy
  |                                            |   |
  |                                     fluent events / package-level convenience
  |
  +-- optional goconfig Manager/Node --> Parse --> fluxlog.Config
                                      --> BindLevel/OnReplace --> SetLevelString
```

Read the actual constructor in `logger.go:New` to resolve any ambiguity.

## Construction

1. `Config.normalized()` copies nested console/file configs, validates level, output modes, color and nonnegative rotation/caller settings, and supplies defaults. Neither input nor global zerolog settings are mutated.
2. Console output wraps the caller-provided `io.Writer` (or `os.Stdout`) in a `zerolog.ConsoleWriter` for human formatting; JSON is written as-is.
3. File output calls `os.MkdirAll` on its parent directory with mode `0750`; it creates a `lumberjack.Logger` and selects JSON (default) or colorless console format.
4. `zerolog.MultiLevelWriter` fans each accepted record out to configured destinations; a `managedWriter` protects writes and closers with one mutex.
5. A `levelWriter` wraps the managed output and holds the threshold in `atomic.Int32`; zerolog is constructed on top of it.
6. Optional timestamp and service context fields are attached; caller-enabled construction produces distinct instance/package caller-skip configurations.

## Event path and level filtering

```text
logger.Debug() / fluxlog.Debug(...)
  -> check levelWriter.enabled(severity) before creating event
  -> (if enabled) zerolog builds and serializes event
  -> levelWriter.WriteLevel(severity, record)
  -> recheck current threshold (drop silently if now disabled)
  -> managedWriter.WriteLevel (serialized under mutex)
  -> fan-out to console/file destination(s)
```

The double check matters when levels change during an in-flight event. Disabled instance methods return a nil `*zerolog.Event`, following zerolog's nil-safe fluent event API. Package functions skip format interpolation when disabled. Global `zerolog.SetGlobalLevel` is never invoked by fluxlog; external changes to zerolog's own global settings are outside this library's ownership.

## Concurrency and lifecycle

- Runtime levels use atomic loads/stores, scoped to each Logger.
- `managedWriter` protects each output write and `Close` using the **same mutex**, so `Close` cannot overtake an already executing managed write.
- `Close` uses `sync.Once` and returns the first stored error, joining multiple closer errors. Subsequent managed writes return `io.ErrClosedPipe`.
- The Logger closes only **owned file rotation sinks**, not an application-supplied console writer.
- The package default is `atomic.Pointer[Logger]`. `SetDefault` swaps a pointer; it **does not** close the old logger or take ownership of the new one.
- Because a goroutine may have obtained the previous default pointer just before a swap, applications must coordinate replacing/closing loggers with their in-flight work.

## Formatting details

- Console part order: timestamp, level, caller, service, message; other fields follow as zerolog fields.
- Service is formatted as `(CapitalizedService)`, with cyan decoration when console color is enabled. The separate `service` field is excluded from the extra console-fields section; JSON keeps it.
- Caller paths are made relative to the current working directory when possible and truncated **from the left** to preserve the filename/line number. `CallerMaxLength` defaults to 15. Truncation is rune-aware.
- `AutoColor` checks whether the output is an `*os.File` with character-device mode; `NeverColor` disables ANSI. The service formatter also checks `NO_COLOR`; exercise complete color handling with tests before changing it.
- If file `Format` is explicitly `console`, the formatter is forced to `NeverColor` and uses the console configuration's time format and caller length when present (otherwise defaults).

## Architectural boundaries

- Core package deliberately does not import goconfig.
- Optional adapter imports both core fluxlog and goconfig; other field changes are not rebound.
- `Fatal(...)` and `Panic(...)` expose underlying zerolog process termination/panic semantics when *called by the application*. Core initialization and internal runtime code return errors instead.
- `Zerolog()` returns a **copy** sharing managed output and dynamic threshold; it does not transfer ownership.
- This path is synchronous: writes go through an in-process mutex; there is no retry queue or remote sink.

Primary evidence: `config.go`, `logger.go`, `writer.go`, `console.go`, `default.go`, `adapters/goconfig/adapter.go`.

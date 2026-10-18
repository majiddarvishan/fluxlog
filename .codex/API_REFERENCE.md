# Core public API and behavioral contracts

## Construction and configuration

```go
logger, err := fluxlog.New(fluxlog.Config{
    Level: fluxlog.InfoLevel,
    Service: "gateway",
    Timestamp: true,
    Caller: true,
    Console: &fluxlog.ConsoleConfig{Format: fluxlog.ConsoleFormat},
    File: &fluxlog.FileConfig{Path: "./logs/app.log", Format: fluxlog.JSONFormat},
})
if err != nil { return err }
defer logger.Close()
```

At least one of `Console` or `File` must be non-nil. `Config` construction settings are not mutable through the Logger after creation, except for the minimum log level.

`DefaultConfig()` yields an info-level, timestamped console logger; `New` returns validation or filesystem errors.

## Severity APIs

| API | Behavior |
| --- | --- |
| `ParseLevel(string) (Level, error)` | Case-insensitive, trimmed; accepts `warning`, `off` aliases |
| `SetLevel(Level) error` | Validates and atomically updates one instance |
| `SetLevelString(string) error` | Parses then applies; invalid input leaves level intact |
| `Level() Level` | Reads current instance threshold |
| `Trace/Debug/Info/Warn/Error/Fatal/Panic() *zerolog.Event` | Fluent structured events; may return nil when filtered |
| `Close() error` | Safe and idempotent; releases owned file outputs |

`DisabledLevel` suppresses events. `Fatal().Msg(...)` and `Panic().Msg(...)` intentionally trigger zerolog exit/panic behavior when enabled. A zero/empty `Config.Level` defaults to Info, whereas `ParseLevel("")` is invalid.

## Structured and advanced APIs

```go
logger.Info().Str("request_id", "abc").Int("attempt", 2).Msg("accepted")
z := logger.Zerolog()
z.With().Str("component", "worker").Logger().Info().Msg("ready")
ctx := logger.WithContext(parentCtx) // retrieve through zerolog.Ctx(ctx)
fields := logger.With()             // zerolog.Context
```

`Zerolog()` is a value copy still writing into the managed destination. Custom children and contexts must not outlive the owning Logger if they need to emit logs.

## Package-level convenience API

```go
previous := fluxlog.SetDefault(logger)
defer fluxlog.SetDefault(previous)
fluxlog.Info("started")
fluxlog.Debug("request %d", 7)
fluxlog.Error(err)
```

`Default() *Logger` retrieves the atomic package pointer. `SetDefault(*Logger) *Logger` returns its previous value, **panics on nil**, and neither closes nor owns either instance.

The seven package-level funcs accept `...any` with these exact rules:

- No arguments: empty message.
- One string: literal message (not used as printf format).
- One non-string: `fmt.Sprint(value)` (including `error`).
- At least two arguments, first a string: `Msgf(first, rest...)`.
- At least two arguments, first non-string: `fmt.Sprint(values...)`.
- No argument formatting is performed if the event was already filtered out.

For structured fields, use the instance `*zerolog.Event` API instead of package functions.

## Output errors and consumption

The application must manage output paths, permissions and lifecycle. `Close` errors can be observed; normal zerolog `Msg` calls do not give a per-write error return through this wrapper. Do not assume remote delivery, durability or a flush/ACK protocol. Avoid closing a Logger while other goroutines still need it.

## Legacy adapter API (different module)

`goconfig.New(manager, node) (*fluxlog.Logger, error)`, `Parse(node) (fluxlog.Config, *config.Node, error)`, and `BindLevel(manager, levelNode, logger) error`.

The adapter maps the legacy keys, creates a logger, applies the level, and registers `OnReplace` for that node. Other changes require a new Logger or separate application logic. Errors are returned for absent nodes/types and unknown modes.

Sources: `logger.go`, `default.go`, `level.go`, `writer.go`, `adapters/goconfig/adapter.go`; examples in `examples/`.

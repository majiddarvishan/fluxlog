# Fluxlog Integration Recipes

Use [Integration Guide](INTEGRATION_GUIDE.md) to pick a pattern and
[Configuration Reference](CONFIGURATION_REFERENCE.md) for supported
fields and defaults. Examples below use the published core `v0.3.0`
and, where relevant, independent goconfig adapter `v0.3.0`.

## Recipe 1 — Small application, package-level Info

**Choose when:** there is one application-wide Logger and concise messages
are sufficient.

~~~go
logger, err := fluxlog.New(fluxlog.Config{
    Level:     fluxlog.InfoLevel,
    Service:   "api",
    Timestamp: true,
    Caller:    true,
    Console: &fluxlog.ConsoleConfig{
        Format:     fluxlog.ConsoleFormat,
        Color:      fluxlog.AutoColor,
        TimeFormat: "2006-01-02 15:04:05",
    },
})
if err != nil {
    return err
}
previous := fluxlog.SetDefault(logger)
defer func() {
    fluxlog.SetDefault(previous)
    _ = logger.Close()
}()
fluxlog.Info("server started on port %d", 8080)
~~~

Put the deferred cleanup in the **long-lived owner** (often `main`).
The stock package default works without `SetDefault`, but it does not
automatically inherit custom service, output, or timestamps.

**Run:** from the repository root `go run ./examples/package-level`.
The [full example](../examples/package-level/main.go) includes manual
runtime level changes and structured events.

## Recipe 2 — Structured events and isolated components

**Choose when:** you need searchable event fields or separately owned
Logger instances.

~~~go
logger.Info().
    Str("request_id", requestID).
    Str("component", "http").
    Int("status_code", 200).
    Msg("request completed")
~~~

Prefer a `*fluxlog.Logger` injected into a component rather than having
each component modify the process's package default. Use type-appropriate
structured fields; `fluxlog.Info("request %s", id)` remains a text
message and does not create an indexed `request_id` field.

**Run:** `go run ./examples/basic` at the repository root.

## Recipe 3 — JSON stdout for container collection

**Choose when:** the deployment collects stdout/stderr and needs machine-
readable events; avoid a local file unless a file is explicitly required.

~~~go
logger, err := fluxlog.New(fluxlog.Config{
    Level:     fluxlog.InfoLevel,
    Service:   "api",
    Timestamp: true,
    Console: &fluxlog.ConsoleConfig{
        Format: fluxlog.JSONFormat,
        Color:  fluxlog.NeverColor,
    },
})
~~~

Handle `err` and close the Logger at the end of its owned lifetime.
This is a **core** feature. The stock goconfig adapter selects human-readable
console formatting; it does not accept `console_format` as a JSON key.
For a goconfig-based app needing JSON stdout, customize the parsed core
config before construction and bind its level separately.

## Recipe 4 — Both console and rotating JSON file

**Choose when:** an operational requirement explicitly calls for both
human-readable local output and locally rotated structured files.

~~~go
logger, err := fluxlog.New(fluxlog.Config{
    Level:     fluxlog.InfoLevel,
    Timestamp: true,
    Service:   "api",
    Console: &fluxlog.ConsoleConfig{
        Format: fluxlog.ConsoleFormat,
        Color:  fluxlog.AutoColor,
    },
    File: &fluxlog.FileConfig{
        Path:       "./logs/api.log",
        Format:     fluxlog.JSONFormat,
        MaxSizeMB:  100,
        MaxBackups: 10,
        MaxAgeDays: 28,
        Compress:   true,
    },
})
~~~

Check destination permissions and retention policy. `New` creates the
parent directory but can fail on a read-only filesystem. The file uses
lumberjack size rotation, not a remote/guaranteed delivery service.
This fully customizable setup uses the core module.

**Run:** `go run ./examples/file` from the root.

## Recipe 5 — Existing goconfig application with runtime level changes

**Choose when:** the project already owns a `*goconfig.Manager`, and the
adapter's accepted JSON fields cover the logging requirements.

~~~json
{
  "logger": {
    "output_mode": "console",
    "level": "info",
    "service": "gateway",
    "caller_max_length": 15
  }
}
~~~

~~~go
node, err := manager.Config().At("logger")
if err != nil {
    return err
}
logger, err := fluxgoconfig.New(manager, node)
if err != nil {
    return err
}
previous := fluxlog.SetDefault(logger)
defer func() {
    fluxlog.SetDefault(previous)
    _ = logger.Close()
}()

fluxlog.Info("application started")
if err := manager.Replace("/logger/level", "debug"); err != nil {
    return err
}
fluxlog.Debug("debug enabled after committed update")
~~~

The `manager` above is created by the **consuming application**. A complete
standalone source/manager setup is in the
[runnable adapter example](../adapters/goconfig/examples/basic/main.go):

~~~bash
cd adapters/goconfig
go run ./examples/basic
~~~

This example uses an in-memory JSON `StrSource`. Do not interpret it as
automatic on-disk file watching. The example's `both` mode writes a JSON
log file under `./logs/` relative to the running working directory.
Pre-commit validation must be registered first for that level path;
failed persistence and invalid levels do not update runtime loggers.

## Recipe 6 — goconfig, but custom console time layout

**Choose when:** the application uses goconfig yet needs settings beyond
the stock adapter's JSON mapping (e.g. a 24-hour console time layout).

~~~go
cfg, levelNode, err := fluxgoconfig.Parse(node)
if err != nil {
    return err
}
if cfg.Console != nil {
    cfg.Console.TimeFormat = "2006-01-02 15:04:05"
}
logger, err := fluxlog.New(cfg)
if err != nil {
    return err
}
if err := fluxgoconfig.BindLevel(manager, levelNode, logger); err != nil {
    _ = logger.Close()
    return err
}
// Own logger and manager until the application stops.
~~~

This retains the adapter's **level-only** runtime binding; it does not
create a new `time_format` goconfig key or make TimeFormat hot-reloadable.
Manage default installation, request lifetimes, and cleanup just as with
the other recipes.

## Common mistakes to avoid

- Importing goconfig into a core-only application for no operational reason.
- Assuming a package function returns an event; use
  `logger.Info().Str(...).Msg(...)` for structured fields.
- Setting `TimeFormat` in the stock adapter's JSON; the adapter does not
  consume it.
- Calling `logger.Close()` while it remains the installed default or while
  active goroutines still use it.
- Deferring Logger cleanup in a short-lived initialization helper.
- Calling `manager.Replace` and assuming **all** construction fields
  update live; only `level` is bound.
- Registering another legacy `OnReplace` for the same path before the
  adapter's level validator.
- Assuming a new README on `main` exists in a previously tagged release.
- Using `Fatal` or `Panic` in reusable logging or cleanup code.

For a reusable AI coding-agent prompt, see the
[Integration Guide](INTEGRATION_GUIDE.md#7-copypaste-request-for-a-coding-agent).

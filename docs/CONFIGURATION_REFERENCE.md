# Fluxlog Configuration Reference

**Applies to:** current public API supported by Fluxlog core `v0.3.0` and
optional goconfig adapter `v0.3.0`. The library's
[construction code](../config.go) and
[adapter parser](../adapters/goconfig/adapter.go) are authoritative if a
future release changes behavior. **Go constructor fields and JSON adapter keys
are different interfaces.**

## Core constructor: `fluxlog.Config`

| Go field | Type | Rule / default |
| --- | --- | --- |
| `Level` | `fluxlog.Level` | Empty -> `info`; otherwise validated |
| `Service` | `string` | Trimmed; omitted when blank |
| `Timestamp` | `bool` | Zero-value `false` (but `DefaultConfig()` sets `true`) |
| `Caller` | `bool` | Zero-value `false` |
| `Console` | `*fluxlog.ConsoleConfig` | Non-nil enables console/stream output |
| `File` | `*fluxlog.FileConfig` | Non-nil enables size-rotated file output |

**At least one** of `Console` and `File` is required. Use both for dual
output. `DefaultConfig()` supplies info-level, timestamped, human-readable
console output with auto color, RFC3339 time and a caller limit of 15. It does
**not** enable `Caller`.

### `ConsoleConfig`

| Field | Type | Rule / default |
| --- | --- | --- |
| `Writer` | `io.Writer` | `os.Stdout` when nil; caller owns the writer |
| `Format` | `OutputFormat` | `ConsoleFormat` when empty; `ConsoleFormat` or `JSONFormat` |
| `Color` | `ColorMode` | `AutoColor` when empty; `AutoColor`, `AlwaysColor`, `NeverColor` |
| `TimeFormat` | `string` | Go time layout; RFC3339 when empty |
| `CallerMaxLength` | `int` | 15 when zero; negative rejected; truncation applies to human console |

Console format is human-readable. A JSON-formatted console `Writer`
provides structured lines for stdout collectors. `TimeFormat` is interpreted
by the console formatter; it does not change JSON timestamp encoding.
Use `Timestamp: true` to include timestamps. For local 24-hour time:

~~~go
Console: &fluxlog.ConsoleConfig{
    Format:     fluxlog.ConsoleFormat,
    Color:      fluxlog.NeverColor,
    TimeFormat: "2006-01-02 15:04:05",
},
~~~

The console formatter converts supported zerolog numeric Unix timestamps
(seconds, milliseconds, microseconds or nanoseconds) to human-readable
timestamps in the process's local timezone. Color behavior also depends on
output/terminal conditions; explicitly use `NeverColor` when plain text is
required.

### `FileConfig`

| Field | Type | Rule / default |
| --- | --- | --- |
| `Path` | `string` | Required; whitespace trimmed; parent directory created |
| `Format` | `OutputFormat` | `JSONFormat` when empty; either JSON or console is valid |
| `MaxSizeMB` | `int` | 100 when zero; negative rejected; **megabytes** |
| `MaxBackups` | `int` | Passed to lumberjack; negative rejected |
| `MaxAgeDays` | `int` | Passed to lumberjack; negative rejected; days |
| `Compress` | `bool` | Default `false` |
| `LocalTime` | `bool` | Default `false`; applies to rotation timestamps |

File output uses lumberjack and needs write permissions. Zero
`MaxBackups`/`MaxAgeDays` values retain **lumberjack** semantics; they are
not Fluxlog-specific rotation presets. File output does not include console
ANSI colors. For normal file ingestion prefer `JSONFormat`.

## Core level values

`trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`,
`disabled`. Parsing is trimmed/case-insensitive and additionally accepts
`warning` for `warn`, `off` for `disabled`. Zero-value `Config.Level`
defaults to info; `ParseLevel("")` returns an error.

`logger.SetLevel(level)` and `logger.SetLevelString(value)` change
**only that instance's minimum level** and are safe alongside concurrent
logging. They do not mutate `Config` fields, reconfigure outputs or defeat
external zerolog global-level filtering.

## goconfig adapter input mapping (nested Go module)

The adapter expects a **flat logger node**; a consuming application may
place that node under a parent such as `logger`.

~~~json
{
  "logger": {
    "output_mode": "both",
    "level": "info",
    "service": "gateway",
    "caller_max_length": 15,
    "file_name": "./logs/application.log",
    "max_file_size": 100,
    "max_files": 10
  }
}
~~~

| goconfig key | Required? | Mapped core setting |
| --- | --- | --- |
| `output_mode` | Yes | `console`, `file` or `both`, case-insensitive and trimmed |
| `level` | Yes | `Level`, validated with `ParseLevel` |
| `service` | No | `Service`; string |
| `caller_max_length` | No | `Console.CallerMaxLength`; integer |
| `file_name` | For `file`/`both` | `File.Path`; string |
| `max_file_size` | For `file`/`both` | `File.MaxSizeMB`; integer, MB |
| `max_files` | For `file`/`both` | `File.MaxBackups`; integer |

The adapter fixes these construction defaults (not JSON fields):

| Core setting | Adapter behavior |
| --- | --- |
| `Timestamp`, `Caller` | Both **true** |
| Console format | `ConsoleFormat` (human readable) |
| Console color | `AutoColor` |
| Console time layout | Core default RFC3339 |
| File format | `JSONFormat` |
| File `MaxAgeDays` | 28 days |
| File `Compress`, `LocalTime` | Core booleans default to false |

**Not accepted as stock adapter configuration keys:**
`time_format`, `console_format`, `console_color`, `timestamp`,
`caller`, `max_age_days`, `compress`, `local_time` and general
`ConsoleConfig` / `FileConfig` fields. Do not invent them in an application
schema and assume the adapter consumes them. For a custom core-only option,
use `fluxgoconfig.Parse(node)` to get a `fluxlog.Config` and level node,
customize the returned core config, construct using `fluxlog.New`, then
`fluxgoconfig.BindLevel(manager, levelNode, logger)`. Handle a failed bind
by closing the newly created Logger.

**Only the `level` node is dynamically bound.** `service`, output mode,
file rotation settings and formatting do not change just because goconfig
commits their new values.

### Runtime persistence and registration

`BindLevel` uses `Manager.OnReplace` for pre-commit validation and
`Manager.OnReplaceCommitted` to apply the final committed value.
Invalid levels or failed persistence do not change the runtime Logger level;
successful replacements and batches notify all bound Logger instances.

The goconfig legacy `OnReplace` mechanism selects the **first** matching
pre-commit handler. Bind/register Fluxlog's validating handler first for the
level path. Committed observers remain registered until the Manager's lifetime
ends; there is no unsubscribe API. An in-memory `StrSource` demonstration
does **not** implement automatic file-watching or remote configuration
delivery. Inspect goconfig's actual source and persistence configuration
before claiming these capabilities.

## Package-level and instance APIs

| API | Behavior |
| --- | --- |
| `fluxlog.Info("ready")` | Text logged via the default instance |
| `fluxlog.Info("port %d", 8080)` | Printf-style formatting when a string plus additional values are supplied |
| `logger.Info().Str("request_id", id).Msg("ready")` | Structured event using a specific instance |
| `fluxlog.Default()` | Returns the current default instance |
| `fluxlog.SetDefault(logger)` | Atomically swaps default and returns previous; panics on nil |
| `logger.Close()` | Releases owned outputs; caller manages lifecycle |

A package-level function **does not** return a `zerolog.Event` and cannot
be chained with `.Str()`. `SetDefault` does not close or own the previous
or new Logger. Stop all users before closing an instance, and restore a
previous default before closing the instance it replaced.

## Versioning and verification

~~~bash
go get github.com/majiddarvishan/fluxlog@v0.3.0
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.0
~~~

Import the adapter only when required. The adapter is independently published
under `adapters/goconfig/v0.3.0` and requires core `v0.3.0` plus
`goconfig v1.7.0`. No adapter `v0.2.0` was published. Distinct tag commits
are valid under the [lockstep release policy](../RELEASING.md).

See [Integration Guide](INTEGRATION_GUIDE.md) and [Recipes](RECIPES.md).

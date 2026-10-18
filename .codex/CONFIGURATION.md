# Configuration contract and defaults

## Core `Config`

| Field | Default / rule |
| --- | --- |
| `Level` | Empty → `info`; otherwise a validated `Level` |
| `Service` | Trimmed string; omitted if empty |
| `Timestamp` | Boolean, default false in a zero-value `Config` |
| `Caller` | Boolean, default false in a zero-value `Config` |
| `Console` | Optional `*ConsoleConfig` |
| `File` | Optional `*FileConfig` |

`Console == nil && File == nil` returns an error.

### ConsoleConfig

| Field | Default / rule |
| --- | --- |
| `Writer` | `os.Stdout` when nil; caller retains ownership |
| `Format` | `console` when empty; allowed: `console`, `json` |
| `Color` | `auto` when empty; allowed: `auto`, `always`, `never` |
| `TimeFormat` | `time.RFC3339` when empty |
| `CallerMaxLength` | 15 when zero; negative rejected |

`DefaultConfig()` is an explicit convenience preset with Info, timestamp=true, console output, auto color, RFC3339, and caller max length=15. It does **not** automatically turn on `Caller`.

### FileConfig

| Field | Default / rule |
| --- | --- |
| `Path` | Required non-whitespace path; trimmed |
| `Format` | `json` when empty; `console` also allowed |
| `MaxSizeMB` | 100 when zero; negative rejected |
| `MaxBackups` | Passed to lumberjack; negative rejected |
| `MaxAgeDays` | Passed to lumberjack; negative rejected |
| `Compress` | Boolean, default false |
| `LocalTime` | Boolean, default false |

File output uses lumberjack and creates the parent directory at construction. Zero backup/age values are passed through rather than assigned bespoke fluxlog defaults; inspect lumberjack's documentation for its retention semantics. There is no periodic arbitrary hot-reload of these settings.

## Adapter input shape

The separate module `adapters/goconfig` reads a goconfig `*Node` with these exact flat keys:

```json
{
  "output_mode": "both",
  "level": "info",
  "service": "gateway",
  "caller_max_length": 15,
  "file_name": "./logs/application.log",
  "max_file_size": 100,
  "max_files": 10
}
```

| Key | Rules |
| --- | --- |
| `output_mode` | Required string: `console`, `file`, or `both` (case-insensitive, trimmed) |
| `level` | Required string, validated through core `ParseLevel` |
| `service` | Optional string |
| `caller_max_length` | Optional int; defaults via core when zero |
| `file_name` | Required when file/both |
| `max_file_size` | Required int when file/both, MB |
| `max_files` | Required int when file/both |

The adapter forces `Timestamp=true`, `Caller=true`, console human format/auto color, file JSON format, and `MaxAgeDays=28`. For `file`/ `both`, no other file settings are accepted by the adapter API as currently written.

`Parse` translates the document into core config and returns the actual `level` node. `BindLevel` applies its current value and registers `Manager.OnReplace` on that node; only that logger's minimum level changes at runtime. The adapter's callback returns an error on invalid replacement (its comment states the manager rolls back failed replacements, but dedicated end-to-end rollback verification is not present in its test file).

## Modules and installation references

Root module name: `github.com/majiddarvishan/fluxlog` (root `go.mod`, Go 1.23). Adapter module name: `github.com/majiddarvishan/fluxlog/adapters/goconfig` (nested `go.mod`, Go 1.23).

**Documentation discrepancy at the baseline:** root `README.md` suggests installing the adapter at `@v0.1.0`, whereas `adapters/goconfig/README.md` suggests `@v0.2.1` and the adapter `go.mod` requires fluxlog `v0.2.1`. Verify actual module tags and correct documentation before making a release recommendation; this is **not** proof of a runtime defect.

Evidence: `config.go`, `logger.go`, `adapters/goconfig/adapter.go`, both `go.mod` files, both README files.

## 2026-10-18 dependency direction

The maintainer selected `github.com/majiddarvishan/goconfig v1.6.0` for the nested adapter. This is a dependency change, not a change in the legacy adapter's accepted JSON fields or runtime-reload contract. The post-tidy `go.mod`/`go.sum` pair was regenerated using the Go 1.23 GitHub Actions output and passed both CI jobs in [run 39144064913](https://github.com/majiddarvishan/fluxlog/actions/runs/39144064913). The upgrade removes the unused adapter-module indirect dependencies `gorilla/mux` and `rs/cors`. This verification applies to the PR branch until merged.

**Tag verification:** GitHub listed `adapters/goconfig/v0.1.0` but not `adapters/goconfig/v0.2.1` on 2026-10-18. A root `v0.2.1` tag does not publish the nested adapter module. See `TASKS.md:REL-01`.

## Phase 1 source-review caveat (2026-10-18)

At goconfig v1.6.0, an invalid `level` replacement is rejected by `applyLevel` and the config node is reverted. However, if the callback **succeeds** and subsequent `FileSource.setConfig` **fails**, goconfig reverts its node without compensating the externally changed `fluxlog.Logger` level. This is a source-backed consistency risk, distinct from invalid-value rollback; the previously missing post-persist error-path regression test has since reproduced the mismatch (CFG-03 remains open).

Registering two `BindLevel` callbacks for the same goconfig path does not make them broadcast listeners: v1.6.0 manager selects the first matching handler on replacement. See `REVIEW_2026-10-18.md`; define a supported binding policy before making changes.

## Console human-readable timestamp contract (FMT-01)

Human-readable `console` output always formats a timestamp according to `ConsoleConfig.TimeFormat` (Go time layout, default `time.RFC3339`). Example compact 24-hour local time: `"2006-01-02 15:04:05"`. Set `Config.Timestamp = true` to include a timestamp in emitted events.

When zerolog's `TimeFieldFormat` encodes Unix seconds, milliseconds, microseconds or nanoseconds, the custom console writer converts its numeric JSON timestamp to a `time.Time` in the **process local timezone** and displays the configured human-readable layout. Subsecond precision is preserved when that layout includes it. String-formatted timestamps continue to use the original parsing path. JSON output remains raw/unchanged.

FMT-01 is implemented and tested in [PR #13](https://github.com/majiddarvishan/fluxlog/pull/13). See [Actions run 39147165603](https://github.com/majiddarvishan/fluxlog/actions/runs/39147165603). No global timestamp format mutation was added to the core library.

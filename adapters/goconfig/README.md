# fluxlog goconfig adapter

This optional module creates a `fluxlog.Logger` from
`github.com/majiddarvishan/goconfig` and updates only that logger instance when
the `level` node is successfully persisted.

## Install

The independent goconfig adapter is prepared for lockstep `v0.3.2`.
**Use this version only after verifying the nested
`adapters/goconfig/v0.3.2` tag exists and resolves as a Go module**:

```bash
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.3.2
```

The `v0.3.2` adapter requires published Fluxlog core `v0.3.2` and
upstream `goconfig v1.7.0`. Historical published `v0.3.1` remains
available. Core never depends on goconfig. For new integrations, consult
the [latest stable tag selection guide](../../docs/INTEGRATION_GUIDE.md#3-resolve-versions-then-install)
and select only verified, compatible published tags.

The earlier adapter `v0.1.0` remains published as historical compatibility
with core `v0.1.0` and goconfig `v1.3.0`; it does not include the
v0.3.0 integration improvements.

## Choose and validate an integration

- [Integration Guide](../../docs/INTEGRATION_GUIDE.md) — decision matrix
  for core vs adapter, default vs instance API and runtime-level management.
- [Configuration Reference](../../docs/CONFIGURATION_REFERENCE.md) — the
  **exact** adapter JSON mapping; it does not expose every core setting.
- [Recipes](../../docs/RECIPES.md) — examples for an existing goconfig
  Manager and custom core settings via `Parse` + `BindLevel`.

The stock adapter dynamically binds only the `level` field. For example,
`time_format` is **not** an adapter JSON option. To customize the console
time layout, parse the node into a core Config, set
`Config.Console.TimeFormat` before construction, then bind the returned
level node explicitly.

## Current-source configuration

Expected configuration:

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

`file_name`, `max_file_size`, and `max_files` are required only when
`output_mode` is `file` or `both`. `service` and `caller_max_length` are
optional; caller length defaults to `15`.

```go
source, err := config.NewStrSource(configJSON, schemaJSON)
if err != nil {
	return err
}
manager, err := config.NewManager(source)
if err != nil {
	return err
}

logger, err := fluxgoconfig.New(manager, manager.Config())
if err != nil {
	return err
}

previous := fluxlog.SetDefault(logger)
defer func() {
	fluxlog.SetDefault(previous) // Restore the default before closing.
	_ = logger.Close()
}()

fluxlog.Info("server started")
// Given a "/level" node at the root, this enables package-level debug:
// manager.Replace("/level", "debug")
```

The core `github.com/majiddarvishan/fluxlog` module has no dependency on
`goconfig`; only applications importing this adapter module receive it.

A complete runnable program is available in [`examples/basic`](./examples/basic).
From the repository's nested module, run:

```bash
cd adapters/goconfig
go run ./examples/basic
```

The example configures both console output and a JSON file at
`./logs/application.log`, prints package-level `fluxlog.Info(...)` messages,
then calls `manager.Replace("/logger/level", "debug")`. The adapter applies
the new level after goconfig commits the replacement, so `fluxlog.Debug(...)`
becomes visible without rebuilding the logger. An invalid replacement is
rejected and leaves the last valid level unchanged. The example also
demonstrates a structured event with `logger.Info().Str(...).Msg(...)`.

`SetDefault(logger)` changes only the package-level default; it does **not**
take ownership of the Logger or close the previous default. Keep the Manager
and Logger alive for the intended logging lifetime, restore the previous
default before closing the Logger, and do not put this defer in a short-lived
initialization helper. goconfig retains the registered callbacks for the
Manager's lifetime.

## Runtime level updates (v0.3.0 and later)

The published adapter requires the official upstream `goconfig v1.7.0` tag,
verified against its annotated tag on commit
`213e87a1677af862fa9d6226f4e6ea8d0b0edfbc`.

`BindLevel` applies the initial level, validates proposed replacements
without changing Logger state, and uses `OnReplaceCommitted` to apply a new
level only **after configuration persistence succeeds**. Multiple loggers
bound to the same node receive the committed update. A failed FileSource
write leaves all runtime levels unchanged. The same contract applies to a
successful `Batch`; the final value is delivered once per replaced path.

**Registration limitation:** goconfig's original `OnReplace` chooses only
the first callback registered for a path. To guarantee invalid-level
rejection, the adapter must own the first pre-commit validation callback
for that path. Do not register an unrelated `OnReplace` callback before
`BindLevel` for the same level node. Committed observers are retained
until the Manager is discarded (there is no unsubscribe method).

This behavior was introduced in `v0.3.0` and remains in `v0.3.1`.
The adapter and core have independent tags which may point to distinct
CI-tested commits. No source checkout or local replacement is needed
to install the published modules.

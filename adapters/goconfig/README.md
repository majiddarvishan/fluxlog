# fluxlog goconfig adapter

This optional module creates a `fluxlog.Logger` from
`github.com/majiddarvishan/goconfig` and updates only that logger instance when
the `level` node is successfully persisted.

```bash
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.1.0
```


## Published module vs current source

The command above installs the latest **tagged adapter release**, `v0.1.0`,
from the independently versioned Go module under `adapters/goconfig`.
The corresponding repository tag is `adapters/goconfig/v0.1.0`.
A root-module tag such as `v0.2.1` is **not** an adapter tag, and
`adapters/goconfig/v0.2.1` has not been published as of 2026-10-18.

| | Published adapter `v0.1.0` | Current `main` source |
| --- | --- | --- |
| Fluxlog core dependency | `v0.1.0` | `v0.2.1` |
| goconfig dependency | `v1.3.0` | `v1.6.1-0.20261018202246-20a58eb9a426` (immutable pre-release revision) |
| Adapter `service` / `caller_max_length` settings | Not supported by adapter mapping | Supported as optional fields |
| Release status | Existing tag | Unreleased changes |

The configuration and usage below describe the **current `main` source**.
For the tagged `v0.1.0` version, consult its
[release-tag README](https://github.com/majiddarvishan/fluxlog/blob/adapters/goconfig/v0.1.0/adapters/goconfig/README.md).
To test unreleased source, check out the repository and set up a Go workspace
or local module replacements for the root and nested adapter modules; do not
assume the published `v0.1.0` dependency contains the current changes.

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
defer logger.Close()

previous := fluxlog.SetDefault(logger)
defer fluxlog.SetDefault(previous)

fluxlog.Info("server started")
```

The core `github.com/majiddarvishan/fluxlog` module has no dependency on
`goconfig`; only applications importing this adapter module receive it.

A complete runnable program is available in [`examples/basic`](./examples/basic).

## Runtime level updates in current development source

This development branch pins an **immutable upstream commit**
(`goconfig@20a58eb9a426d0ed7d522b1e5848a0971ce58198`) using a Go
pseudo-version. This is **not a published goconfig release** and the tagged
adapter `v0.1.0` still uses the legacy behavior.

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

No release has been tagged for the new adapter behavior. To use the code
here before an adapter release, work from a local checkout or a Go workspace
and follow the exact source revision in this module's `go.mod`.

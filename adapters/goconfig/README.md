# fluxlog goconfig adapter

This optional module creates a `fluxlog.Logger` from
`github.com/majiddarvishan/goconfig` and updates only that logger instance when
the `level` node is replaced at runtime.

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
| goconfig dependency | `v1.3.0` | `v1.6.0` |
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

## Runtime replacement limitations in goconfig v1.6.0

The current adapter binds the `level` key through `Manager.OnReplace`.
With goconfig v1.6.0, replacement invokes only the **first** matching
registered callback for a configuration path. Multiple loggers bound to the
same `level` node are therefore not a supported broadcast configuration.

A separate known consistency issue exists for `FileSource`: after a level
callback succeeds, a subsequent file-persistence error can roll back the
configuration node without reverting the Logger's changed runtime level.
These limitations are covered by known-issue regression test SKIPs
(`CFG-02` and `CFG-03`); **they are not fixed by the release documentation
change**. Use one Logger binding per managed level path and avoid relying on
automatic rollback of external side effects on persistence errors.

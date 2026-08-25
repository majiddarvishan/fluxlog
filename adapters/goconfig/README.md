# fluxlog goconfig adapter

This optional module creates a `fluxlog.Logger` from
`github.com/majiddarvishan/goconfig` and updates only that logger instance when
the `level` node is replaced at runtime.

```bash
go get github.com/majiddarvishan/fluxlog/adapters/goconfig@v0.2.1
```

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

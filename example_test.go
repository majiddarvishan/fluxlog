package fluxlog_test

import (
	"bytes"

	"github.com/majiddarvishan/fluxlog"
)

func Example() {
	var output bytes.Buffer
	log, err := fluxlog.New(fluxlog.Config{
		Level: fluxlog.InfoLevel,
		Console: &fluxlog.ConsoleConfig{
			Writer: &output,
			Format: fluxlog.JSONFormat,
		},
	})
	if err != nil {
		panic(err)
	}
	defer log.Close()

	log.Info().Str("request_id", "req-1").Msg("request received")
	_ = log.SetLevel(fluxlog.DebugLevel)
	log.Debug().Msg("runtime level changed")
}

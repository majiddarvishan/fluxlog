package goconfig

import (
	"bytes"
	"testing"

	"github.com/majiddarvishan/fluxlog"
	config "github.com/majiddarvishan/goconfig"
)

func managerFromJSON(t *testing.T, document string) (*config.Manager, *config.Node) {
	t.Helper()
	source, err := config.NewStrSource(document, `{}`)
	if err != nil {
		t.Fatal(err)
	}
	manager, err := config.NewManager(source)
	if err != nil {
		t.Fatal(err)
	}
	return manager, manager.Config()
}

func TestParseLegacyBoth(t *testing.T) {
	_, node := managerFromJSON(t, `{
		"output_mode":"both",
		"level":"debug",
		"file_name":"application.log",
		"max_file_size":25,
		"max_files":7
	}`)

	actual, levelNode, err := Parse(node)
	if err != nil {
		t.Fatal(err)
	}
	if levelNode == nil || actual.Level != fluxlog.DebugLevel {
		t.Fatalf("unexpected level: %#v", actual.Level)
	}
	if actual.Console == nil || actual.File == nil {
		t.Fatal("both mode must configure console and file outputs")
	}
	if actual.File.Path != "application.log" || actual.File.MaxSizeMB != 25 || actual.File.MaxBackups != 7 {
		t.Fatalf("unexpected file config: %#v", actual.File)
	}
	if actual.File.Format != fluxlog.JSONFormat {
		t.Fatalf("file format = %q, want json", actual.File.Format)
	}
}

func TestParseConsoleDoesNotRequireFileFields(t *testing.T) {
	_, node := managerFromJSON(t, `{"output_mode":"console","level":"info"}`)
	actual, _, err := Parse(node)
	if err != nil {
		t.Fatal(err)
	}
	if actual.Console == nil || actual.File != nil {
		t.Fatalf("unexpected outputs: %#v", actual)
	}
}

func TestParseOptionalConsoleIdentity(t *testing.T) {
	_, node := managerFromJSON(t, `{
		"output_mode":"console",
		"level":"info",
		"service":"gateway",
		"caller_max_length":15
	}`)
	actual, _, err := Parse(node)
	if err != nil {
		t.Fatal(err)
	}
	if actual.Service != "gateway" {
		t.Fatalf("service = %q, want gateway", actual.Service)
	}
	if actual.Console == nil || actual.Console.CallerMaxLength != 15 {
		t.Fatalf("unexpected console config: %#v", actual.Console)
	}
}

func TestParseRejectsInvalidOptionalConsoleIdentity(t *testing.T) {
	_, node := managerFromJSON(t, `{
		"output_mode":"console",
		"level":"info",
		"service":42
	}`)
	if _, _, err := Parse(node); err == nil {
		t.Fatal("expected invalid optional service to fail")
	}
}

func TestParseRejectsUnknownOutputMode(t *testing.T) {
	_, node := managerFromJSON(t, `{"output_mode":"remote","level":"info"}`)
	if _, _, err := Parse(node); err == nil {
		t.Fatal("expected unknown output mode to fail")
	}
}

func TestApplyLevel(t *testing.T) {
	_, root := managerFromJSON(t, `{"level":"error"}`)
	levelNode, err := root.At("level")
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	logger, err := fluxlog.New(fluxlog.Config{
		Level:   fluxlog.InfoLevel,
		Console: &fluxlog.ConsoleConfig{Writer: &output, Format: fluxlog.JSONFormat},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	if err := applyLevel(levelNode, logger); err != nil {
		t.Fatal(err)
	}
	if logger.Level() != fluxlog.ErrorLevel {
		t.Fatalf("level = %q, want error", logger.Level())
	}
}

func TestBindLevelValidatesArguments(t *testing.T) {
	if err := BindLevel(nil, nil, nil); err == nil {
		t.Fatal("expected nil manager to fail")
	}
}

func TestBindLevelRegistersNode(t *testing.T) {
	manager, root := managerFromJSON(t, `{"level":"warn"}`)
	levelNode, err := root.At("level")
	if err != nil {
		t.Fatal(err)
	}

	logger, err := fluxlog.New(fluxlog.Config{
		Level:   fluxlog.InfoLevel,
		Console: &fluxlog.ConsoleConfig{Writer: &bytes.Buffer{}, Format: fluxlog.JSONFormat},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	if err := BindLevel(manager, levelNode, logger); err != nil {
		t.Fatal(err)
	}
	if logger.Level() != fluxlog.WarnLevel {
		t.Fatalf("level = %q, want warn", logger.Level())
	}
}

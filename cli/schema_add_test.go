package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
)

func TestMakeSchemaAddCommand(t *testing.T) {
	tempDir := t.TempDir()

	store, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	schema := "type SampleType { id: String }"

	cmd := cli.MakeSchemaAddCommand()
	cmd.SetArgs([]string{schema})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	ctx := cli.WithSchemaStore(context.Background(), store)
	cmd.SetContext(ctx)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "Schema added successfully.") {
		t.Errorf("unexpected output: %s", result)
	}

	loaded, err := store.LoadCustom()
	if err != nil {
		t.Fatalf("failed to load custom schema: %v", err)
	}

	if !strings.Contains(loaded, "SampleType") {
		t.Error("expected schema type 'SampleType' not found in custom schema")
	}
}

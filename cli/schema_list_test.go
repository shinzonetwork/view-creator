package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
)

func TestMakeSchemaListCommand(t *testing.T) {
	tempDir := t.TempDir()

	store, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	customSchema := "type CustomListType { name: String }"
	if err := store.SaveCustom(customSchema); err != nil {
		t.Fatalf("failed to save custom schema: %v", err)
	}

	cmd := cli.MakeSchemaListCommand()
	cmd.SetArgs([]string{})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	ctx := cli.WithSchemaStore(context.Background(), store)
	cmd.SetContext(ctx)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "Schemas:") {
		t.Error("expected 'Schemas:' header in output")
	}
	if !strings.Contains(result, "CustomListType (custom)") {
		t.Errorf("expected custom type 'CustomListType' in output, got:\n%s", result)
	}
}

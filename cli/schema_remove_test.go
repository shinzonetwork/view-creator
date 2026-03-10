package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
)

func TestMakeSchemaRemoveCommand(t *testing.T) {
	tempDir := t.TempDir()

	store, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	schema := "type ToRemove { key: String }"
	if err := store.SaveCustom(schema); err != nil {
		t.Fatalf("failed to save custom schema: %v", err)
	}

	cmd := cli.MakeSchemaRemoveCommand()
	cmd.SetArgs([]string{"ToRemove"})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	ctx := cli.WithSchemaStore(context.Background(), store)
	cmd.SetContext(ctx)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "Schema removed.") {
		t.Errorf("unexpected output: %s", result)
	}

	custom, err := store.LoadCustom()
	if err != nil {
		t.Fatalf("failed to reload custom schema: %v", err)
	}
	if strings.Contains(custom, "ToRemove") {
		t.Error("type 'ToRemove' should have been removed from the custom schema")
	}
}

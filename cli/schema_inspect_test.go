package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
)

func TestMakeSchemaInspectCommand(t *testing.T) {
	tempDir := t.TempDir()

	store, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	schema := "type InspectableType { id: String }"
	if err := store.SaveCustom(schema); err != nil {
		t.Fatalf("failed to add custom schema: %v", err)
	}

	cmd := cli.MakeSchemaInspectCommand()
	cmd.SetArgs([]string{"InspectableType"})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	ctx := cli.WithSchemaStore(context.Background(), store)
	cmd.SetContext(ctx)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "type InspectableType") {
		t.Errorf("unexpected output: %s", result)
	}
}

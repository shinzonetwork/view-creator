package cli_test

import (
	"context"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
)

func TestMakeSchemaUpdateCommand(t *testing.T) {
	t.Skip("TODO: implement test for MakeSchemaUpdateCommand")

	tempDir := t.TempDir()

	store, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	cmd := cli.MakeSchemaUpdateCommand()
	cmd.SetArgs([]string{"--version", "main"})

	ctx := cli.WithSchemaStore(context.Background(), store)
	cmd.SetContext(ctx)

	// TODO: stub HTTP call or test against real GitHub branch (but url is not public)
	// var out bytes.Buffer
	// cmd.SetOut(&out)
	// cmd.SetErr(&out)

	// err = cmd.Execute()
	// if err != nil {
	//     t.Fatalf("command execution failed: %v", err)
	// }

	// TODO: assert output contains expected update message
	// TODO: assert schema file was updated
}

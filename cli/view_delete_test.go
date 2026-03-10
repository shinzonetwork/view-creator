package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
	"github.com/spf13/cobra"
)

func TestDeleteViewSuccess(t *testing.T) {
	tempDir := t.TempDir()
	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create local store: %v", err)
	}

	// First: create a view so we can delete it
	view, err := service.InitView("testview", store)
	if err != nil {
		t.Fatalf("failed to init view: %v", err)
	}
	if view.Name != "testview" {
		t.Fatalf("unexpected view name: %s", view.Name)
	}

	// Now: delete the view via CLI
	cmd := cli.MakeViewDeleteCommand()
	cmd.SetArgs([]string{"testview"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	if err := cmd.Execute(); err != nil {
		t.Fatalf("delete command failed: %v", err)
	}

	out := buf.String()
	expected := "deleted view testview Successfully"
	if !strings.Contains(out, expected) {
		t.Errorf("expected confirmation message, got:\n%s", out)
	}
}

func TestDeleteViewAlreadyDeletedFails(t *testing.T) {
	tempDir := t.TempDir()
	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create local store: %v", err)
	}

	// Create and delete the view first
	if _, err := service.InitView("testview", store); err != nil {
		t.Fatalf("failed to init view: %v", err)
	}

	deleteCmd := func() *cobra.Command {
		cmd := cli.MakeViewDeleteCommand()
		cmd.SetArgs([]string{"testview"})
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetContext(cli.WithViewStore(context.Background(), store))
		return cmd
	}

	// First delete (should succeed)
	if err := deleteCmd().Execute(); err != nil {
		t.Fatalf("first delete failed: %v", err)
	}

	// Second delete (should fail)
	err = deleteCmd().Execute()
	if err == nil {
		t.Fatal("expected error on second delete, got none")
	}

	if !strings.Contains(err.Error(), "view does not exists") {
		t.Errorf("expected error to mention 'not found', got: %v", err)
	}
}

package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
	"github.com/spf13/cobra"
)

func TestInitViewDirectWithTempStore(t *testing.T) {
	tempDir := t.TempDir()

	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create temp store: %v", err)
	}

	cmd := cli.MakeViewInitCommand()

	cmd.SetArgs([]string{"testview"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Inject store into context
	ctx := cli.WithViewStore(context.Background(), store)
	cmd.SetContext(ctx)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	out := buf.String()

	expected := `📄 View: testview
🔍 Query: <none>
📐 SDL: <none>
🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: `

	if !strings.Contains(out, expected) {
		t.Errorf("expected view name in output, got:\n%s \nexpected:\n%s", out, expected)
	}
}

func TestInitViewDuplicateFails(t *testing.T) {
	tempDir := t.TempDir()
	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create temp store: %v", err)
	}

	createCmd := func() *cobra.Command {
		cmd := cli.MakeViewInitCommand()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetContext(cli.WithViewStore(context.Background(), store))
		return cmd
	}

	// First creation (should succeed)
	cmd1 := createCmd()
	cmd1.SetArgs([]string{"testview"})

	if err := cmd1.Execute(); err != nil {
		t.Fatalf("first creation failed unexpectedly: %v", err)
	}

	// Second creation (should fail)
	cmd2 := createCmd()
	cmd2.SetArgs([]string{"testview"})

	err = cmd2.Execute()
	if err == nil {
		t.Fatal("expected second creation to fail, but it succeeded")
	}

	// Optionally check for specific error message
	expected := "view already exists"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("expected error to contain %q, got: %v", expected, err)
	}
}

package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestViewInspectCommandSuccess(t *testing.T) {
	tempDir := t.TempDir()
	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create local store: %v", err)
	}

	// Create the view first
	_, err = service.InitView("testview", store)
	if err != nil {
		t.Fatalf("failed to initialize view: %v", err)
	}

	// Create the inspect command
	cmd := cli.MakeViewInspectCommand()
	cmd.SetArgs([]string{"testview"})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("inspect command failed: %v", err)
	}

	output := out.String()

	expected := `📄 View: testview
🔍 Query: <none>
📐 SDL: <none>
🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: `

	if !strings.Contains(output, expected) {
		t.Errorf("expected view name in output, got:\n%s \nexpected:\n%s", output, expected)
	}
}

func TestViewInspectCommandNotFound(t *testing.T) {
	tempDir := t.TempDir()
	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create local store: %v", err)
	}

	cmd := cli.MakeViewInspectCommand()
	cmd.SetArgs([]string{"ghostview"})

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err == nil {
		t.Fatal("expected error when inspecting non-existent view, got nil")
	}

	if !strings.Contains(err.Error(), "view does not exists") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

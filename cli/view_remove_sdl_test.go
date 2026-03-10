package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestRemoveSdlFromView(t *testing.T) {
	tempDir := t.TempDir()

	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create temp store: %v", err)
	}

	viewName := "testview"

	cmd := cli.MakeViewInitCommand()
	cmd.SetArgs([]string{viewName})

	var initBuf bytes.Buffer
	cmd.SetOut(&initBuf)
	cmd.SetErr(&initBuf)
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	if err := cmd.Execute(); err != nil {
		t.Fatalf("view init command failed: %v", err)
	}

	sdl := `type ReturnedLog { address: String }`
	cmd = cli.MakeAddSdlCommand(&viewName)

	var addBuf bytes.Buffer
	cmd.SetOut(&addBuf)
	cmd.SetErr(&addBuf)
	cmd.SetArgs([]string{sdl})
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	if err := cmd.Execute(); err != nil {
		t.Fatalf("add SDL command failed: %v", err)
	}

	cmd = cli.MakeRemoveSdlCommand(&viewName)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	if err := cmd.Execute(); err != nil {
		t.Fatalf("remove SDL command failed: %v", err)
	}

	out := buf.String()

	expected := `📄 View: testview
🔍 Query: <none>
📐 SDL: <none>
🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: `

	if !strings.HasPrefix(out, expected) {
		t.Errorf("unexpected output after SDL removal.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

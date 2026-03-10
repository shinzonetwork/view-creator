package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestAddSdlToExistingView(t *testing.T) {
	tempDir := t.TempDir()

	// Create local store
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

	// Inject store into context
	ctx := cli.WithViewStore(context.Background(), store)
	cmd.SetContext(ctx)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	sdl := `type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String}`

	cmd = cli.MakeAddSdlCommand(&viewName)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{sdl})
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("add SDL command failed: %v", err)
	}

	out := buf.String()

	expected := `📄 View: testview
🔍 Query: <none>
📐 SDL:
type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String}

🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: `

	if !strings.HasPrefix(out, expected) {
		t.Errorf("unexpected output.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

func TestUpdateSdlOfExistingView(t *testing.T) {
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

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	initialSDL := `type FilteredAndDecodedLogs @materialized(if: false) {hash: String}`
	cmd = cli.MakeAddSdlCommand(&viewName)

	var firstBuf bytes.Buffer
	cmd.SetOut(&firstBuf)
	cmd.SetErr(&firstBuf)
	cmd.SetArgs([]string{initialSDL})
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("add initial SDL command failed: %v", err)
	}

	updatedSDL := `type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String}`

	cmd = cli.MakeAddSdlCommand(&viewName)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{updatedSDL})
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("update SDL command failed: %v", err)
	}

	out := buf.String()

	expected := `📄 View: testview
🔍 Query: <none>
📐 SDL:
type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String}

🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: `

	if !strings.HasPrefix(out, expected) {
		t.Errorf("unexpected output.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

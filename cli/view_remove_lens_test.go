package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestRemoveLensFromView(t *testing.T) {
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
		t.Fatalf("view init command failed: %v", err)
	}

	wasmPath := filepath.Join(tempDir, "decode.wasm")
	if err := os.WriteFile(wasmPath, []byte("\x00asm\x01\x00\x00\x00"), 0644); err != nil {
		t.Fatalf("failed to write dummy wasm file: %v", err)
	}

	cmd = cli.MakeAddLensCommand(&viewName)

	var addBuf bytes.Buffer
	cmd.SetOut(&addBuf)
	cmd.SetErr(&addBuf)
	cmd.SetArgs([]string{
		"--label", "decode_usdt",
		"--path", wasmPath,
		"--args", `{"token": "USDT"}`,
	})
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("add lens command failed: %v", err)
	}

	cmd = cli.MakeRemoveLensCommand(&viewName)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--label", "decode_usdt"})
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("remove lens command failed: %v", err)
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
		t.Errorf("unexpected output after lens removal.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

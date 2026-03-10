package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/cli"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestRemoveQueryFromView(t *testing.T) {
	tempDir := t.TempDir()

	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create temp store: %v", err)
	}

	schemastore, err := fileschema.NewFileSchemaStore(tempDir)
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

	query := `Log { address topics }`
	cmd = cli.MakeAddQueryCommand(&viewName)

	var addBuf bytes.Buffer
	cmd.SetOut(&addBuf)
	cmd.SetErr(&addBuf)
	cmd.SetArgs([]string{query})
	ctx := context.Background()
	ctx = cli.WithViewStore(ctx, store)
	ctx = cli.WithSchemaStore(ctx, schemastore)
	cmd.SetContext(ctx)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("add query command failed: %v", err)
	}

	cmd = cli.MakeRemoveQueryCommand(&viewName)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetContext(cli.WithViewStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("remove query command failed: %v", err)
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
		t.Errorf("unexpected output after query removal.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

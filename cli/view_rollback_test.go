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

func TestRollbackView(t *testing.T) {
	tempDir := t.TempDir()

	store, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	schemaStore, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	viewName := "testrollback"

	initCmd := cli.MakeViewInitCommand()
	initCmd.SetArgs([]string{viewName})
	initCmd.SetContext(cli.WithViewStore(context.Background(), store))
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	initialQuery := `Log { address }`
	queryCmd := cli.MakeAddQueryCommand(&viewName)
	var firstBuf bytes.Buffer
	queryCmd.SetArgs([]string{initialQuery})
	queryCmd.SetOut(&firstBuf)
	queryCmd.SetContext(cli.WithSchemaStore(cli.WithViewStore(context.Background(), store), schemaStore))
	if err := queryCmd.Execute(); err != nil {
		t.Fatalf("add initial query failed: %v", err)
	}

	updatedQuery := `Log { address topics }`
	queryCmd = cli.MakeAddQueryCommand(&viewName)
	var secondBuf bytes.Buffer
	queryCmd.SetArgs([]string{updatedQuery})
	queryCmd.SetOut(&secondBuf)
	queryCmd.SetContext(cli.WithSchemaStore(cli.WithViewStore(context.Background(), store), schemaStore))
	if err := queryCmd.Execute(); err != nil {
		t.Fatalf("update query failed: %v", err)
	}

	rollbackCmd := cli.MakeViewRollbackCommand()
	rollbackCmd.SetArgs([]string{viewName, "--version=1"})
	var rollbackBuf bytes.Buffer
	rollbackCmd.SetOut(&rollbackBuf)
	rollbackCmd.SetContext(cli.WithViewStore(context.Background(), store))
	if err := rollbackCmd.Execute(); err != nil {
		t.Fatalf("rollback command failed: %v", err)
	}

	out := rollbackBuf.String()
	expected := `📄 View: testrollback
🔍 Query:
Log { address }

📐 SDL: <none>
🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version:`

	if !strings.HasPrefix(out, expected) {
		t.Errorf("unexpected rollback output.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

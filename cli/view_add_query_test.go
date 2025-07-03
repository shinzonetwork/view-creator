package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shinzonetwork/view-creator/cli"
	"github.com/shinzonetwork/view-creator/core/store/local"
)

func TestAddQueryToExistingView(t *testing.T) {
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
	cmd.SetContext(cli.WithStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	query := `Log {address topics data transactionHash blockNumber}`

	cmd = cli.MakeAddQueryCommand(&viewName)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{query})
	cmd.SetContext(cli.WithStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("add query command failed: %v", err)
	}

	out := buf.String()

	expected := `📄 View: testview
🔍 Query:
Log {address topics data transactionHash blockNumber}

📐 SDL: <none>
🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: 0
 - Total: 0
 - Created At: `

	if !strings.HasPrefix(out, expected) {
		t.Errorf("unexpected output.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

func TestUpdateQueryOfExistingView(t *testing.T) {
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
	cmd.SetContext(cli.WithStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	initialQuery := `Log {address}`
	cmd = cli.MakeAddQueryCommand(&viewName)

	var firstBuf bytes.Buffer
	cmd.SetOut(&firstBuf)
	cmd.SetErr(&firstBuf)
	cmd.SetArgs([]string{initialQuery})
	cmd.SetContext(cli.WithStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("add initial query command failed: %v", err)
	}

	updatedQuery := `Log {address topics data transactionHash blockNumber}`
	cmd = cli.MakeAddQueryCommand(&viewName)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{updatedQuery})
	cmd.SetContext(cli.WithStore(context.Background(), store))

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("update query command failed: %v", err)
	}

	out := buf.String()

	expected := `📄 View: testview
🔍 Query:
Log {address topics data transactionHash blockNumber}

📐 SDL: <none>
🔧 Lenses:
 - (empty)

🗂  Metadata:
 - Version: 0
 - Total: 0
 - Created At: `

	if !strings.HasPrefix(out, expected) {
		t.Errorf("unexpected output.\nGot:\n%s\nExpected prefix:\n%s", out, expected)
	}
}

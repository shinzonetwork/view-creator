package fileschema_test

import (
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
)

func TestFileSchemaStore_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()

	s, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	defaultSchema, err := s.LoadDefault()
	if err != nil {
		t.Fatalf("failed to load default schema: %v", err)
	}
	if len(defaultSchema) == 0 {
		t.Error("default schema should not be empty")
	}

	if err := s.SaveCustom("type TestType { field: String }"); err != nil {
		t.Fatalf("failed to save custom schema: %v", err)
	}

	customSchema, err := s.LoadCustom()
	if err != nil {
		t.Fatalf("failed to load custom schema: %v", err)
	}
	if !strings.Contains(customSchema, "TestType") {
		t.Error("custom schema should contain TestType")
	}

	defaultTypes, customTypes, err := s.ListTypes()
	if err != nil {
		t.Fatalf("failed to list types: %v", err)
	}
	if len(defaultTypes) == 0 {
		t.Error("expected at least one default type")
	}
	found := false
	for _, t := range customTypes {
		if t == "TestType" {
			found = true
			break
		}
	}
	if !found {
		t.Error("TestType should appear in custom types list")
	}

	def, err := s.GetTypeDefinition("TestType")
	if err != nil {
		t.Fatalf("failed to get type definition: %v", err)
	}
	if !strings.Contains(def, "field: String") {
		t.Error("definition should include field definition")
	}

	if err := s.ResetCustom(); err != nil {
		t.Fatalf("failed to reset custom schema: %v", err)
	}

	afterReset, err := s.LoadCustom()
	if err != nil {
		t.Fatalf("failed to reload custom schema: %v", err)
	}
	if strings.TrimSpace(afterReset) != "" {
		t.Error("custom schema should be empty after reset")
	}
}

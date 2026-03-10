package service_test

import (
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
)

func TestSchemaService_AddGetListRemoveReset(t *testing.T) {
	tempDir := t.TempDir()

	schemaStore, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	sdl := "type TestServiceType { id: String }"

	err = service.AddCustomSchema(schemaStore, sdl)
	if err != nil {
		t.Fatalf("AddCustomSchema failed: %v", err)
	}

	def, err := service.GetSchemaTypeDefinition(schemaStore, "TestServiceType")
	if err != nil {
		t.Fatalf("GetSchemaTypeDefinition failed: %v", err)
	}
	if !strings.Contains(def, "id: String") {
		t.Errorf("unexpected definition output: %s", def)
	}

	_, customs, err := service.ListSchemas(schemaStore)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	found := false
	for _, c := range customs {
		if c == "TestServiceType" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("custom type 'TestServiceType' not found in list")
	}

	err = service.RemoveCustomSchema(schemaStore, "TestServiceType")
	if err != nil {
		t.Fatalf("RemoveCustomSchema failed: %v", err)
	}

	customAfter, err := schemaStore.LoadCustom()
	if err != nil {
		t.Fatalf("LoadCustom after removal failed: %v", err)
	}
	if strings.Contains(customAfter, "TestServiceType") {
		t.Error("expected type to be removed from custom schema")
	}

	err = service.AddCustomSchema(schemaStore, sdl)
	if err != nil {
		t.Fatalf("re-adding schema before reset failed: %v", err)
	}

	err = service.ResetCustomSchemas(schemaStore)
	if err != nil {
		t.Fatalf("ResetCustomSchemas failed: %v", err)
	}

	cleared, err := schemaStore.LoadCustom()
	if err != nil {
		t.Fatalf("LoadCustom after reset failed: %v", err)
	}
	if strings.TrimSpace(cleared) != "" {
		t.Error("expected custom schema to be empty after reset")
	}
}

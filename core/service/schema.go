package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	"github.com/shinzonetwork/shinzo-view-creator/core/util"
)

func AddCustomSchema(schemaStore store.SchemaStore, newSchema string) error {
	customContent, err := schemaStore.LoadCustom()
	if err != nil {
		return fmt.Errorf("failed to load custom schema: %w", err)
	}

	typeName, err := extractStrictTypeName(newSchema)
	if err != nil {
		return err
	}

	if strings.Contains(customContent, fmt.Sprintf("type %s", typeName)) {
		return fmt.Errorf("type %s already exists", typeName)
	}

	if err := util.ValidateSchemaBlock(newSchema); err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	customContent = strings.TrimSpace(customContent) + "\n\n" + strings.TrimSpace(newSchema) + "\n"

	if err := schemaStore.SaveCustom(customContent); err != nil {
		return fmt.Errorf("failed to save updated schema: %w", err)
	}

	return nil
}

func GetSchemaTypeDefinition(schemaStore store.SchemaStore, name string) (string, error) {
	return schemaStore.GetTypeDefinition(name)
}

func ListSchemas(schemaStore store.SchemaStore) (defaultTypes []string, customTypes []string, err error) {
	return schemaStore.ListTypes()
}

func RemoveCustomSchema(schemaStore store.SchemaStore, name string) error {
	custom, err := schemaStore.LoadCustom()
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?ms)^type\s+` + regexp.QuoteMeta(name) + `\s*\{[^}]*\}\n*`)
	modified := re.ReplaceAllString(custom, "")
	if modified == custom {
		return fmt.Errorf("type %s not found in custom schema", name)
	}

	return schemaStore.SaveCustom(modified)
}

func ResetCustomSchemas(schemaStore store.SchemaStore) error {
	return schemaStore.ResetCustom()
}

func UpdateDefaultSchemas(schemaStore store.SchemaStore, version string) error {
	if version == "" {
		version = "main"
	}
	return schemaStore.UpdateDefaultFromRemote(version)
}

func extractStrictTypeName(schema string) (string, error) {
	re := regexp.MustCompile(`(?m)^\s*type\s+(\w+)`)
	match := re.FindStringSubmatch(schema)
	if len(match) < 2 {
		return "", fmt.Errorf("only 'type <Name> { ... }' is supported")
	}
	return match[1], nil
}

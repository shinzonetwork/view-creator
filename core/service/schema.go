package service

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shinzonetwork/view-creator/core/util"
)

const customMarker = "# --- CUSTOM SCHEMAS ---"
const defaultMarker = "# --- DEFAULT SCHEMAS ---"

func AddCustomSchema(schema string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	schemaPath := filepath.Join(home, ".shinzo", "tools", "schema.graphql")

	content := ""
	if b, err := os.ReadFile(schemaPath); err == nil {
		content = string(b)
	} else {
		content = fmt.Sprintf("%s\n\n%s\n", defaultMarker, customMarker)
	}

	typeName, err := extractStrictTypeName(schema)
	if err != nil {
		return err
	}

	if strings.Contains(content, fmt.Sprintf("type %s", typeName)) {
		return fmt.Errorf("type %s already exists", typeName)
	}

	if err := util.ValidateSchemaBlock(schema); err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	parts := strings.Split(content, customMarker)
	if len(parts) != 2 {
		return fmt.Errorf("malformed schema file: missing custom section")
	}

	before := strings.TrimRight(parts[0], "\n")
	after := strings.TrimLeft(parts[1], "\n")

	newBlock := fmt.Sprintf("\n\n%s\n\n%s\n", customMarker, strings.TrimSpace(schema))
	newContent := before + newBlock + after

	if err := os.MkdirAll(filepath.Dir(schemaPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(schemaPath, []byte(newContent), 0644)
}

func RemoveCustomSchema(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	schemaPath := filepath.Join(home, ".shinzo", "tools", "schema.graphql")

	contentBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("could not read schema file: %w", err)
	}
	content := string(contentBytes)

	parts := strings.Split(content, customMarker)
	if len(parts) != 2 {
		return fmt.Errorf("schema file is malformed or missing custom section")
	}

	before := parts[0]
	customBlock := parts[1]

	re := regexp.MustCompile(`(?ms)^type\s+` + regexp.QuoteMeta(name) + `\s*\{[^}]*\}\n*`)
	modified := re.ReplaceAllString(customBlock, "")

	if modified == customBlock {
		return fmt.Errorf("type %s not found in custom schemas", name)
	}

	newContent := before + customMarker + "\n" + strings.TrimLeft(modified, "\n")

	return os.WriteFile(schemaPath, []byte(newContent), 0644)
}

func ListSchemas() ([]string, []string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}
	schemaPath := filepath.Join(home, ".shinzo", "tools", "schema.graphql")

	contentBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read schema file: %w", err)
	}
	content := string(contentBytes)

	parts := strings.Split(content, customMarker)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("schema file is malformed")
	}

	defaultBlock := parts[0]
	customBlock := parts[1]

	typeRegex := regexp.MustCompile(`(?m)^type\s+(\w+)`)

	defaults := []string{}
	for _, match := range typeRegex.FindAllStringSubmatch(defaultBlock, -1) {
		defaults = append(defaults, match[1])
	}

	customs := []string{}
	for _, match := range typeRegex.FindAllStringSubmatch(customBlock, -1) {
		customs = append(customs, match[1])
	}

	return defaults, customs, nil
}

func GetSchemaTypeDefinition(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	schemaPath := filepath.Join(home, ".shinzo", "tools", "schema.graphql")

	b, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("could not read schema: %w", err)
	}
	content := string(b)

	re := regexp.MustCompile(`(?ms)^type\s+` + regexp.QuoteMeta(name) + `\s*\{[^}]*\}`)
	match := re.FindString(content)

	if match == "" {
		return "", fmt.Errorf("type '%s' not found in schema", name)
	}
	return match, nil
}

func UpdateDefaultSchemas() error {
	// WIP

	// TODO: we need to make this public
	// const remoteURL = "https://raw.githubusercontent.com/shinzonetwork/viewkit/main/tools/default_schema.graphql"

	return nil
}

func extractStrictTypeName(schema string) (string, error) {
	re := regexp.MustCompile(`(?m)^\s*type\s+(\w+)`)
	match := re.FindStringSubmatch(schema)
	if len(match) < 2 {
		return "", fmt.Errorf("only 'type <Name> { ... }' is supported")
	}
	return match[1], nil
}

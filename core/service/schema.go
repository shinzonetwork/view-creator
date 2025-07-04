package service

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shinzonetwork/view-creator/core/util"
)

func AddCustomSchema(schema string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	customPath := filepath.Join(home, ".shinzo", "tools", "schema", "custom_schema.graphql")

	content := ""
	if b, err := os.ReadFile(customPath); err == nil {
		content = string(b)
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

	content = strings.TrimSpace(content) + "\n\n" + strings.TrimSpace(schema) + "\n"
	return os.WriteFile(customPath, []byte(content), 0644)
}

func RemoveCustomSchema(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	customPath := filepath.Join(home, ".shinzo", "tools", "schema", "custom_schema.graphql")

	contentBytes, err := os.ReadFile(customPath)
	if err != nil {
		return fmt.Errorf("could not read custom schema file: %w", err)
	}
	content := string(contentBytes)

	re := regexp.MustCompile(`(?ms)^type\s+` + regexp.QuoteMeta(name) + `\s*\{[^}]*\}\n*`)
	modified := re.ReplaceAllString(content, "")

	if modified == content {
		return fmt.Errorf("type %s not found in custom schemas", name)
	}

	return os.WriteFile(customPath, []byte(strings.TrimSpace(modified)+"\n"), 0644)
}

func ListSchemas() ([]string, []string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}

	defaultPath := filepath.Join(home, ".shinzo", "tools", "schema", "default_schema.graphql")
	customPath := filepath.Join(home, ".shinzo", "tools", "schema", "custom_schema.graphql")

	defaultBlock, _ := os.ReadFile(defaultPath)
	customBlock, _ := os.ReadFile(customPath)

	typeRegex := regexp.MustCompile(`(?m)^type\s+(\w+)`)

	var defaults, customs []string

	for _, match := range typeRegex.FindAllStringSubmatch(string(defaultBlock), -1) {
		defaults = append(defaults, match[1])
	}
	for _, match := range typeRegex.FindAllStringSubmatch(string(customBlock), -1) {
		customs = append(customs, match[1])
	}

	return defaults, customs, nil
}

func GetSchemaTypeDefinition(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	paths := []string{
		filepath.Join(home, ".shinzo", "tools", "schema", "default_schema.graphql"),
		filepath.Join(home, ".shinzo", "tools", "schema", "custom_schema.graphql"),
	}

	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		content := string(b)
		re := regexp.MustCompile(`(?ms)^type\s+` + regexp.QuoteMeta(name) + `\s*\{([^}]*)\}`)
		matches := re.FindStringSubmatch(content)

		if len(matches) >= 2 {
			formatted := formatTypeBlock(name, matches[1])
			return formatted, nil
		}
	}

	return "", fmt.Errorf("type '%s' not found in schema", name)
}

func ResetCustomSchemas() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	customPath := filepath.Join(home, ".shinzo", "tools", "schema", "custom_schema.graphql")
	return os.WriteFile(customPath, []byte(""), 0644)
}

func formatTypeBlock(typeName, rawBody string) string {
	rawBody = strings.TrimSpace(rawBody)
	lines := strings.Split(rawBody, "\n")

	for i, line := range lines {
		lines[i] = "  " + strings.TrimSpace(line)
	}

	return fmt.Sprintf("type %s {\n%s\n}", typeName, strings.Join(lines, "\n"))
}

func extractStrictTypeName(schema string) (string, error) {
	re := regexp.MustCompile(`(?m)^\s*type\s+(\w+)`)
	match := re.FindStringSubmatch(schema)
	if len(match) < 2 {
		return "", fmt.Errorf("only 'type <Name> { ... }' is supported")
	}
	return match[1], nil
}

func UpdateDefaultSchemas() error {
	// WIP

	// TODO: we need to make this public
	// const remoteURL = "https://raw.githubusercontent.com/shinzonetwork/viewkit/main/tools/default_schema.graphql"

	return nil
}

package fileschema

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shinzonetwork/shinzo-view-creator/tools"
)

type FileSchemaStore struct {
	BasePath string
}

func NewFileSchemaStore(dir ...string) (*FileSchemaStore, error) {
	var base string

	if len(dir) == 0 || dir[0] == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("unable to get home directory: %w", err)
		}
		base = filepath.Join(home, ".shinzo", "schema")
	} else {
		base = filepath.Join(dir[0], ".shinzo", "schema")
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, fmt.Errorf("unable to create schema directory: %w", err)
	}

	defaultPath := filepath.Join(base, "default_schema.graphql")
	customPath := filepath.Join(base, "custom_schema.graphql")

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) || isFileEmpty(defaultPath) {
		if err := os.WriteFile(defaultPath, []byte(strings.TrimSpace(tools.DefaultSchema)), 0644); err != nil {
			return nil, fmt.Errorf("failed to write default schema: %w", err)
		}
	}

	if _, err := os.Stat(customPath); os.IsNotExist(err) || isFileEmpty(customPath) {
		if err := os.WriteFile(customPath, []byte(""), 0644); err != nil {
			return nil, fmt.Errorf("failed to write custom schema: %w", err)
		}
	}

	return &FileSchemaStore{BasePath: base}, nil
}

func (s *FileSchemaStore) Load() (string, error) {
	defaultSchema, err := s.LoadDefault()
	if err != nil {
		return "", fmt.Errorf("failed to load default schema: %w", err)
	}

	customSchema, err := s.LoadCustom()
	if err != nil {
		return "", fmt.Errorf("failed to load custom schema: %w", err)
	}

	return defaultSchema + "\n\n" + customSchema, nil
}

func (s *FileSchemaStore) LoadDefault() (string, error) {
	return read(filepath.Join(s.BasePath, "default_schema.graphql"))
}

func (s *FileSchemaStore) LoadCustom() (string, error) {
	return read(filepath.Join(s.BasePath, "custom_schema.graphql"))
}

func (s *FileSchemaStore) SaveCustom(schema string) error {
	customPath := filepath.Join(s.BasePath, "custom_schema.graphql")
	tempPath := customPath + ".tmp"

	content := strings.TrimSpace(schema) + "\n"

	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write temp schema: %w", err)
	}

	if err := os.Rename(tempPath, customPath); err != nil {
		return fmt.Errorf("failed to replace schema file: %w", err)
	}

	return nil
}

func (s *FileSchemaStore) ResetCustom() error {
	customPath := filepath.Join(s.BasePath, "custom_schema.graphql")
	tempPath := customPath + ".tmp"

	if err := os.WriteFile(tempPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to write temp empty schema: %w", err)
	}

	if err := os.Rename(tempPath, customPath); err != nil {
		return fmt.Errorf("failed to clear custom schema: %w", err)
	}

	return nil
}

func (s *FileSchemaStore) ListTypes() ([]string, []string, error) {
	defaultSchema, err := s.LoadDefault()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load default schema: %w", err)
	}

	customSchema, err := s.LoadCustom()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load custom schema: %w", err)
	}

	typeRegex := regexp.MustCompile(`(?m)^type\s+(\w+)`)

	var defaultTypes, customTypes []string

	for _, match := range typeRegex.FindAllStringSubmatch(defaultSchema, -1) {
		defaultTypes = append(defaultTypes, match[1])
	}
	for _, match := range typeRegex.FindAllStringSubmatch(customSchema, -1) {
		customTypes = append(customTypes, match[1])
	}

	return defaultTypes, customTypes, nil
}

func (s *FileSchemaStore) GetTypeDefinition(typeName string) (string, error) {
	paths := []string{
		filepath.Join(s.BasePath, "default_schema.graphql"),
		filepath.Join(s.BasePath, "custom_schema.graphql"),
	}

	re := regexp.MustCompile(`(?ms)^type\s+` + regexp.QuoteMeta(typeName) + `\s*\{([^}]*)\}`)

	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		content := string(b)
		matches := re.FindStringSubmatch(content)

		if len(matches) >= 2 {
			return formatTypeBlock(typeName, matches[1]), nil
		}
	}

	return "", fmt.Errorf("type '%s' not found in schema", typeName)
}

func (s *FileSchemaStore) UpdateDefaultFromRemote(version string) error {
	if version == "" {
		version = "main"
	}

	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/shinzonetwork/viewkit/%s/tools/default_schema.graphql",
		version,
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch schema from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response from %s: %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read schema response: %w", err)
	}

	path := filepath.Join(s.BasePath, "default_schema.graphql")
	temp := path + ".tmp"

	// Safe write using temp file
	if err := os.WriteFile(temp, []byte(strings.TrimSpace(string(body))+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := os.Rename(temp, path); err != nil {
		return fmt.Errorf("failed to replace schema file: %w", err)
	}

	return nil
}

func isFileEmpty(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Size() == 0
}

func read(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func formatTypeBlock(typeName, rawBody string) string {
	rawBody = strings.TrimSpace(rawBody)
	lines := strings.Split(rawBody, "\n")

	for i, line := range lines {
		lines[i] = "  " + strings.TrimSpace(line)
	}

	return fmt.Sprintf("type %s {\n%s\n}", typeName, strings.Join(lines, "\n"))
}

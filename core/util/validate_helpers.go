package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

func ValidateQueryAgainstSchema(rawQuery string) error {
	wrappedQuery := fmt.Sprintf("query { %s }", strings.TrimSpace(rawQuery))

	home, _ := os.UserHomeDir()
	schemaPath := filepath.Join(home, ".shinzo", "tools", "schema.graphql")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}

	finalSchema := buildSchemaWithRoot(string(schemaBytes))

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "schema.graphql",
		Input: finalSchema,
	})
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	queryDoc, err := parser.ParseQuery(&ast.Source{
		Name:  "query.graphql",
		Input: wrappedQuery,
	})
	if err != nil {
		return fmt.Errorf("query parse error: %w", err)
	}

	errs := validator.ValidateWithRules(schema, queryDoc, nil)
	if len(errs) > 0 {
		return formatValidationErrors(errs)
	}

	return nil
}

func buildSchemaWithRoot(original string) string {
	var rootFields []string

	lines := strings.Split(original, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "type ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				typeName := parts[1]
				if typeName != "Query" {
					rootFields = append(rootFields, fmt.Sprintf("  %s: [%s]", typeName, typeName))
				}
			}
		}
	}

	rootType := "type Query {\n" + strings.Join(rootFields, "\n") + "\n}"
	schemaBlock := "schema {\n  query: Query\n}"

	return strings.TrimSpace(original) + "\n\n" + schemaBlock + "\n\n" + rootType
}

func formatValidationErrors(errs gqlerror.List) error {
	if len(errs) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("GraphQL validation failed:\n")
	for _, err := range errs {
		b.WriteString(fmt.Sprintf("• %s\n", err.Error()))
	}
	return fmt.Errorf("%s", b.String())
}

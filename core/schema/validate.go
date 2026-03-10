package schema

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	"github.com/vektah/gqlparser/v2/validator/rules"
)

func ValidateQuery(schemaStore store.SchemaStore, rawQuery string) error {
	defaultSchema, err := schemaStore.LoadDefault()
	if err != nil {
		return fmt.Errorf("failed to load default schema: %w", err)
	}

	customSchema, err := schemaStore.LoadCustom()
	if err != nil {
		return fmt.Errorf("failed to load custom schema: %w", err)
	}

	fullSchema := buildSchemaWithRoot(defaultSchema + "\n\n" + customSchema)

	schemaAST, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "combined.graphql",
		Input: fullSchema,
	})
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	wrappedQuery := fmt.Sprintf("query { %s }", strings.TrimSpace(rawQuery))

	queryDoc, err := parser.ParseQuery(&ast.Source{
		Name:  "query.graphql",
		Input: wrappedQuery,
	})
	if err != nil {
		return fmt.Errorf("query parse error: %w", err)
	}

	errs := validator.ValidateWithRules(schemaAST, queryDoc, rules.NewDefaultRules())
	if len(errs) > 0 {
		var lines []string
		for _, e := range errs {
			lines = append(lines, fmt.Sprintf("• %s", e.Error()))
		}
		return fmt.Errorf("GraphQL validation failed:\n%s", strings.Join(lines, "\n"))
	}

	return nil
}

func buildSchemaWithRoot(original string) string {
	re := regexp.MustCompile(`(?m)^type\s+(\w+)\s*\{`)
	matches := re.FindAllStringSubmatch(original, -1)

	var rootFields []string
	for _, match := range matches {
		typeName := match[1]
		if typeName == "Query" {
			continue
		}
		rootFields = append(rootFields, fmt.Sprintf("  %s: [%s]", typeName, typeName))
	}

	rootType := "type Query {\n" + strings.Join(rootFields, "\n") + "\n}"
	schemaBlock := "schema {\n  query: Query\n}"

	return strings.TrimSpace(original) + "\n\n" + schemaBlock + "\n\n" + rootType
}

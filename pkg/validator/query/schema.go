package query

import (
	"fmt"
	"os"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

var loadedSchema *ast.Schema

func LoadSchema(path string) error {
	schemaBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}
	source := &ast.Source{Input: string(schemaBytes), Name: "schema.graphql"}
	schema, gqlErr := gqlparser.LoadSchema(source)
	if gqlErr != nil {
		return fmt.Errorf("failed to parse schema: %w", gqlErr)
	}
	loadedSchema = schema
	return nil
}

func getSchema() *ast.Schema {
	return loadedSchema
}

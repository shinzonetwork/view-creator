package util

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func ValidateSDL(sdl string) error {
	directives := `
directive @materialized(if: Boolean) on OBJECT
`
	sdlWithDirectives := directives + "\n" + sdl

	schema, err := gqlparser.LoadSchema(&ast.Source{
		Name:  "sdl.graphql",
		Input: sdlWithDirectives,
	})
	if err != nil {
		return fmt.Errorf("invalid SDL: %w", err)
	}

	definedTypes := make(map[string]bool)
	for typeName := range schema.Types {
		definedTypes[typeName] = true
	}

	builtins := map[string]bool{
		"Int": true, "Float": true, "String": true, "Boolean": true, "ID": true,
	}

	for _, def := range schema.Types {
		if strings.HasPrefix(def.Name, "__") || builtins[def.Name] {
			continue
		}

		for _, field := range def.Fields {
			baseType := unwrapType(field.Type)
			if !builtins[baseType] && !definedTypes[baseType] {
				return fmt.Errorf("undefined type used in SDL: %s (in %s.%s)", baseType, def.Name, field.Name)
			}
		}
	}

	return nil
}

func ValidateSchemaBlock(schema string) error {
	src := &ast.Source{
		Name:  "new-type.graphql",
		Input: schema,
	}

	_, err := parser.ParseSchema(src)
	if err != nil {
		return err
	}
	return nil
}

func unwrapType(t *ast.Type) string {
	if t == nil {
		return ""
	}
	for t.Elem != nil {
		t = t.Elem
	}
	return t.NamedType
}

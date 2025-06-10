package sdl

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

func ValidateFieldTypes(doc *ast.SchemaDocument) error {
	definedTypes := map[string]bool{}
	for _, def := range doc.Definitions {
		definedTypes[def.Name] = true
	}

	for _, def := range doc.Definitions {
		for _, field := range def.Fields {
			baseType := unwrapType(field.Type)
			if !isBuiltinType(baseType) && !definedTypes[baseType] {
				return fmt.Errorf("field '%s' in '%s' references unknown type '%s'",
					field.Name, def.Name, baseType)
			}
		}
	}
	return nil
}

func unwrapType(t *ast.Type) string {
	for t.Elem != nil {
		t = t.Elem
	}
	return t.Name()
}

func isBuiltinType(name string) bool {
	switch name {
	case "Int", "String", "Float", "Boolean", "ID":
		return true
	default:
		return false
	}
}

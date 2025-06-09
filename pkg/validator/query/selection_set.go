package query

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateSelectionSet(query string) (string, []string, error) {
	schema := getSchema()
	if schema == nil {
		return "", nil, fmt.Errorf("schema not loaded")
	}

	re := regexp.MustCompile(`(?m)^\s*(\w+)\s*{\s*([\w\s]+)\s*}\s*$`)
	matches := re.FindStringSubmatch(query)
	if len(matches) != 3 {
		return "", nil, fmt.Errorf("invalid query format. Expected: Type { field1 field2 }")
	}

	entity := matches[1]
	fields := strings.Fields(matches[2])

	obj := schema.Types[entity]
	if obj == nil || obj.Kind != "OBJECT" {
		return "", nil, fmt.Errorf("type '%s' not found in schema", entity)
	}

	for _, field := range fields {
		if obj.Fields.ForName(field) == nil {
			return "", nil, fmt.Errorf("field '%s' not found in type '%s'", field, entity)
		}
	}

	return entity, fields, nil
}

package query

type QueryCheck struct {
	Entity string
	Fields []string
}

func Validate(query string) (*QueryCheck, error) {
	entity, fields, err := ValidateSelectionSet(query)
	if err != nil {
		return nil, err
	}

	// TODO: add future checks here: filters, nested fields, args, etc.

	return &QueryCheck{
		Entity: entity,
		Fields: fields,
	}, nil
}

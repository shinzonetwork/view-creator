package sdl

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func ParseSDL(sdl string) (*ast.SchemaDocument, error) {
	src := &ast.Source{
		Input: sdl,
		Name:  "inline-sdl",
	}

	doc, err := parser.ParseSchema(src)
	if err != nil {
		return nil, fmt.Errorf("invalid SDL syntax: %w", err)
	}

	return doc, nil
}

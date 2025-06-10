package sdl

import "fmt"

func Validate(sdl string) error {
	doc, err := ParseSDL(sdl)
	if err != nil {
		return err
	}

	if err := ValidateFieldTypes(doc); err != nil {
		return fmt.Errorf("field type validation failed: %w", err)
	}

	// Future validations can go here:

	return nil
}

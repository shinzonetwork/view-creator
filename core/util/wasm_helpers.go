package util

import (
	"bytes"
	"fmt"
	"io"
)

func IsValidWasm(file io.Reader) (io.Reader, error) {
	const wasmMagic = "\x00asm"

	header := make([]byte, 4)
	n, err := io.ReadFull(file, header)
	if err != nil || n != 4 {
		return nil, fmt.Errorf("unable to read wasm header: %w", err)
	}

	if string(header) != wasmMagic {
		return nil, fmt.Errorf("file is not a valid wasm")
	}

	return io.MultiReader(bytes.NewReader(header), file), nil
}

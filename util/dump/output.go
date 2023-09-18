package dump

import (
	"io"
	"os"
)

// output pre-defined
const (
	OutputStdout = "" // represents stdout by ""
)

// CreateOutput creates output writer if needed, return os.Stdout for 'stdout'.
func CreateOutput(output string) (io.Writer, io.Closer, error) {
	if output == OutputStdout {
		return os.Stdout, nil, nil // stdout doesn't want to be closed
	}

	w, err := os.Create(output)
	return w, w, err
}

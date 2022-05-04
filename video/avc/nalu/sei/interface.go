package sei

import "io"

type payloadParser interface {

	// input reader and expect size,
	// output parsed bytes or error
	parse(io.Reader, int) (uint64, error)
}

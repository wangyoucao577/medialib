package nalu

import "io"

// NALUParser defines NALU Parse interface.
type NALUParser interface {

	// Parse parses bytes to NALU data.
	// @param r io.Reader: Reader where to read data from
	// @param size int: how many bytes expect to read, 0 means no limit
	// @return read size or error
	Parse(r io.Reader, size int) (uint64, error)
}

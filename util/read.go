package util

import (
	"io"
)

const (
	// use '-' to represent stdin as input
	InputStdin = "-"
)

// ReadOrError reads a specifed amount of data, otherwise error.
func ReadOrError(r io.Reader, data []byte) error {

	l := len(data)
	readN := 0
	for {
		n, err := r.Read(data[readN:])
		if err != nil {
			return err
		}
		readN += n
		if readN == l {
			break
		}
	}

	return nil
}

// ReadByteOrError reads a byte in success, otherwise error.
func ReadByteOrError(r io.Reader) (byte, error) {
	nextByte := make([]byte, 1)
	if err := ReadOrError(r, nextByte); err != nil {
		return 0, err
	}
	return nextByte[0], nil
}

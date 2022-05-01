package util

import (
	"errors"
	"fmt"
	"io"
)

// ReadOrError reads a specifed amount of data, otherwise error.
func ReadOrError(r io.Reader, data []byte) error {

	l := len(data)
	n, err := r.Read(data)
	if err != nil {
		return err
	} else if n != l {
		s := fmt.Sprintf("expect to read %d bytes but got %d bytes", l, n)
		return errors.New(s)
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

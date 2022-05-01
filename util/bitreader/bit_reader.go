// Package bitreader deifnes reader per bit for easier bit level parsing.
package bitreader

import (
	"io"

	"github.com/wangyoucao577/medialib/util"
)

const (
	bitsPerByte = 8
)

// Reader implements bit level reader.
type Reader struct {
	cache      byte
	cachedBits uint // [0,8]

	r io.Reader
}

// New creates a new bit reader.
func New(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

// ReadBit reads one bit, return error if occurred.
func (r *Reader) ReadBit() (byte, error) {
	if r.cachedBits == 0 {
		if nextByte, err := util.ReadByteOrError(r.r); err != nil {
			return 0, err
		} else {
			r.cache = nextByte
			r.cachedBits = bitsPerByte
		}
	}

	newBit := r.cache >> (r.cachedBits - 1) & 0x1
	r.cachedBits--
	return newBit, nil
}

// ReadBits reads count specified bits, return error if occurred.
func (r *Reader) ReadBits(count uint) ([]byte, error) {
	bits := []byte{}
	bitsIndex := -1

	for i := 0; i < int(count); i++ {
		if i%bitsPerByte == 0 {
			bits = append(bits, 0)
			bitsIndex++
		}

		if nextBit, err := r.ReadBit(); err != nil {
			return nil, err
		} else {
			bits[bitsIndex] = (bits[bitsIndex] << 1) | nextBit
		}
	}
	return bits, nil
}

// ReadByte reads a byte.
func (r *Reader) ReadByte() (byte, error) {
	bits, err := r.ReadBits(bitsPerByte)
	if err != nil {
		return 0, err
	}
	return bits[0], nil
}

// ReadBytes reads specific count of bytes.
func (r *Reader) ReadBytes(count uint) ([]byte, error) {
	return r.ReadBits(bitsPerByte * count)
}

// CachedBitsCount returns cached bits count.
func (r *Reader) CachedBitsCount() int {
	return int(r.cachedBits)
}

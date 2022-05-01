// Package expgolombcoding implements Exponential-Golomb Coding,
// see more in https://en.wikipedia.org/wiki/Exponential-Golomb_coding or ISO/IEC-14496-10 7.2
package expgolombcoding

import (
	"encoding/json"
	"fmt"

	"github.com/wangyoucao577/medialib/util/bitreader"
)

// Unsigned represents unsigned Exponential-Golomb coding.
type Unsigned struct {

	// parsing
	leadingZeroBit int

	// store parsed
	value uint64
}

// MarshalJSON implements json.Marshaler.
func (u Unsigned) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Value())
}

// Value returns parsed value.
func (u *Unsigned) Value() uint64 {
	return u.value
}

// Parse parses unsigned value of Exponential-Golomb Coding.
// return the cost bits(NOT Byte) if succeed, otherwise error.
func (u *Unsigned) Parse(r *bitreader.Reader) (uint64, error) {
	if r == nil {
		return 0, fmt.Errorf("invalid bit reader")
	}

	var parsedBits uint64

	u.leadingZeroBit = -1
	for b := uint8(0); b == 0; u.leadingZeroBit++ {
		var err error
		if b, err = r.ReadBit(); err != nil {
			return uint64(parsedBits), err
		} else {
			parsedBits++
		}
	}

	for i := 0; i < u.leadingZeroBit; i++ {
		if b, err := r.ReadBit(); err != nil {
			return uint64(parsedBits), err
		} else {
			u.value = (u.value << 1) | uint64(b)
			parsedBits++
		}
	}

	u.value += (1 << u.leadingZeroBit) - 1
	return uint64(parsedBits), nil
}

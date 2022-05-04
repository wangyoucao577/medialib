package expgolombcoding

import (
	"encoding/json"
	"math"

	"github.com/wangyoucao577/medialib/util/bitreader"
)

// Signed contains signed Exponential-Golomb coded integer.
type Signed struct {
	unsigned Unsigned
}

// MarshalJSON implements json.Marshaler.
func (s *Signed) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value())
}

// Value returns parsed value.
func (s *Signed) Value() int64 {
	// ISO/IEC-14496-10 9.1.1
	v := math.Ceil(float64(s.unsigned.value) / 2)
	v *= math.Pow(-1, float64(s.unsigned.value+1))
	return int64(v)
}

// Parse parses unsigned value of Exponential-Golomb Coding.
// return the cost bits(NOT Byte) if succeed, otherwise error.
func (s *Signed) Parse(r *bitreader.Reader) (uint64, error) {
	return s.unsigned.Parse(r)
}

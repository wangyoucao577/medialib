package expgolombcoding

import (
	"encoding/json"

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
	//TODO: correct signed value
	return int64(s.unsigned.value)
}

// Parse parses unsigned value of Exponential-Golomb Coding.
// return the cost bits(NOT Byte) if succeed, otherwise error.
func (s *Signed) Parse(r *bitreader.Reader) (uint64, error) {
	return s.unsigned.Parse(r)
}

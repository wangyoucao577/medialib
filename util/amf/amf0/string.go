package amf0

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// StringPayload represents AMF0 String type payload.
type StringPayload struct {
	Length uint16 `json:"length"`
	Str    string `json:"string"`
}

// Decode decodes AMF0 string type payload.
func (s *StringPayload) Decode(r io.Reader) (int, error) {

	var parsedBytes int

	data := make([]byte, 2)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		s.Length = binary.BigEndian.Uint16(data)
		parsedBytes += 2
	}

	data = make([]byte, s.Length)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		s.Str = string(data)
		parsedBytes += int(s.Length)
	}

	return parsedBytes, nil
}

// Encode encodes StringPayload to AMF0 byte stream.
func (s StringPayload) Encode() ([]byte, error) {
	if int(s.Length) != len(s.Str) {
		return nil, fmt.Errorf("string length %d is not equal to len(str) %d", s.Length, len(s.Str))
	}

	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, s.Length)

	data = append(data, []byte(s.Str)...)

	return data, nil
}

package amf0

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// StrictArrayPayload represents payload of AMF0 strict array.
type StrictArrayPayload struct {
	Count     uint32      `json:"count"`
	ValueType []ValueType `json:"value_type"`
}

// Decode implements decoder.
func (s *StrictArrayPayload) Decode(r io.Reader) (int, error) {
	var parsedBytes int

	data := make([]byte, 4)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		s.Count = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	for i := 0; i < int(s.Count); i++ {
		v := ValueType{}
		if bytes, err := v.Decode(r); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
		s.ValueType = append(s.ValueType, v)

		if v.TypeMarker == TypeMarkerObjectEnd { // object-end-marker if exist
			break
		}
	}
	return parsedBytes, nil
}

// Encode implements encoder interface.
func (s StrictArrayPayload) Encode() ([]byte, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, s.Count)

	for _, v := range s.ValueType {
		if d, err := v.Encode(); err != nil {
			return data, err
		} else {
			data = append(data, d...)
		}
	}

	return data, nil
}

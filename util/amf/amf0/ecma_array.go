package amf0

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// ECMAArrayPayload represents payload of AMF0 strict array.
type ECMAArrayPayload struct {
	Count          uint32           `json:"count"`
	ObjectProperty []ObjectProperty `json:"object-property"`
}

// Decode implements decoder.
func (e *ECMAArrayPayload) Decode(r io.Reader) (int, error) {
	var parsedBytes int

	data := make([]byte, 4)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		e.Count = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	for {
		op := ObjectProperty{}
		if bytes, err := op.Decode(r); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
		e.ObjectProperty = append(e.ObjectProperty, op)

		if op.ValueType.TypeMarker == TypeMarkerObjectEnd { // object-end-marker if exist
			break
		}
	}
	return parsedBytes, nil
}

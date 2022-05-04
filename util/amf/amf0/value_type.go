package amf0

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/wangyoucao577/medialib/util"
)

// ValueType represents an AMF0 value-type.
type ValueType struct {
	TypeMarker uint8 `json:"type_marker"`

	Value interface{} `json:"value,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (v ValueType) MarshalJSON() ([]byte, error) {
	var dj = struct {
		TypeMarker            uint8  `json:"type_marker"`
		TypeMarkerDescription string `json:"type_marker_description"`

		Value interface{} `json:"value,omitempty"`
	}{
		TypeMarker:            v.TypeMarker,
		TypeMarkerDescription: TypeMarkerDescription(int(v.TypeMarker)),

		Value: v.Value,
	}
	return json.Marshal(dj)
}

// Decode decodes bytes to AMF0 value type.
func (v *ValueType) Decode(r io.Reader) (int, error) {
	var parsedBytes int

	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		v.TypeMarker = data[0]
		parsedBytes += 1
	}

	var dec decoder
	switch v.TypeMarker {
	case TypeMarkerNumber:
		data := make([]byte, 8)
		if err := util.ReadOrError(r, data); err != nil {
			return parsedBytes, err
		} else {
			u := binary.BigEndian.Uint64(data)
			v.Value = math.Float64frombits(u)
			parsedBytes += 8
		}
	case TypeMarkerBoolean:
		if nextByte, err := util.ReadByteOrError(r); err != nil {
			return parsedBytes, err
		} else {
			v.Value = false
			if nextByte != 0 {
				v.Value = true
			}
			parsedBytes += 1
		}
	case TypeMarkerString:
		dec = &StringPayload{}
	case TypeMarkerObject:
		dec = &ObjectPayload{}
	case TypeMarkerNull: // nothing to do
	case TypeMarkerUndefined: // nothing to do
	case TypeMarkerReference:
		data := make([]byte, 2)
		if err := util.ReadOrError(r, data); err != nil {
			return parsedBytes, err
		} else {
			v.Value = binary.BigEndian.Uint16(data)
			parsedBytes += 2
		}
	case TypeMarkerECMAArray:
		dec = &ECMAArrayPayload{}
	case TypeMarkerObjectEnd: // nothing to do
	case TypeMarkerStrictArray:
		dec = &StrictArrayPayload{}
	case TypeMarkerDate:
		dec = &Date{}
	default:
		return parsedBytes, fmt.Errorf("AMF0 type %d(%s) unsupported", v.TypeMarker, TypeMarkerDescription(int(v.TypeMarker)))
	}

	if dec != nil {
		if bytes, err := dec.Decode(r); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
		v.Value = dec
	}

	return parsedBytes, nil
}

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

// Encode encodes AMF0 value type to bytes.
func (v ValueType) Encode() ([]byte, error) {
	data := []byte{v.TypeMarker}

	var enc encoder
	switch v.TypeMarker {
	case TypeMarkerNumber:
		if f, err := v.AsNumber(); err != nil {
			return data, err
		} else {
			numData := make([]byte, 8)
			binary.BigEndian.PutUint64(numData, math.Float64bits(f))
			data = append(data, numData...)
		}
	case TypeMarkerBoolean:
		if b, err := v.AsBoolean(); err != nil {
			return data, err
		} else {
			var bData byte
			if b {
				bData = 0x1
			}
			data = append(data, bData)
		}
	case TypeMarkerReference:
		if u, ok := v.Value.(uint16); !ok {
			return data, fmt.Errorf("value type %v as uint16 failed", v)
		} else {
			uData := make([]byte, 2)
			binary.BigEndian.PutUint16(uData, u)
			data = append(data, uData...)
		}

	case TypeMarkerNull: // nothing to do
	case TypeMarkerUndefined: // nothing to do
	case TypeMarkerObjectEnd: // nothing to do

	case TypeMarkerString:
		fallthrough
	case TypeMarkerObject:
		fallthrough
	case TypeMarkerECMAArray:
		fallthrough
	case TypeMarkerStrictArray:
		fallthrough
	case TypeMarkerDate:
		if converted, ok := v.Value.(encoder); !ok {
			return data, fmt.Errorf("value type %v as encoder interface failed", v)
		} else {
			enc = converted
		}
	default:
		return data, fmt.Errorf("AMF0 type %d(%s) unsupported", v.TypeMarker, TypeMarkerDescription(int(v.TypeMarker)))
	}

	if enc != nil {
		if encodedData, err := enc.Encode(); err != nil {
			return data, err
		} else {
			data = append(data, encodedData...)
		}
	}

	return data, nil
}

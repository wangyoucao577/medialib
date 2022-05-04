package amf0

import (
	"bytes"
	"reflect"
	"testing"
)

const floatEpsilon = 0.000001

func TestValueTypeDecode(t *testing.T) {
	var cases = []struct {
		in          []byte
		vt          ValueType
		parsedBytes int
	}{
		{
			[]byte{0x00, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			ValueType{
				TypeMarker: TypeMarkerNumber,
				Value:      float64(0x01),
			},
			9,
		},
		{
			[]byte{0x01, 0x01},
			ValueType{
				TypeMarker: TypeMarkerBoolean,
				Value:      true,
			},
			2,
		},
		{
			[]byte{0x02, 0x00, 0x07, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63, 0x74},
			ValueType{
				TypeMarker: TypeMarkerString,
				Value: &StringPayload{
					Length: 7,
					Str:    "connect",
				},
			},
			10,
		},
	}

	for _, c := range cases {
		vt := ValueType{}
		if _, err := vt.Decode(bytes.NewReader(c.in)); err != nil {
			t.Error(err)
		}
		switch c.vt.TypeMarker {
		case TypeMarkerNumber:
			cNum, err := c.vt.AsNumber()
			if err != nil {
				t.Error(err)
			}
			num, err := vt.AsNumber()
			if err != nil {
				t.Error(err)
			}
			if cNum-num > floatEpsilon {
				t.Errorf("decode %v expect %v but got %v", c.in, c.vt, vt)
			}
		default:
			if !reflect.DeepEqual(vt, c.vt) {
				t.Errorf("decode %v expect %v but got %v", c.in, c.vt, vt)
			}
		}

	}
}

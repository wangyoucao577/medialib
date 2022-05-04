package script

import (
	"encoding/json"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// DataValue represents ScriptDataValue.
type DataValue struct {

	// 	Type of the ScriptDataValue. The following types are defined:
	// 0 = Number
	// 1 = Boolean
	// 2 = String
	// 3 = Object
	// 4 = MovieClip (reserved, not supported)
	// 5 = Null
	// 6 = Undefined
	// 7 = Reference
	// 8 = ECMA array
	// 9 = Object end marker
	// 10 = Strict array
	// 11 = Date
	// 12 = Long string
	Type uint8 `json:"Type"`

	// IF Type == 0 DOUBLE
	// IF Type == 1 UI8
	// IF Type == 2 SCRIPTDATASTRING
	// IF Type == 3 SCRIPTDATAOBJECT
	// IF Type == 7 UI16
	// IF Type == 8 SCRIPTDATAECMAARRAY
	// IF Type == 10 SCRIPTDATASTRICTARRAY
	// IF Type == 11 SCRIPTDATADATE
	// IF Type == 12 SCRIPTDATALONGSTRING
	// Only one of below will not empty.
	DataString *DataString `json:"string,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (d DataValue) MarshalJSON() ([]byte, error) {
	var dj = struct {
		Type            uint8  `json:"Type"`
		TypeDescription string `json:"TypeDescription"`

		DataString *DataString `json:"string,omitempty"`
	}{
		Type:            d.Type,
		TypeDescription: DataTypeDescription(int(d.Type)),

		DataString: d.DataString,
	}
	return json.Marshal(dj)
}

func (d *DataValue) parse(r io.Reader) (uint64, error) {
	var parsedBytes uint64

	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.Type = data[0]
		parsedBytes += 1
	}

	if d.Type == 2 {
		d.DataString = &DataString{}
		if bytes, err := d.DataString.parse(r); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
	}

	return parsedBytes, nil
}

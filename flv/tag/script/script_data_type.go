package script

// data types
const (
	DataTypeDouble = 0
	DataTypeUint8  = 1
	DataTypeString = 2
	DataTypeObject = 3

	DataTypeUint16    = 7
	DataTypeECMAArray = 8

	DataTypeStrictArray = 10
	DataTypeDate        = 11
	DataTypeLongString  = 12
)

var dataTypeDescriptions = map[int]string{
	DataTypeDouble: "double",
	DataTypeUint8:  "uint8",
	DataTypeString: "SCRIPTDATASTRING",
	DataTypeObject: "SCRIPTDATAOBJECT",

	DataTypeUint16:    "uint16",
	DataTypeECMAArray: "SCRIPTDATAECMAARRAY",

	DataTypeStrictArray: "SCRIPTDATASTRICTARRAY",
	DataTypeDate:        "SCRIPTDATADATE",
	DataTypeLongString:  "SCRIPTDATALONGSTRING",
}

// DataTypeDescription returns description of data type.
func DataTypeDescription(t int) string {
	d, ok := dataTypeDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

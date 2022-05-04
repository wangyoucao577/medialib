// Package amf0 implements AMF0 specification.
package amf0

// type markers
const (
	TypeMarkerNumber      = 0x00
	TypeMarkerBoolean     = 0x01
	TypeMarkerString      = 0x02
	TypeMarkerObject      = 0x03
	TypeMarkerMoiveClip   = 0x04 // reserved, not supported
	TypeMarkerNull        = 0x05
	TypeMarkerUndefined   = 0x06
	TypeMarkerReference   = 0x07
	TypeMarkerECMAArray   = 0x08
	TypeMarkerObjectEnd   = 0x09
	TypeMarkerStrictArray = 0x0A
	TypeMarkerDate        = 0x0B
	TypeMarkerLongString  = 0x0C
	TypeMarkerUnsupported = 0x0D
	TypeMarkerRecordSet   = 0x0E // reserved, not supported
	TypeMarkerXMLDocument = 0x0F
	TypeMarkerTypedObject = 0x10
)

var typeMarkerDescriptions = map[int]string{
	TypeMarkerNumber:      "number",
	TypeMarkerBoolean:     "boolean",
	TypeMarkerString:      "string",
	TypeMarkerObject:      "object",
	TypeMarkerMoiveClip:   "movieclip, reserved, not supported",
	TypeMarkerNull:        "null",
	TypeMarkerUndefined:   "undefined",
	TypeMarkerReference:   "reference",
	TypeMarkerECMAArray:   "ecma-array",
	TypeMarkerObjectEnd:   "object-end",
	TypeMarkerStrictArray: "strict-array",
	TypeMarkerDate:        "date",
	TypeMarkerLongString:  "long-string",
	TypeMarkerUnsupported: "unsupported",
	TypeMarkerRecordSet:   "recordset, reserved, not supported",
	TypeMarkerXMLDocument: "xml-document",
	TypeMarkerTypedObject: "typed-object",
}

// TypeMarkerDescription returns description of type marker.
func TypeMarkerDescription(t int) string {
	d, ok := typeMarkerDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

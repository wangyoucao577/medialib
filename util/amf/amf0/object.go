package amf0

import "io"

// ObjectProperty represents
type ObjectProperty struct {
	String    StringPayload `json:"name"`
	ValueType ValueType     `json:"value_type"`
}

// ObjectPayload represents AMF0 object payload.
type ObjectPayload struct {
	ObjectProperty []ObjectProperty `json:"object-property"`
}

// Decode implements decoder.
func (o *ObjectProperty) Decode(r io.Reader) (int, error) {
	var parsedBytes int

	if bytes, err := o.String.Decode(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	if bytes, err := o.ValueType.Decode(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	return parsedBytes, nil
}

// Decode implements decoder.
func (o *ObjectPayload) Decode(r io.Reader) (int, error) {
	var parsedBytes int

	for {
		op := ObjectProperty{}
		if bytes, err := op.Decode(r); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
		o.ObjectProperty = append(o.ObjectProperty, op)

		if op.ValueType.TypeMarker == TypeMarkerObjectEnd { // object-end-marker
			break
		}
	}
	return parsedBytes, nil
}

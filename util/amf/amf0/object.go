package amf0

import "io"

// ObjectProperty represents AMF0 object property.
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

// Encode implements encoder interface.
func (o ObjectProperty) Encode() ([]byte, error) {
	var data []byte

	if s, err := o.String.Encode(); err != nil {
		return data, err
	} else {
		data = append(data, s...)
	}

	if v, err := o.ValueType.Encode(); err != nil {
		return data, err
	} else {
		data = append(data, v...)
	}

	return data, nil
}

// Encode implements encoder interface.
func (o ObjectPayload) Encode() ([]byte, error) {
	var data []byte

	for _, op := range o.ObjectProperty {
		if d, err := op.Encode(); err != nil {
			return data, err
		} else {
			data = append(data, d...)
		}
	}

	return data, nil
}

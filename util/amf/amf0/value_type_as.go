package amf0

import "fmt"

// AsNumber returns number value if type matched.
func (v ValueType) AsNumber() (float64, error) {
	if v.Value == nil {
		return 0, fmt.Errorf("empty value")
	}
	if v.TypeMarker != TypeMarkerNumber {
		return 0, fmt.Errorf("type 0x%x(%s) unmatched with expect %d(%s)",
			v.TypeMarker, TypeMarkerDescription(int(v.TypeMarker)), TypeMarkerNumber, TypeMarkerDescription(TypeMarkerNumber))
	}

	f, ok := v.Value.(float64)
	if !ok {
		return 0, fmt.Errorf("value as float64 failed")
	}
	return f, nil
}

// AsBoolean returns boolean value if type matched.
func (v ValueType) AsBoolean() (bool, error) {
	if v.Value == nil {
		return false, fmt.Errorf("empty value")
	}
	if v.TypeMarker != TypeMarkerBoolean {
		return false, fmt.Errorf("type 0x%x(%s) unmatched with expect %d(%s)",
			v.TypeMarker, TypeMarkerDescription(int(v.TypeMarker)), TypeMarkerBoolean, TypeMarkerDescription(TypeMarkerBoolean))
	}

	b, ok := v.Value.(bool)
	if !ok {
		return false, fmt.Errorf("value as bool failed")
	}
	return b, nil
}

// AsString returns string value if type matched.
func (v ValueType) AsString() (string, error) {
	if v.Value == nil {
		return "", fmt.Errorf("empty value")
	}
	if v.TypeMarker != TypeMarkerString {
		return "", fmt.Errorf("type 0x%x(%s) unmatched with expect %d(%s)",
			v.TypeMarker, TypeMarkerDescription(int(v.TypeMarker)), TypeMarkerString, TypeMarkerDescription(TypeMarkerString))
	}

	s, ok := v.Value.(StringPayload)
	if !ok {
		return "", fmt.Errorf("value as string failed")
	}
	return s.Str, nil
}

// AsReference returns reference value if type matched.
func (v ValueType) AsReference() (uint16, error) {
	if v.Value == nil {
		return 0, fmt.Errorf("empty value")
	}
	if v.TypeMarker != TypeMarkerReference {
		return 0, fmt.Errorf("type 0x%x(%s) unmatched with expect %d(%s)",
			v.TypeMarker, TypeMarkerDescription(int(v.TypeMarker)), TypeMarkerReference, TypeMarkerDescription(TypeMarkerReference))
	}

	r, ok := v.Value.(uint16)
	if !ok {
		return 0, fmt.Errorf("value as string failed")
	}
	return r, nil
}

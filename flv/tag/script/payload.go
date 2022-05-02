package script

import "io"

// TagBody represents script data tag payload.
type TagBody struct {
	Name  DataValue `json:"name"`
	Value DataValue `json:"value"`
}

func (t *TagBody) parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	if bytes, err := t.Name.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	if bytes, err := t.Value.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	return parsedBytes, nil
}

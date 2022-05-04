package script

import (
	"io"

	"github.com/wangyoucao577/medialib/util/amf/amf0"
)

// TagBody represents script data tag payload.
type TagBody struct {
	Name  amf0.ValueType `json:"name"`
	Value amf0.ValueType `json:"value"`
}

func (t *TagBody) parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	if bytes, err := t.Name.Decode(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += uint64(bytes)
	}

	if bytes, err := t.Value.Decode(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += uint64(bytes)
	}

	return parsedBytes, nil
}

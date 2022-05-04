package amf0

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// Date reprsents AMF0 date type.
type Date struct {
	Timestamp uint64 `json:"timestamp"`
	TimeZone  int16  `json:"time_zone"`
}

// Decode implements decoder.
func (d *Date) Decode(r io.Reader) (int, error) {
	var parsedBytes int

	data := make([]byte, 8)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.Timestamp = binary.BigEndian.Uint64(data)
		parsedBytes += 8
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return parsedBytes, err
	} else {
		d.TimeZone = int16(binary.BigEndian.Uint16(data[:2]))
		parsedBytes += 2
	}

	return parsedBytes, nil
}

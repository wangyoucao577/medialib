package script

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

type DataString struct {
	Length uint16 `json:"length"`
	Data   string `json:"data"` // String data, up to 65535 bytes, with no terminating NUL
}

func (d *DataString) parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	data := make([]byte, 2)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.Length = binary.BigEndian.Uint16(data)
		parsedBytes += 2
	}

	data = make([]byte, d.Length)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.Data = string(data)
		parsedBytes += uint64(d.Length)
	}

	return parsedBytes, nil
}

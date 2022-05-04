package esds

import (
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// DecoderSpecificInfo represents DecoderSpecificInfo.
type DecoderSpecificInfo struct {
	Descriptor Descriptor `json:"descriptor"`

	Data []byte `json:"data"`
}

func (d *DecoderSpecificInfo) parse(r io.Reader) (uint64, error) {
	var parsedBytes uint64

	// parse descriptor header
	if bytes, err := d.Descriptor.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	if d.Descriptor.Size > 0 {
		d.Data = make([]byte, d.Descriptor.Size)
		if err := util.ReadOrError(r, d.Data); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += uint64(d.Descriptor.Size)
		}
	}

	return parsedBytes, nil
}

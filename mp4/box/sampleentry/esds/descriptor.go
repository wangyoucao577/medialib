package esds

import (
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// Descriptor represents base descriptor.
type Descriptor struct {
	Tag  uint8  `json:"tag"`
	Size uint32 `json:"size"`
}

func (d *Descriptor) parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64
	data := make([]byte, 4)

	// first bytes is tag
	if err := util.ReadOrError(r, data[:1]); err != nil {
		return parsedBytes, err
	} else {
		d.Tag = data[0]
		parsedBytes += 1
	}

	if bytes, err := d.readSize(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	return parsedBytes, nil
}

func (d *Descriptor) readSize(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	data := make([]byte, 1)

	// The elementary stream size is specific by up to 4 bytes.
	// The MSB of a byte indicates if there are more bytes for the size.
	// Reference to below chromium implementation for more details
	// 	https://chromium.googlesource.com/chromium/src/media/+/16ba1c56b860d53d7354c0ec9538650cf1f20e2d/mp4/es_descriptor.cc
	for i := 0; i < 4; i++ {
		var msb uint8

		if err := util.ReadOrError(r, data); err != nil {
			return parsedBytes, err
		} else {
			msb = (data[0] >> 7) & 0x1
			s := uint32(data[0] & 0x7F)
			d.Size = d.Size<<7 + s
			parsedBytes += 1
		}

		if msb == 0 {
			break
		}
	}

	return parsedBytes, nil
}

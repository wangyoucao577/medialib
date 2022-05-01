// Package cprt represents Copyright Box.
package cprt

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a cprt box.
type Box struct {
	box.FullHeader `json:"full_header"`

	Pad      uint8    `json:"pad"`      // 1 bit
	Language [3]uint8 `json:"language"` // 5 bytes per uint
	Notice   string   `json:"notice"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		FullHeader: box.FullHeader{
			Header: h,
		},
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	// parse full header additional information first
	if err := b.FullHeader.ParseVersionFlag(r); err != nil {
		return err
	}

	// start to parse payload
	var parsedBytes uint64

	data := make([]byte, 2)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.Pad = (uint8(data[0]) >> 7) & 0x1 // 1 bit

		// 5 bits per Language
		b.Language[0] = (uint8(data[0]) >> 2) & 0x1F
		b.Language[1] = ((uint8(data[0]) & 0x3) << 3) | ((uint8(data[1]) >> 5) & 0x7)
		b.Language[2] = uint8(data[1]) & 0x1F

		parsedBytes += 2
	}

	data = make([]byte, b.PayloadSize()-parsedBytes)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.Notice = string(data)
	}

	return nil
}

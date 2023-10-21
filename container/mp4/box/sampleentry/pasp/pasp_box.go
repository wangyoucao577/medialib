// Package pasp represents PixelAspectRatio Box.
package pasp

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a PixelAspectRatio box.
type Box struct {
	box.Header `json:"header"`

	HSpacing uint32 `json:"hSpacing"`
	VSpacing uint32 `json:"vSpacing"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Header.Type, err)
		return nil
	}

	// start to parse payload
	var parsedBytes uint64

	data := make([]byte, 8)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.HSpacing = binary.BigEndian.Uint32(data[:4])
		b.VSpacing = binary.BigEndian.Uint32(data[4:])
		parsedBytes += 8
	}

	return nil
}

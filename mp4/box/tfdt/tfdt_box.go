// Package tfdt represents Track Fragment Base Media Decode Time Box.
package tfdt

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a tfdt box.
type Box struct {
	box.FullHeader `json:"full_header"`

	BaseMediaDecodeTime uint64 `json:"base_media_decode_time"`
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

	var parsedBytes uint64

	data := make([]byte, 8)
	if b.Version != 1 {
		data = data[:4]
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		if b.Version == 1 {
			b.BaseMediaDecodeTime = binary.BigEndian.Uint64(data)
			parsedBytes += 8
		} else {
			b.BaseMediaDecodeTime = uint64(binary.BigEndian.Uint32(data))
			parsedBytes += 4
		}
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

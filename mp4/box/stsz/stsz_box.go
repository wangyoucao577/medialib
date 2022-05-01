// Package stsz represents Sample Size Box.
package stsz

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a stsz box.
type Box struct {
	box.FullHeader `json:"full_header"`

	SampleSize  uint32   `json:"sample_size"`
	SampleCount uint32   `json:"sample_count"`
	EntrySizes  []uint32 `json:"entry_size,omitempty"`
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

	data := make([]byte, 4)

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.SampleSize = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.SampleCount = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if b.SampleSize == 0 {
		for i := 0; i < int(b.SampleCount); i++ {
			if err := util.ReadOrError(r, data); err != nil {
				return err
			} else {
				size := binary.BigEndian.Uint32(data)
				b.EntrySizes = append(b.EntrySizes, size)
				parsedBytes += 4
			}
		}
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

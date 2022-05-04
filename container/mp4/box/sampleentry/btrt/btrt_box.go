// Package btrt represents MPEG4 Bit Rate Box.
package btrt

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a btrt box.
type Box struct {
	box.Header `json:"header"`

	BufferSizeDB uint32 `json:"buffer_size_db"` // gives the size of the decoding buffer for the elementary stream in bytes.
	MaxBitrate   uint32 `json:"max_bitrate"`    // gives the maximum rate in bits/second over any window of one second.
	AvgBitrate   uint32 `json:"avg_bitrate"`    // gives the average rate in bits/second over the entire presentation.
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
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	var parsedBytes uint64

	data := make([]byte, 4)

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.BufferSizeDB = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.MaxBitrate = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.AvgBitrate = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

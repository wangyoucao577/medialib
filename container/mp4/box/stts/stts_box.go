// Package stts represents Decoding Time to Sample Box.
package stts

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a stts box.
type Box struct {
	box.FullHeader `json:"full_header"`

	EntryCount   uint32   `json:"entry_count"`
	SampleCounts []uint32 `json:"sample_count,omitempty"`
	SampleDeltas []uint32 `json:"sample_delta,omitempty"`
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
		b.EntryCount = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	for i := 0; i < int(b.EntryCount); i++ {

		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			count := binary.BigEndian.Uint32(data)
			b.SampleCounts = append(b.SampleCounts, count)
			parsedBytes += 4
		}

		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			delta := binary.BigEndian.Uint32(data)
			b.SampleDeltas = append(b.SampleDeltas, delta)
			parsedBytes += 4
		}
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

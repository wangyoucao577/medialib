// Package trun represents Track Fragment Run Box.
package trun

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a trun box.
type Box struct {
	box.FullHeader `json:"full_header"`

	SampleCount                 uint32   `json:"sample_count"`
	DataOffset                  int32    `json:"data_offset,omitempty"`
	FirstSampleFlags            uint32   `json:"first_sample_flags,omitempty"`
	SampleDuration              []uint32 `json:"sample_duration,omitempty"`
	SampleSize                  []uint32 `json:"sample_size,omitempty"`
	SampleFlags                 []uint32 `json:"sample_flags,omitempty"`
	SampleCompositionTimeOffset []int64  `json:"sample_composition_time_offset,omitempty"` // int32 or uint32 in file
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
	payloadSize := b.PayloadSize()
	data := make([]byte, 4)

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.SampleCount = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if (b.Flags & 0x1) > 0 {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.DataOffset = int32(binary.BigEndian.Uint32(data))
			parsedBytes += 4
		}
	}

	if (b.Flags & 0x4) > 0 {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.FirstSampleFlags = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	for i := 0; i < int(b.SampleCount); i++ {
		if (b.Flags & 0x100) > 0 {
			if err := util.ReadOrError(r, data); err != nil {
				return err
			} else {
				sampleDuration := binary.BigEndian.Uint32(data)
				b.SampleDuration = append(b.SampleDuration, sampleDuration)
				parsedBytes += 4
			}
		}

		if (b.Flags & 0x200) > 0 {
			if err := util.ReadOrError(r, data); err != nil {
				return err
			} else {
				sampleSize := binary.BigEndian.Uint32(data)
				b.SampleSize = append(b.SampleSize, sampleSize)
				parsedBytes += 4
			}
		}

		if (b.Flags & 0x400) > 0 {
			if err := util.ReadOrError(r, data); err != nil {
				return err
			} else {
				sampleFlags := binary.BigEndian.Uint32(data)
				b.SampleFlags = append(b.SampleFlags, sampleFlags)
				parsedBytes += 4
			}
		}

		if (b.Flags & 0x800) > 0 {
			if err := util.ReadOrError(r, data); err != nil {
				return err
			} else {
				// don't need to check b.Version since we can use `int64` to compatible for both
				sampleCompositionTimeOffset := int64(binary.BigEndian.Uint32(data))
				b.SampleCompositionTimeOffset = append(b.SampleCompositionTimeOffset, sampleCompositionTimeOffset)
				parsedBytes += 4
			}
		}
	}

	if parsedBytes != payloadSize {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, payloadSize)
	}

	return nil
}

// Package tfhd represents Track Fragment Header Box.
package tfhd

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a tfhd box.
type Box struct {
	box.FullHeader `json:"full_header"`

	TrackID uint32 `json:"track_id"`

	// optional fields
	BaseDataOffset         uint64 `json:"base_data_offset,omitempty"`
	SampleDescriptionIndex uint32 `json:"sample_description_index,omitempty"`
	DefaultSampleDuration  uint32 `json:"default_sample_duration,omitempty"`
	DefaultSampleSize      uint32 `json:"default_sample_size,omitempty"`
	DefaultSampleFlags     uint32 `json:"default_sample_flags,omitempty"`
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
		b.TrackID = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if (b.Flags & 0x1) > 0 {
		offsetData := make([]byte, 8)
		if err := util.ReadOrError(r, offsetData); err != nil {
			return err
		} else {
			b.BaseDataOffset = binary.BigEndian.Uint64(offsetData)
			parsedBytes += 8
		}
	}

	if (b.Flags & 0x2) > 0 {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.SampleDescriptionIndex = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	if (b.Flags & 0x8) > 0 {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.DefaultSampleDuration = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	if (b.Flags & 0x10) > 0 {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.DefaultSampleSize = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	if (b.Flags & 0x20) > 0 {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.DefaultSampleFlags = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	if parsedBytes != payloadSize {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, payloadSize)
	}

	return nil
}

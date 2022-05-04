// Package trex represents Track Extends Box.
package trex

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a trex box.
type Box struct {
	box.FullHeader `json:"full_header"`

	TrackID                       uint32 `json:"track_id"`
	DefaultSampleDescriptionIndex uint32 `json:"default_sample_description_index"`
	DefaultSampleDuration         uint32 `json:"default_sample_duration"`
	DefaultSampleSize             uint32 `json:"default_sample_size"`
	DefaultSampleFlags            uint32 `json:"default_sample_flags"`
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

	arr := []*uint32{&b.TrackID, &b.DefaultSampleDescriptionIndex, &b.DefaultSampleDuration, &b.DefaultSampleSize, &b.DefaultSampleFlags}
	for i := 0; i < len(arr); i++ {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			*arr[i] = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

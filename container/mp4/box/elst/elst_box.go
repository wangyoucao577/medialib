// Package elst represents Edit List Box.
package elst

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a trak box.
type Box struct {
	box.FullHeader `json:"full_header"`

	EntryCount uint32 `json:"entry_count"`

	SegmentDuration   []uint64 `json:"segment_duration"`
	MediaTime         []int64  `json:"media_time"`
	MediaRateInteger  []int16  `json:"media_rate_integer"`
	MediaRateFraction []int16  `json:"media_rate_fraction"`
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

		if b.FullHeader.Version == 1 {
			data = make([]byte, 16)
		} else {
			data = make([]byte, 8)
		}

		if err := util.ReadOrError(r, data); err != nil {
			return err
		}

		if b.FullHeader.Version == 1 {
			num := binary.BigEndian.Uint64(data[:8])
			b.SegmentDuration = append(b.SegmentDuration, num)

			num = binary.BigEndian.Uint64(data[8:])
			b.MediaTime = append(b.MediaTime, int64(num))
			parsedBytes += 16
		} else {
			num := binary.BigEndian.Uint32(data[:4])
			b.SegmentDuration = append(b.SegmentDuration, uint64(num))

			num = binary.BigEndian.Uint32(data[4:])
			b.MediaTime = append(b.MediaTime, int64(num))
			parsedBytes += 8
		}

		if err := util.ReadOrError(r, data[:4]); err != nil {
			return err
		}

		num := binary.BigEndian.Uint16(data[:2])
		b.MediaRateInteger = append(b.MediaRateInteger, int16(num))

		num = binary.BigEndian.Uint16(data[2:])
		b.MediaRateFraction = append(b.MediaRateFraction, int16(num))
		parsedBytes += 4
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

// Package mdhd represents Media Header Box.
package mdhd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/util/time1904"
)

const (
	languageCodeOffset = 0x60 // ISO/IEC 14496-12:2015 8.4.2
)

// Box represents a mdhd box.
type Box struct {
	box.FullHeader `json:"full_header"`

	CreationTime     uint64 `json:"creation_time,string"`
	ModificationTime uint64 `json:"modification_time,string"`
	Timescale        uint32 `json:"timescale"`
	Duration         uint64 `json:"duration"`

	Pad        uint8   `json:"pad"`      // 1 bit
	Language   [3]byte `json:"language"` // 5 bytes per uint
	Predefined uint16  `json:"pre_defined"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		FullHeader: box.FullHeader{
			Header: h,
		},
	}
}

// MarshalJSON implements json.Marshaler interface.
func (b Box) MarshalJSON() ([]byte, error) {
	jsonBox := struct {
		box.FullHeader `json:"full_header"`

		CreationTime         time.Time `json:"creation_time"`
		ModificationTime     time.Time `json:"modification_time"`
		Timescale            uint32    `json:"timescale"`
		Duration             uint64    `json:"duration"`
		DurationMilliSeconds float64   `json:"duration_ms"`

		Pad        uint8  `json:"pad"`      // 1 bit
		Language   string `json:"language"` // 5 bytes per uint
		Predefined uint16 `json:"pre_defined"`
	}{
		FullHeader: b.FullHeader,

		CreationTime:         time1904.Unix(int64(b.CreationTime), 0).UTC(),
		ModificationTime:     time1904.Unix(int64(b.ModificationTime), 0).UTC(),
		Timescale:            b.Timescale,
		Duration:             b.Duration,
		DurationMilliSeconds: float64(b.Duration) * 1000 / float64(b.Timescale),

		Pad:        b.Pad,
		Language:   string(b.Language[:]),
		Predefined: b.Predefined,
	}

	return json.Marshal(jsonBox)
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

	timeDataSize := 4 // if Version == 0
	if b.FullHeader.Version == 1 {
		timeDataSize = 8
	}
	data := make([]byte, timeDataSize)

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		if timeDataSize == 8 {
			b.CreationTime = binary.BigEndian.Uint64(data)
		} else {
			b.CreationTime = uint64(binary.BigEndian.Uint32(data))
		}
		parsedBytes += uint64(timeDataSize)
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		if timeDataSize == 8 {
			b.ModificationTime = binary.BigEndian.Uint64(data)
		} else {
			b.ModificationTime = uint64(binary.BigEndian.Uint32(data))
		}
		parsedBytes += uint64(timeDataSize)
	}

	if err := util.ReadOrError(r, data[:4]); err != nil {
		return err
	} else {
		b.Timescale = binary.BigEndian.Uint32(data[:4])
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		if timeDataSize == 8 {
			b.Duration = binary.BigEndian.Uint64(data)
		} else {
			b.Duration = uint64(binary.BigEndian.Uint32(data))
		}
		parsedBytes += uint64(timeDataSize)
	}

	if err := util.ReadOrError(r, data[:4]); err != nil {
		return err
	} else {
		b.Pad = (uint8(data[0]) >> 7) & 0x1 // 1 bit

		// 5 bits per Language
		b.Language[0] = (uint8(data[0]) >> 2) & 0x1F
		b.Language[1] = ((uint8(data[0]) & 0x3) << 3) | ((uint8(data[1]) >> 5) & 0x7)
		b.Language[2] = uint8(data[1]) & 0x1F

		for i := 0; i < len(b.Language); i++ {
			b.Language[i] += languageCodeOffset // character is packed as the difference between its ASCII value and 0x60.
		}

		b.Predefined = binary.BigEndian.Uint16(data[2:4])

		parsedBytes += 4
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

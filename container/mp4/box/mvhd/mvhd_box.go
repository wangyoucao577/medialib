// Package mvhd represents Movie Header Box.
package mvhd

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

// Box represents a ftyp box.
type Box struct {
	box.FullHeader `json:"full_header"`

	CreationTime     uint64 `json:"creation_time,string"`
	ModificationTime uint64 `json:"modification_time,string"`
	Timescale        uint32 `json:"timescale"`
	Duration         uint64 `json:"duration"`

	Rate   int32 `json:"rate"`
	Volume int16 `json:"volume"`
	// reserved 16 + 2*32 = 80 bits in here
	Matrix      [9]int32  `json:"matrix"`
	PreDefined  [6]uint32 `json:"pre_defined"`
	NextTrackID uint32    `json:"next_track_id"`
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

		Rate        int32     `json:"rate"`
		Volume      int16     `json:"volume"`
		Matrix      [9]int32  `json:"matrix"`
		PreDefined  [6]uint32 `json:"pre_defined"`
		NextTrackID uint32    `json:"next_track_id"`
	}{
		FullHeader: b.FullHeader,

		CreationTime:         time1904.Unix(int64(b.CreationTime), 0).UTC(),
		ModificationTime:     time1904.Unix(int64(b.ModificationTime), 0).UTC(),
		Timescale:            b.Timescale,
		Duration:             b.Duration,
		DurationMilliSeconds: float64(b.Duration) * 1000 / float64(b.Timescale),

		Rate:        b.Rate,
		Volume:      b.Volume,
		Matrix:      b.Matrix,
		PreDefined:  b.PreDefined,
		NextTrackID: b.NextTrackID,
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
	payloadSize := b.PayloadSize() // need

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
		b.Rate = int32(binary.BigEndian.Uint32(data[:4]))
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return err
	} else {
		b.Volume = int16(binary.BigEndian.Uint16(data[:2]))
		parsedBytes += 2
	}

	// ignore reserved 16 + 2*32 = 80 bits in here
	if err := util.ReadOrError(r, make([]byte, 10)); err != nil {
		return err
	} else {
		parsedBytes += 10
	}

	data = data[:4] // shrink to 4 bytes
	for i := 0; i < len(b.Matrix); i++ {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.Matrix[i] = int32(binary.BigEndian.Uint32(data))
			parsedBytes += 4
		}
	}

	for i := 0; i < len(b.PreDefined); i++ {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.PreDefined[i] = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.NextTrackID = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if parsedBytes != payloadSize {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, payloadSize)
	}

	return nil
}

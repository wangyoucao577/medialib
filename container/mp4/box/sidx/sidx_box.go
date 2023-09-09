// Package sidx represents Segment Index Box.
package sidx

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

type Reference struct {
	ReferenceType      uint8  `json:"reference_type"`      // 1bit
	ReferencedSize     uint32 `json:"referenced_size"`     // 31bits
	SubsegmentDuration uint32 `json:"subsegment_duration"` // 32bits
	StartsWithSAP      uint8  `json:"starts_with_SAP"`     // 1bit
	SAPtype            uint8  `json:"SAP_type"`            // 3bits
	SAPDeltaTime       uint32 `json:"SAP_delta_time"`      // 28bits
}

// Box represents a sidx box.
type Box struct {
	box.FullHeader `json:"full_header"`

	ReferenceID              uint32 `json:"reference_ID"`
	Timescale                uint32 `json:"timescale"`
	EarliestPresentationTime uint64 `json:"earliest_presentation_time"`
	FirstOffset              uint64 `json:"first_offset"`

	// 16 bits reserved here
	ReferenceCount uint16      `json:"reference_count"`
	References     []Reference `json:"references"`
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

	data := make([]byte, 8)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.ReferenceID = binary.BigEndian.Uint32(data[:4])
		b.Timescale = binary.BigEndian.Uint32(data[4:])
	}

	if b.FullHeader.Version != 0 {
		data = make([]byte, 16)
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	}

	if b.FullHeader.Version != 0 {
		b.EarliestPresentationTime = binary.BigEndian.Uint64(data[:8])
		b.FirstOffset = binary.BigEndian.Uint64(data[8:])
	} else {
		b.EarliestPresentationTime = uint64(binary.BigEndian.Uint32(data[:4]))
		b.FirstOffset = uint64(binary.BigEndian.Uint32(data[4:]))
	}

	data = data[:4]
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.ReferenceCount = binary.BigEndian.Uint16(data[2:])
	}

	data = make([]byte, 12)
	for i := 0; i < int(b.ReferenceCount); i++ {
		if err := util.ReadOrError(r, data); err != nil {
			return err
		}

		ref := Reference{}
		ref.ReferenceType = (data[0] >> 7) & 0x1
		data[0] &= 0x7F
		ref.ReferencedSize = binary.BigEndian.Uint32(data[:4])
		ref.SubsegmentDuration = binary.BigEndian.Uint32(data[4:8])
		ref.StartsWithSAP = (data[8] >> 7) & 0x1
		ref.SAPtype = (data[8] >> 4) & 0x7
		data[8] &= 0xF
		ref.SAPDeltaTime = binary.BigEndian.Uint32(data[8:])
		b.References = append(b.References, ref)
	}

	return nil
}

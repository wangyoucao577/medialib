package esds

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// DecoderConfigDescriptor represents DecoderConfigDescriptor.
type DecoderConfigDescriptor struct {
	Descriptor Descriptor `json:"descriptor"`

	ObjectTypeIndication uint8 `json:"object_type_indication"`
	StreamType           uint8 `json:"stream_type"` // 6 bits
	UpStream             uint8 `json:"up_stream"`   // 1 bit
	// 1 bit reserved here
	BufferSizeDB uint32 `json:"buffer_size_db"` // 24 bits
	MaxBitrate   uint32 `json:"max_bitrate"`
	AvgBitrate   uint32 `json:"avg_bitrate"`

	DecoderSpecificInfo DecoderSpecificInfo `json:"decoder_specific_info"`
}

func (d *DecoderConfigDescriptor) parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64
	var parsedHeaderBytes uint64
	data := make([]byte, 4)

	// parse descriptor header
	if bytes, err := d.Descriptor.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedHeaderBytes += bytes
		parsedBytes += bytes
	}

	if err := util.ReadOrError(r, data[:1]); err != nil {
		return parsedBytes, err
	} else {
		d.ObjectTypeIndication = uint8(data[0])
		parsedBytes += 1
	}

	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.StreamType = (data[0] >> 2) & 0x3F
		d.UpStream = (data[0] >> 1) & 0x1

		data[0] = 0
		d.BufferSizeDB = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.MaxBitrate = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		d.AvgBitrate = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if d.Descriptor.Size > uint32(parsedBytes-parsedHeaderBytes) {
		if bytes, err := d.DecoderSpecificInfo.parse(r); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
	}

	if parsedBytes-parsedHeaderBytes != uint64(d.Descriptor.Size) {
		glog.Warningf("descriptor %s(%d) still has %d bytes need to parse, parsed payload bytes != payload size: %d != %d", classTagName(d.Descriptor.Tag), d.Descriptor.Tag, uint64(d.Descriptor.Size)-uint64(parsedBytes-parsedHeaderBytes), parsedBytes-parsedHeaderBytes, d.Descriptor.Size)
	}

	return parsedBytes, nil
}

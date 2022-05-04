// Package chunk represents RTMP chunk stream formats.
package chunk

import (
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// Fmt to decide Chunk Message Header Type(0-3)
const (
	MessageHeaderFmt0 = iota
	MessageHeaderFmt1
	MessageHeaderFmt2
	MessageHeaderFmt3
)

// BasicHeader represents RTMP chunk basic header.
type BasicHeader struct {
	Fmt      uint8  `json:"fmt"`       // chunk format, 2 bits
	StreamID uint32 `json:"stream_id"` // chunk stream id, 6 bits or 14 bits or 22 bits
}

const (
	chunkStreamIDThreshold6bits  = 64
	chunkStreamIDThreshold14bits = 320

	chunkStreamIDCalculateBase22bits = 256
)

// Serialize serializes basic header to binary format.
func (b *BasicHeader) Serialize() []byte {
	data := make([]byte, 3)

	data[0] = (b.Fmt << 6)
	if b.StreamID < chunkStreamIDThreshold6bits {
		data[0] |= byte(b.StreamID)
		return data[:1]
	} else if b.StreamID < chunkStreamIDThreshold14bits {
		data[1] = byte(b.StreamID - chunkStreamIDThreshold6bits)
		return data[:2]
	} else {
		data[0] |= 0x1
		data[1] = byte(b.StreamID % chunkStreamIDCalculateBase22bits)
		data[2] = byte(b.StreamID / chunkStreamIDCalculateBase22bits)
		return data[:]
	}
}

// Parse parses basic header from binary format.
func (b *BasicHeader) Parse(r io.Reader) (uint64, error) {
	var parsedBytes uint

	if nextByte, err := util.ReadByteOrError(r); err != nil {
		return uint64(parsedBytes), err
	} else {
		b.Fmt = (nextByte >> 6) & 0x3
		b.StreamID = uint32(nextByte) & 0x3F
		parsedBytes += 1
	}

	if b.StreamID == 0 { // parse 14 bits stream id
		if nextByte, err := util.ReadByteOrError(r); err != nil {
			return uint64(parsedBytes), err
		} else {
			b.StreamID += uint32(nextByte) + chunkStreamIDThreshold6bits
			parsedBytes += 1
		}
	} else if b.StreamID == 1 { // parse 22 bits stream id
		data := make([]byte, 2)
		if err := util.ReadOrError(r, data); err != nil {
			return uint64(parsedBytes), err
		} else {
			b.StreamID = uint32(data[1])*chunkStreamIDCalculateBase22bits +
				uint32(data[1]) + chunkStreamIDThreshold6bits
			parsedBytes += 2
		}
	}

	return uint64(parsedBytes), nil
}

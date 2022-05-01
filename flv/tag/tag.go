// Package tag represents FLV tags.
package tag

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// Header presents FLV tag shared data.
type Header struct {
	// 2 bits reserved here

	// 	Indicates if packets are filtered. 0 = No pre-processing required.
	// 1 = Pre-processing (such as decryption) of the packet is required before it can be rendered.
	// Shall be 0 in unencrypted files, and 1 for encrypted tags. See Annex F. FLV Encryption for the use of filters.
	Filter uint8 `json:"Filter"` // 1 bit

	// 	Type of contents in this tag. The following types are defined:
	// 8 = audio
	// 9 = video
	// 18 = script data
	TagType uint8 `json:"TagType"` // 5 bits

	// Length of the message. Number of bytes after StreamID to
	// end of tag (Equal to length of the tag – 11)
	DataSize uint32 `json:"DataSize"` // 24 bits

	// Time in milliseconds at which the data in this tag applies.
	// This value is relative to the first tag in the FLV file, which
	// always has a timestamp of 0.
	Timestamp24bits uint32 `json:"Timestamp"` // 24 bits

	// Extension of the Timestamp field to form a SI32 value.
	// This field represents the upper 8 bits, while the previous
	// Timestamp field represents the lower 24 bits of the time in milliseconds.
	TimestampExtended uint8 `json:"TimestampExtended"` // 8 bits

	// Timestamp concats TimestampExtended and Timestamp24bits together to represent the real timestamp.
	// In contrast of other fields are stored in byte stream, it's calculated by TimestampExtended and Timestamp24bits.
	Timestamp int32

	StreamID uint32 `json:"StreamID"` // Always 0. 24 bits
}

// Parse parses Tag Header.
func (h *Header) Parse(r io.Reader) error {
	data := make([]byte, 11) // fixed-size header 11 bytes
	if err := util.ReadOrError(r, data); err != nil {
		return err
	}

	h.Filter = (data[0] >> 5) & 0x1
	h.TagType = data[0] & 0x1F
	h.DataSize = binary.BigEndian.Uint32([]byte{0x00, data[1], data[2], data[3]})

	h.Timestamp24bits = binary.BigEndian.Uint32([]byte{0x00, data[4], data[5], data[6]})
	h.TimestampExtended = data[7]
	h.Timestamp = int32((uint32(data[7]) << 24) | h.Timestamp24bits)
	h.StreamID = binary.BigEndian.Uint32([]byte{0x00, data[8], data[9], data[10]})
	return nil
}

// Validate checks whether Header is valid or not.
func (h Header) Validate() error {
	if h.TagType != TypeAudio &&
		h.TagType != TypeVideo &&
		h.TagType != TypeSriptData {
		return fmt.Errorf("unknown tag type %d", h.TagType)
	}

	if h.DataSize == 0 {
		return fmt.Errorf("tag type %d data size is 0", h.TagType)
	}

	return nil
}

package flv

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/wangyoucao577/medialib/util"
)

const (
	HeaderSize = 9 // FLV header is fixed 9 bytes
)

// Header represents FLV Header.
type Header struct {
	Signature string `json:"Signature"` // 3 bytes, always be 'FLV'
	Version   uint8  `json:"Version"`   // File version (for example, 0x01 for FLV version 1)
	// 5 bits reserved here
	TypeFlagsAudio uint8 `json:"TypeFlagsAudio"` // 1 bit, 1 = Audio tags are presents
	// 1 bit reserved here
	TypeFlagsVideo uint8  `json:"TypeFlagsVideo"` // 1 bit, 1 = Video tags are presents
	DataOffset     uint32 `json:"DataOffset"`     // The length of this header in bytes
}

// Parse parses FLV Header.
func (h *Header) Parse(r io.Reader) error {

	data := make([]byte, HeaderSize) // fixed 9 bytes
	if err := util.ReadOrError(r, data); err != nil {
		return err
	}

	h.Signature = string(data[:3])
	h.Version = data[3]
	h.TypeFlagsAudio = (data[4] >> 2) & 0x1
	h.TypeFlagsVideo = data[4] & 0x1
	h.DataOffset = binary.BigEndian.Uint32(data[5:])

	if !strings.EqualFold(h.Signature, "FLV") {
		return fmt.Errorf("invalid signature %s", string(h.Signature))
	}
	if h.DataOffset != HeaderSize {
		return fmt.Errorf("invalid data offset %d, should be %d", h.DataOffset, HeaderSize)
	}

	return nil
}

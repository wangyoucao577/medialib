// Package sei represents AVC Supplemental enhancement information defined in ISO/IEC-14496-10 7.3.2.3.1
package sei

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// SEIMessage represents AVC Supplemental enhancement information.
type SEIMessage struct {
	PayloadType         int   `json:"payload_type"`
	LastPayloadTypeByte uint8 `json:"last_payload_type_byte"`
	PayloadSize         int   `json:"payload_size"`
	LastPayloadSizeByte uint8 `json:"last_payload_size_byte"`

	UserDataUnregistered *UserDataUnregistered `json:"user_data_unregistered,omitempty"`
}

// Parse parses bytes to AVC NAL Unit, return parsed bytes or error.
// The NAL Unit syntax defined in ISO/IEC-14496-10 7.3.1.
func (s *SEIMessage) Parse(r io.Reader, size int) (uint64, error) {
	var parsedBytes uint64

	nextByte := make([]byte, 1)

	// payload type
	for {
		if err := util.ReadOrError(r, nextByte); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += 1
		}

		if nextByte[0] != 0xFF {
			s.LastPayloadTypeByte = nextByte[0]
			s.PayloadType += int(s.LastPayloadTypeByte)
			break
		}
		s.PayloadType += 255
	}

	// payload size
	for {
		if err := util.ReadOrError(r, nextByte); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += 1
		}

		if nextByte[0] != 0xFF {
			s.LastPayloadSizeByte = nextByte[0]
			s.PayloadSize += int(s.LastPayloadSizeByte)
			break
		}
		s.PayloadSize += 255
	}

	switch s.PayloadType {
	case PayloadTypeUserDataUnregistered:
		s.UserDataUnregistered = &UserDataUnregistered{}
		if bytes, err := s.UserDataUnregistered.parse(r, s.PayloadSize); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
	default:
		glog.Warningf("unknown SEI payload type %d, ignore %d bytes", s.PayloadType, s.PayloadSize)
	}

	return parsedBytes, nil
}

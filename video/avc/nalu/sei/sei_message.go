// Package sei represents AVC Supplemental enhancement information defined in ISO/IEC-14496-10 7.3.2.3.1
package sei

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// SEIMessage represents AVC Supplemental enhancement information.
type SEIMessage struct {
	PayloadType         int   `json:"payload_type"`
	LastPayloadTypeByte uint8 `json:"last_payload_type_byte"`
	PayloadSize         int   `json:"payload_size"`
	LastPayloadSizeByte uint8 `json:"last_payload_size_byte"`

	BufferingPeriod      *BufferingPeriod      `json:"buffering_period,omitempty"`
	PicTiming            *PicTiming            `json:"pic_timing,omitempty"`
	UserDataUnregistered *UserDataUnregistered `json:"user_data_unregistered,omitempty"`

	// store for some internal parsing
	sps *sps.SequenceParameterSetData `json:"-"`
	pps *pps.PictureParameterSet      `json:"-"`
}

// SetSequenceHeaders sets both SPS and PPS for parsing.
func (s *SEIMessage) SetSequenceHeaders(sps *sps.SequenceParameterSetData, pps *pps.PictureParameterSet) {
	s.sps = sps
	s.pps = pps
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

	var parser payloadParser
	switch s.PayloadType {
	case PayloadTypeBufferingPeriod:
		s.BufferingPeriod = &BufferingPeriod{}
		s.BufferingPeriod.setSequenceHeaders(s.sps, s.pps)
		parser = s.BufferingPeriod
	case PayloadTypePicTiming:
		s.PicTiming = &PicTiming{}
		s.PicTiming.setSequenceHeaders(s.sps, s.pps)
		parser = s.PicTiming
	case PayloadTypeUserDataUnregistered:
		s.UserDataUnregistered = &UserDataUnregistered{}
		parser = s.UserDataUnregistered
	default:
		glog.Warningf("unknown SEI payload type %d, ignore %d bytes", s.PayloadType, s.PayloadSize)
	}

	if parser != nil {
		if bytes, err := parser.parse(r, s.PayloadSize); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
	}

	return parsedBytes, nil
}

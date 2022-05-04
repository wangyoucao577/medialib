package video

import (
	avcc "github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/avcC"
	"github.com/wangyoucao577/medialib/video/avc/es"
)

type AVCVideoPacket struct {
	AVCDecoderConfigurationRecord *avcc.AVCDecoderConfigurationRecord `json:"avc_config,omitempty"`
	LengthNALU                    []es.LengthNALU                     `json:"length_nalu,omitempty"`
}

// TagBody represents video tag payload.
type TagBody struct {
	AVCVideoPacket *AVCVideoPacket `json:"AVCVideoPacket"`
}

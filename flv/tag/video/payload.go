package video

import (
	avcc "github.com/wangyoucao577/medialib/mp4/box/sampleentry/avcC"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
)

type AVCVideoPacket struct {
	AVCDecoderConfigurationRecord *avcc.AVCDecoderConfigurationRecord `json:"avc_config,omitempty"`
	NALUnits                      []nalu.NALUnit                      `json:"nal_units,omitempty"`
}

// TagBody represents video tag payload.
type TagBody struct {
	AVCVideoPacket *AVCVideoPacket `json:"AVCVideoPacket"`
}

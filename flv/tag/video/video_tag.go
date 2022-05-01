// Package video represents Video Tag.
package video

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv/tag"
	avcc "github.com/wangyoucao577/medialib/mp4/box/sampleentry/avcC"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/es"
)

// Tag represents video tag.
type Tag struct {
	Header         tag.Header `json:"TagHeader"`
	VideoTagHeader TagHeader  `json:"VideoTagHeader"`
	TagBody        *TagBody   `json:"VideoTagBody"`

	avcConfig *avcc.AVCDecoderConfigurationRecord `json:"-"`
}

// SetAVCConfig pass in AVCDecoderConfigurationRecord for slice parsing.
func (t *Tag) SetAVCConfig(avcConfig *avcc.AVCDecoderConfigurationRecord) {
	t.avcConfig = avcConfig
}

// GetTagHeader returns tag header.
func (t *Tag) GetTagHeader() tag.Header {
	return t.Header
}

// ParsePayload parses VideoTagHeader and TayBody data with preset tag.Header.
func (t *Tag) ParsePayload(r io.Reader) error {
	if err := t.Header.Validate(); err != nil {
		return err
	}

	var parsedBytes uint64
	if bytes, err := t.VideoTagHeader.parse(r); err != nil {
		return err
	} else {
		parsedBytes += bytes
	}

	if parsedBytes > uint64(t.Header.DataSize) {
		return fmt.Errorf("tag type %d(%s) data size %d but already parsed %d",
			t.Header.TagType, tag.TypeDescription(int(t.Header.TagType)),
			t.Header.DataSize, parsedBytes)
	}

	if t.VideoTagHeader.CodecID != CodecIDAVC {
		//TODO: parse payload
		glog.Warningf("tag type %d(%s) codec %d(%s) doesn't implemented yet, ignore size %d",
			t.Header.TagType, tag.TypeDescription(int(t.Header.TagType)),
			t.VideoTagHeader.CodecID, CodecIDDescription(int(t.VideoTagHeader.CodecID)),
			t.Header.DataSize-uint32(parsedBytes))
		if err := util.ReadOrError(r, make([]byte, t.Header.DataSize-uint32(parsedBytes))); err != nil {
			return err
		}
		return nil
	}
	t.TagBody = &TagBody{AVCVideoPacket: &AVCVideoPacket{}}

	if *t.VideoTagHeader.AVCPacketType == AVCPacketTypeSequenceHeader {
		t.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord = &avcc.AVCDecoderConfigurationRecord{}
		if bytes, err := t.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.Parse(r); err != nil {
			return err
		} else {
			parsedBytes += bytes
		}
	} else if *t.VideoTagHeader.AVCPacketType == AVCPacketTypeNALU {
		videoES := &es.ElementaryStream{}
		if t.avcConfig != nil {
			videoES.SetLengthSize(t.avcConfig.LengthSize())
			if len(t.avcConfig.LengthSPSNALU) > 0 && len(t.avcConfig.LengthPPSNALU) > 0 {
				videoES.SetSequenceHeaders(t.avcConfig.LengthSPSNALU[0].NALUnit.SequenceParameterSetData,
					t.avcConfig.LengthPPSNALU[0].NALUnit.PictureParameterSet)
			}
		}
		if bytes, err := videoES.Parse(r, int(t.Header.DataSize-uint32(parsedBytes))); err != nil {
			return err
		} else {
			parsedBytes += bytes
		}
		t.TagBody.AVCVideoPacket.LengthNALU = videoES.LengthNALU
	} else {
		glog.Warningf("unhandled avc packet type %d(%s)",
			*t.VideoTagHeader.AVCPacketType, AVCPacketTypeDescription(int(*t.VideoTagHeader.AVCPacketType)))
	}

	if parsedBytes < uint64(t.Header.DataSize) {
		remainBytes := uint64(t.Header.DataSize) - parsedBytes
		glog.Warningf("tag type %d(%s) still has %d bytes NOT parse",
			t.Header.TagType, tag.TypeDescription(int(t.Header.TagType)), remainBytes)
		if err := util.ReadOrError(r, make([]byte, remainBytes)); err != nil {
			return err
		}
	}

	return nil
}

// Package flv reprensents The FLV File Format, defines in Adobe Flash Video File Format Specification Version 10.1 AnnexE.
package flv

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/flv/tag"
	"github.com/wangyoucao577/medialib/container/flv/tag/audio"
	"github.com/wangyoucao577/medialib/container/flv/tag/script"
	"github.com/wangyoucao577/medialib/container/flv/tag/video"
	avcc "github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/avcC"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/annexbes"
	"github.com/wangyoucao577/medialib/video/avc/es"
)

// FLV represents the FLV file format.
type FLV struct {
	Header Header `json:"header"`

	// PreviousTagSize[0] always == 0
	// PreviousTagSize[i+1] = len(Tags[i])
	PreviousTagSize []uint32  `json:"PreviousTagSize"`
	Tags            []tag.Tag `json:"Tags"`
}

// Parse parses FLV data.
func (f *FLV) Parse(r io.Reader) error {
	if err := f.Header.Parse(r); err != nil {
		return err
	}

	var avcConfig *avcc.AVCDecoderConfigurationRecord
	tagSizeData := make([]byte, 4) // fixed 4 bytes
	var lastParsedTagSize int64

	for {
		// parse previous tag size
		if err := util.ReadOrError(r, tagSizeData); err != nil {
			return err
		}
		size := binary.BigEndian.Uint32(tagSizeData)
		f.PreviousTagSize = append(f.PreviousTagSize, size)
		if size != uint32(lastParsedTagSize) {
			glog.Warningf("PreviousTagSize %d != LastParsedTagSize %d", size, lastParsedTagSize)
		}

		// parse tag header
		tagHeader := tag.Header{}
		if err := tagHeader.Parse(r); err != nil {
			return err
		}

		// parse tag data
		var t tag.Tag
		if tagHeader.TagType == tag.TypeAudio {
			t = &audio.Tag{Header: tagHeader}
		} else if tagHeader.TagType == tag.TypeVideo {
			videoTag := &video.Tag{Header: tagHeader}
			videoTag.SetAVCConfig(avcConfig)
			t = videoTag
		} else if tagHeader.TagType == tag.TypeSriptData {
			t = &script.Tag{Header: tagHeader}
		}
		if tagHeader.Filter == 1 {

			//TODO: parse payload
			glog.Warningf("unsupported tag type %d with filter=1, ignore payload size %d", tagHeader.TagType, tagHeader.DataSize)
			if err := util.ReadOrError(r, make([]byte, tagHeader.DataSize)); err != nil {
				return err
			}
			continue
		}

		if err := t.ParsePayload(r); err != nil {
			return err
		}
		f.Tags = append(f.Tags, t)

		lastParsedTagSize = t.Size()                   // cache parsed tag size for checking
		if t.GetTagHeader().TagType == tag.TypeVideo { // cache avcConfig for later slice parsing
			videoTag, ok := t.(*video.Tag)
			if !ok {
				return fmt.Errorf("invalid video tag %v", videoTag)
			}
			if videoTag.VideoTagHeader.AVCPacketType != nil &&
				*videoTag.VideoTagHeader.AVCPacketType == video.AVCPacketTypeSequenceHeader &&
				videoTag.TagBody.AVCVideoPacket != nil &&
				videoTag.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord != nil {
				avcConfig = videoTag.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord
			}
		}
	}
}

// JSON marshals FLV to JSON representation
func (f FLV) JSON() ([]byte, error) {
	return json.Marshal(f)
}

// JSONIndent marshals FLV to JSON representation with customized indent.
func (f FLV) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(f, prefix, indent)
}

// YAML formats FLV to YAML representation.
func (f FLV) YAML() ([]byte, error) {
	j, err := json.Marshal(f)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CSV formats FLV to CSV representation, which isn't supported at the moment.
func (f FLV) CSV() ([]byte, error) {
	return nil, fmt.Errorf("csv representation does not support yet")
}

// ExtractES extracts AVC or HEVC Elementary Stream.
func (f *FLV) ExtractES() (*es.ElementaryStream, error) {

	if len(f.Tags) == 0 {
		return nil, fmt.Errorf("tags not found")
	}

	e := es.ElementaryStream{}
	for _, t := range f.Tags {
		if t.GetTagHeader().TagType != tag.TypeVideo {
			continue
		}

		vt, ok := t.(*video.Tag)
		if !ok {
			return nil, fmt.Errorf("tag %#v should be video tag but cannot convert", t)
		}

		if vt.VideoTagHeader.AVCPacketType == nil || vt.TagBody == nil {
			return nil, fmt.Errorf("tag %#v empty AVCPacketType or TagBody", t)
		}
		if *vt.VideoTagHeader.AVCPacketType == video.AVCPacketTypeSequenceHeader {
			if vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord == nil ||
				len(vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthSPSNALU) == 0 ||
				len(vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthPPSNALU) == 0 {
				return nil, fmt.Errorf("tag %#v expect avc config but empty", t)
			}
			e.SetLengthSize(vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthSize())
			e.SetSequenceHeaders(vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthSPSNALU[0].NALUnit.SequenceParameterSetData,
				vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthPPSNALU[0].NALUnit.PictureParameterSet)

			// also add sps,pps nalu since they're real NALU in flv tag
			var spsLengthNALU, ppsLengthNALU []es.LengthNALU
			for _, n := range vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthSPSNALU {
				spsLengthNALU = append(spsLengthNALU, es.LengthNALU{Length: uint32(n.Length), NALU: n.NALUnit})
			}
			for _, n := range vt.TagBody.AVCVideoPacket.AVCDecoderConfigurationRecord.LengthPPSNALU {
				ppsLengthNALU = append(ppsLengthNALU, es.LengthNALU{Length: uint32(n.Length), NALU: n.NALUnit})
			}
			e.LengthNALU = append(e.LengthNALU, spsLengthNALU...)
			e.LengthNALU = append(e.LengthNALU, ppsLengthNALU...)

		} else if *vt.VideoTagHeader.AVCPacketType == video.AVCPacketTypeNALU {
			if len(vt.TagBody.AVCVideoPacket.LengthNALU) == 0 {
				return nil, fmt.Errorf("tag %#v expect nal units but empty", t)
			}
			e.LengthNALU = append(e.LengthNALU, vt.TagBody.AVCVideoPacket.LengthNALU...)
		} // video.AVCPacketTypeEOS doesn't need to handle

	}

	return &e, nil
}

// ExtractAnnexBES extracts AVC or HEVC Elementary Stream with AnnexB byte format.
func (f FLV) ExtractAnnexBES() (*annexbes.ElementaryStream, error) {
	mp4ES, err := f.ExtractES()
	if err != nil {
		return nil, err
	}

	annexbES := annexbes.ElementaryStream{}
	for i := range mp4ES.LengthNALU {
		annexbES.NALU = append(annexbES.NALU, mp4ES.LengthNALU[i].NALU)
	}

	return &annexbES, nil
}

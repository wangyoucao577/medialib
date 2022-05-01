// Package video represents Video Tag.
package video

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv/tag"
	"github.com/wangyoucao577/medialib/util"
)

// TagHeader reprensets Audio Tag Header.
type TagHeader struct {

	// Type of video frame. The following values are defined:
	// 1 = key frame (for AVC, a seekable frame)
	// 2 = inter frame (for AVC, a non-seekable frame)
	// 3 = disposable inter frame (H.263 only)
	// 4 = generated key frame (reserved for server use only)
	// 5 = video info/command frame
	FrameType uint8 `json:"FrameType"` // 4 bits

	// 	Codec Identifier. The following values are defined:
	// 2 = Sorenson H.263
	// 3 = Screen video
	// 4 = On2 VP6
	// 5 = On2 VP6 with alpha channel 6 = Screen video version 2
	// 7 = AVC
	CodecID uint8 `json:"CodecID"` // 4 bits

	// The following values are defined:
	// 0 = AVC sequence header
	// 1 = AVC NALU
	// 2 = AVC end of sequence (lower level NALU sequence ender is not required or supported)
	AVCPacketType *uint8 `json:"AVCPacketType,omitempty"`

	// IF AVCPacketType == 1
	//   Composition time offset
	// ELSE
	//   0
	// See ISO 14496-12, 8.15.3 for an explanation of composition times.
	// The offset in an FLV file is always in milliseconds.
	CompositionTime *int32 `json:"CompositionTime,omitempty"`
}

// Tag represents video tag.
type Tag struct {
	Header         tag.Header `json:"TagHeader"`
	VideoTagHeader TagHeader  `json:"VideoTagHeader"`
	TagBody        *TagBody   `json:"VideoTagBody"`
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

	//TODO: parse payload
	glog.Warningf("tag type %d doesn't implemented yet, ignore payload size %d", t.Header.TagType, t.Header.DataSize)
	if err := util.ReadOrError(r, make([]byte, t.Header.DataSize)); err != nil {
		return err
	}
	return nil
}

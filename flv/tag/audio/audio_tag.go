// Package audio represents Audio Tag.
package audio

import (
	"fmt"
	"io"

	"github.com/wangyoucao577/medialib/flv/tag"
	"github.com/wangyoucao577/medialib/util"
)

// Tag represents audio tag.
type Tag struct {
	Header         tag.Header `json:"TagHeader"`
	AudioTagHeader TagHeader  `json:"AudioTagHeader"`
	Body           *TagBody   `json:"AudioTagBody"`
}

// GetTagHeader returns tag header.
func (t *Tag) GetTagHeader() tag.Header {
	return t.Header
}

// ParsePayload parses AudioTagHeader and TayBody data with preset tag.Header.
func (t *Tag) ParsePayload(r io.Reader) error {
	if err := t.Header.Validate(); err != nil {
		return err
	}

	var parsedBytes uint64
	if bytes, err := t.AudioTagHeader.parse(r); err != nil {
		return err
	} else {
		parsedBytes += bytes
	}

	if parsedBytes > uint64(t.Header.DataSize) {
		return fmt.Errorf("tag type %d(%s) data size %d but already parsed %d",
			t.Header.TagType, tag.TypeDescription(int(t.Header.TagType)),
			t.Header.DataSize, parsedBytes)
	}
	remainBytes := uint64(t.Header.DataSize) - parsedBytes

	data := make([]byte, remainBytes)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	}

	// MUST AAC here, and MUST have AACPacketType
	aacAudioData := &AACAudioData{}
	if *t.AudioTagHeader.AACPacketType == AACPacketTypeSequenceHeader {
		aacAudioData.AudioSpecificConfig = data
	} else {
		aacAudioData.RawAACFrameData = data
	}
	t.Body = &TagBody{AACAudioData: aacAudioData}

	return nil
}

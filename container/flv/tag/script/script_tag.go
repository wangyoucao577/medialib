// Package script represents FLV script tag.
package script

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/flv/tag"
	"github.com/wangyoucao577/medialib/util"
)

type Tag struct {
	Header  tag.Header `json:"TagHeader"`
	TagBody *TagBody   `json:"TagBody,omitempty"`
}

// GetTagHeader returns tag header.
func (t *Tag) GetTagHeader() tag.Header {
	return t.Header
}

// Size returns total bytes of the tag, equal to (HeaderSize(11bytes) + DataSize)
func (t Tag) Size() int64 {
	return int64(t.Header.DataSize) + tag.HeaderSize
}

// ParsePayload parses VideoTagHeader and TayBody data with preset tag.Header.
func (t *Tag) ParsePayload(r io.Reader) error {
	if err := t.Header.Validate(); err != nil {
		return err
	}

	var parsedBytes uint64

	t.TagBody = &TagBody{}
	if bytes, err := t.TagBody.parse(r); err != nil {
		return err
	} else {
		parsedBytes += bytes
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

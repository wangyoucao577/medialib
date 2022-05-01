// Package script represents FLV script tag.
package script

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv/tag"
	"github.com/wangyoucao577/medialib/util"
)

type Tag struct {
	Header tag.Header `json:"TagHeader"`
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

	//TODO: parse payload
	glog.Warningf("tag type %d doesn't implemented yet, ignore payload size %d", t.Header.TagType, t.Header.DataSize)
	if err := util.ReadOrError(r, make([]byte, t.Header.DataSize)); err != nil {
		return err
	}

	return nil
}

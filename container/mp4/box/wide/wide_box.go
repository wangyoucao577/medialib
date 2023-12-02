// Package wide represents Wide Box.
// https://developer.apple.com/documentation/quicktime-file-format/wide_atom
package wide

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
)

// Box represents a wide box.
type Box struct {
	box.Header `json:"header"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	return nil
}

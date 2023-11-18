// Package sdtp represents Independent and Disposable Samples Box.
package sdtp

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a sdtp box.
type Box struct {
	box.FullHeader `json:"full_header"`

	//TODO: payloads
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		FullHeader: box.FullHeader{
			Header: h,
		},
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	// parse full header additional information first
	if err := b.FullHeader.ParseVersionFlag(r); err != nil {
		return err
	}

	if b.PayloadSize() > 0 {
		//TODO: parse payload
		glog.Warningf("sdtp payload size %d but ignoring", b.PayloadSize())
		if err := util.ReadOrError(r, make([]byte, b.PayloadSize())); err != nil {
			return err
		}
	}

	return nil
}

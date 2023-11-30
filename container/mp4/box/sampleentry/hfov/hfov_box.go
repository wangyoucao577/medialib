// Package hfov represents hfov box.
package hfov

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a hfov box.
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

	//TODO: these can be removed if we have supported full AVCDecoderConfigurationRecord parsing
	glog.Warningf("box type %s still has %d bytes hasn't been parsed yet, ignore them", b.Type, b.PayloadSize())
	if err := util.ReadOrError(r, make([]byte, b.PayloadSize())); err != nil {
		return err
	}

	return nil
}

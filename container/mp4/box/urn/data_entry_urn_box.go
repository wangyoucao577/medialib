// Package urn represents Data Entry Urn Box.
package urn

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a urn box.
type Box struct {
	box.FullHeader `json:"full_header"`
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

	glog.Warningf("box type %s payload bytes %d parsing TODO", b.Type, b.PayloadSize())
	//TODO: parse payload
	if err := util.ReadOrError(r, make([]byte, b.PayloadSize())); err != nil {
		return err
	}

	return nil
}

// Package mdat represents Media Data Box.
package mdat

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a mdat box.
type Box struct {
	box.Header `json:"header"`

	Data []byte `json:"-"`
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

	payloadSize := b.PayloadSize()

	if payloadSize == 0 {
		return fmt.Errorf("TODO: box %s payload size 0, need to read until EOF", b.Type)
	}

	b.Data = make([]byte, payloadSize)
	if err := util.ReadOrError(r, b.Data); err != nil {
		return err
	}

	return nil
}

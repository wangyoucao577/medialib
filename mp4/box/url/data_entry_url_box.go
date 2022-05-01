// Package url represents Data Entry Url Box.
package url

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a url box.
type Box struct {
	box.FullHeader `json:"full_header"`

	Location string `json:"location"`
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
		data := make([]byte, b.PayloadSize())
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			b.Location = string(data)
		}
	}

	return nil
}

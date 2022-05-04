// Package esds represents ES Descriptor Box.
// Most of ES_Descriptor and it's sub descriptors were described in ISO/IEC-14496-1.
// However, I reference a lot from below chromium implementation for many details:
// 	https://chromium.googlesource.com/chromium/src/media/+/16ba1c56b860d53d7354c0ec9538650cf1f20e2d/mp4/es_descriptor.cc
package esds

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a esds box.
type Box struct {
	box.FullHeader `json:"full_header"`

	ESDescriptor `json:"es_descriptor"`
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

	var parsedBytes uint64

	if bytes, err := b.ESDescriptor.parse(r); err != nil {
		return err
	} else {
		parsedBytes += bytes
	}

	if b.PayloadSize() > parsedBytes {
		glog.Warningf("box type %s remain bytes %d parsing TODO", b.Type, b.PayloadSize()-parsedBytes)
		//TODO: parse payload
		if err := util.ReadOrError(r, make([]byte, b.PayloadSize()-parsedBytes)); err != nil {
			return err
		}
	}

	return nil
}

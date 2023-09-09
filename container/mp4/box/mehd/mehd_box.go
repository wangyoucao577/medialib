// Package mehd represents Movie Extends Header Box.
package mehd

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a mehd box.
type Box struct {
	box.FullHeader `json:"full_header"`

	FragmentDuration uint64 `json:"fragment_duration"`
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

	if b.Version == 1 {
		data := make([]byte, 8)
		if err := util.ReadOrError(r, data); err != nil {
			return err
		}
		b.FragmentDuration = binary.BigEndian.Uint64(data)
	} else {
		data := make([]byte, 4)
		if err := util.ReadOrError(r, data); err != nil {
			return err
		}
		b.FragmentDuration = uint64(binary.BigEndian.Uint32(data))
	}

	return nil
}

// Package data represents Data Box (inside ilst).
// https://developer.apple.com/documentation/quicktime-file-format/data_atom
package data

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a data box.
type Box struct {
	box.Header `json:"header"`

	Type     uint32 `json:"type"`
	Language uint32 `json:"lang"`
	Value    string `json:"value"`
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
		glog.Warningf("box %s invalid, err %v", b.Header.Type, err)
		return nil
	}

	// start to parse payload
	var parsedBytes uint64

	data := make([]byte, 8)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.Type = binary.BigEndian.Uint32(data[:4])
		b.Language = binary.BigEndian.Uint32(data[4:])
		parsedBytes += 8
	}

	data = make([]byte, b.PayloadSize()-parsedBytes)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.Value = string(data)
	}

	return nil
}

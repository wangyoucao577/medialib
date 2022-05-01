// Package smhd represents Sound Media Header.
package smhd

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a smhd box.
type Box struct {
	box.FullHeader `json:"full_header"`

	Balance int16 `json:"balance"`
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
	data := make([]byte, 2)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.Balance = int16(binary.BigEndian.Uint16(data))
		parsedBytes += 2
	}

	// ignore reserved 2 bytes
	if err := util.ReadOrError(r, make([]byte, 2)); err != nil {
		return err
	} else {
		parsedBytes += 2
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

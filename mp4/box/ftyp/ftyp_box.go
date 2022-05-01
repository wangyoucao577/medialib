// Package ftyp defines File Type Box.
package ftyp

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a ftyp box.
type Box struct {
	box.Header `json:"header"`

	MajorBrand       box.FixedArray4Bytes   `json:"major_brand"`
	MinorVersion     uint32                 `json:"minor_version"`
	CompatibleBrands []box.FixedArray4Bytes `json:"compatible_brands"`
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

	var parsedBytes uint32
	payloadSize := b.PayloadSize()

	if err := util.ReadOrError(r, b.MajorBrand[:]); err != nil {
		return err
	} else {
		parsedBytes += 4
	}

	minorVersionData := make([]byte, 4)
	if err := util.ReadOrError(r, minorVersionData); err != nil {
		return err
	} else {
		b.MinorVersion = binary.BigEndian.Uint32(minorVersionData)
		parsedBytes += 4
	}

	for parsedBytes < uint32(payloadSize) {
		var data [4]byte
		if err := util.ReadOrError(r, data[:]); err != nil {
			return err
		} else {
			b.CompatibleBrands = append(b.CompatibleBrands, data)
			parsedBytes += 4
		}
	}

	if parsedBytes != uint32(payloadSize) {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, payloadSize)
	}

	return nil
}

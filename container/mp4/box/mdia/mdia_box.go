// Package mdia represents Media Box.
package mdia

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/hdlr"
	"github.com/wangyoucao577/medialib/container/mp4/box/mdhd"
	"github.com/wangyoucao577/medialib/container/mp4/box/minf"
)

// Box represents a mdia box.
type Box struct {
	box.Header `json:"header"`

	Mdhd *mdhd.Box `json:"mdhd,omitempty"`
	Hdlr *hdlr.Box `json:"hdlr,omitempty"`
	Minf *minf.Box `json:"minf,omitempty"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,

		boxesCreator: map[string]box.NewFunc{
			box.TypeMdhd: mdhd.New,
			box.TypeHdlr: hdlr.New,
			box.TypeMinf: minf.New,
		},
	}
}

// CreateSubBox tries to create sub level box.
func (b *Box) CreateSubBox(h box.Header) (box.Box, error) {
	creator, ok := b.boxesCreator[h.Type.String()]
	if !ok {
		glog.V(2).Infof("unknown box type %s, size %d payload %d", h.Type.String(), h.Size, h.PayloadSize())
		return nil, box.ErrUnknownBoxType
	}

	createdBox := creator(h)
	if createdBox == nil {
		glog.Fatalf("create box type %s failed", h.Type.String())
	}

	switch h.Type.String() {
	case box.TypeMdhd:
		b.Mdhd = createdBox.(*mdhd.Box)
	case box.TypeHdlr:
		b.Hdlr = createdBox.(*hdlr.Box)
	case box.TypeMinf:
		b.Minf = createdBox.(*minf.Box)
		// handler_type is required in stsd
		if b.Hdlr != nil {
			b.Minf.SetHdlr(b.Hdlr)
		}
	}
	return createdBox, nil
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	var parsedBytes uint64
	for {
		readBytes, err := box.ParseBox(r, b, b.PayloadSize()-parsedBytes)
		if err != nil {
			if err == io.EOF {
				return err
			} else if err == box.ErrUnknownBoxType || err == box.ErrInsufficientSize {
				// after ignore the box, continue to parse next
			} else {
				return err
			}
		}
		parsedBytes += readBytes

		if parsedBytes == b.PayloadSize() {
			break
		}
	}

	return nil
}

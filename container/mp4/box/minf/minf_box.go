// Package minf represents Media Information Box.
package minf

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/dinf"
	"github.com/wangyoucao577/medialib/container/mp4/box/hdlr"
	"github.com/wangyoucao577/medialib/container/mp4/box/smhd"
	"github.com/wangyoucao577/medialib/container/mp4/box/stbl"
	"github.com/wangyoucao577/medialib/container/mp4/box/vmhd"
)

// Box represents a minf box.
type Box struct {
	box.Header `json:"header"`

	Stbl *stbl.Box `json:"stbl,omitempty"`
	Dinf *dinf.Box `json:"dinf,omitempty"`
	Smhd *smhd.Box `json:"smhd,omitempty"`
	Vmhd *vmhd.Box `json:"vmhd,omitempty"`

	Hdlr       *hdlr.Box `json:"hdlr,omitempty"` // parsed on the fly
	hdlrPassed *hdlr.Box `json:"-"`              // passed from parent for later use

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,

		boxesCreator: map[string]box.NewFunc{
			box.TypeStbl: stbl.New,
			box.TypeDinf: dinf.New,
			box.TypeSmhd: smhd.New,
			box.TypeVmhd: vmhd.New,
			box.TypeHdlr: hdlr.New,
		},
	}
}

// SetHdlr passes hdlr box for later use.
func (b *Box) SetHdlr(h *hdlr.Box) {
	b.hdlrPassed = h
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
	case box.TypeStbl:
		b.Stbl = createdBox.(*stbl.Box)
		// handler_type is required in stsd
		if b.hdlrPassed != nil {
			b.Stbl.SetHdlr(b.hdlrPassed)
		} else if b.Hdlr != nil {
			b.Stbl.SetHdlr(b.Hdlr)
		}
	case box.TypeDinf:
		b.Dinf = createdBox.(*dinf.Box)
	case box.TypeSmhd:
		b.Smhd = createdBox.(*smhd.Box)
	case box.TypeVmhd:
		b.Vmhd = createdBox.(*vmhd.Box)
	case box.TypeHdlr:
		b.Hdlr = createdBox.(*hdlr.Box)
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

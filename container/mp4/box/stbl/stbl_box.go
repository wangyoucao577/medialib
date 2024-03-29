// Package stbl represents Sample Table Box.
package stbl

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/ctts"
	"github.com/wangyoucao577/medialib/container/mp4/box/hdlr"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/sgpd"
	"github.com/wangyoucao577/medialib/container/mp4/box/sdtp"
	"github.com/wangyoucao577/medialib/container/mp4/box/stco"
	"github.com/wangyoucao577/medialib/container/mp4/box/stsc"
	"github.com/wangyoucao577/medialib/container/mp4/box/stsd"
	"github.com/wangyoucao577/medialib/container/mp4/box/stss"
	"github.com/wangyoucao577/medialib/container/mp4/box/stsz"
	"github.com/wangyoucao577/medialib/container/mp4/box/stts"
)

// Box represents a stbl box.
type Box struct {
	box.Header `json:"header"`

	Stsd *stsd.Box  `json:"stsd,omitempty"`
	Stts *stts.Box  `json:"stts,omitempty"`
	Stss *stss.Box  `json:"stss,omitempty"`
	Stsc *stsc.Box  `json:"stsc,omitempty"`
	Stsz *stsz.Box  `json:"stsz,omitempty"`
	Stco *stco.Box  `json:"stco,omitempty"`
	Ctts *ctts.Box  `json:"ctts,omitempty"`
	Sdtp *sdtp.Box  `json:"sdtp,omitempty"`
	Sgpd []sgpd.Box `json:"sgpd,omitempty"`

	// passed from parent for later use
	hdlr *hdlr.Box `json:"-"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// SetHdlr passes hdlr box for later use.
func (b *Box) SetHdlr(h *hdlr.Box) {
	b.hdlr = h
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,

		boxesCreator: map[string]box.NewFunc{
			box.TypeStsd: stsd.New,
			box.TypeStts: stts.New,
			box.TypeStss: stss.New,
			box.TypeStsc: stsc.New,
			box.TypeStsz: stsz.New,
			box.TypeStco: stco.New,
			box.TypeCtts: ctts.New,
			box.TypeSdtp: sdtp.New,
			box.TypeSgpd: sgpd.New,
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
	case box.TypeStsd:
		b.Stsd = createdBox.(*stsd.Box)
		if b.hdlr != nil {
			b.Stsd.SetHdlr(b.hdlr) // requires handler_type for payload parsing
		}
	case box.TypeStts:
		b.Stts = createdBox.(*stts.Box)
	case box.TypeStss:
		b.Stss = createdBox.(*stss.Box)
	case box.TypeStsc:
		b.Stsc = createdBox.(*stsc.Box)
	case box.TypeStsz:
		b.Stsz = createdBox.(*stsz.Box)
	case box.TypeStco:
		b.Stco = createdBox.(*stco.Box)
	case box.TypeCtts:
		b.Ctts = createdBox.(*ctts.Box)
	case box.TypeSdtp:
		b.Sdtp = createdBox.(*sdtp.Box)
	case box.TypeSgpd:
		b.Sgpd = append(b.Sgpd, *createdBox.(*sgpd.Box))
		createdBox = &b.Sgpd[len(b.Sgpd)-1] // reference to the last empty box
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

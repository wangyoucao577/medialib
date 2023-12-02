// Package hev1 represents HEVC Sample Entry.
package hev1

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/btrt"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/colr"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/hfov"
	hvcc "github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/hvcC"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/vexu"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/vide"
)

// HEVCSampleEntry defined HEVC SampleEntry box.
type HEVCSampleEntry struct {
	vide.VisualSampleEntry

	HvccConfig *hvcc.HEVCConfigrationBox `json:"hvcC"`

	LhvcConfig *hvcc.HEVCConfigrationBox `json:"lhvC,omitempty"`
	Colr       *colr.Box                 `json:"colr,omitempty"`
	Hfov       *hfov.Box                 `json:"hfov,omitempty"`
	Vexu       *vexu.Box                 `json:"vexu,omitempty"`

	Btrt *btrt.Box `json:"btrt,omitempty"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &HEVCSampleEntry{
		VisualSampleEntry: vide.VisualSampleEntry{
			SampleEntry: sampleentry.SampleEntry{
				Header: h,
			},
		},

		boxesCreator: map[string]box.NewFunc{
			box.TypeHvcC: hvcc.New,
			box.TypeLhvC: hvcc.New,
			box.TypeColr: colr.New,
			box.TypeHfov: hfov.New,
			box.TypeVexu: vexu.New,
			box.TypeBtrt: btrt.New,
		},
	}
}

// CreateSubBox tries to create sub level box.
func (a *HEVCSampleEntry) CreateSubBox(h box.Header) (box.Box, error) {
	creator, ok := a.boxesCreator[h.Type.String()]
	if !ok {
		glog.V(2).Infof("unknown box type %s, size %d payload %d", h.Type.String(), h.Size, h.PayloadSize())
		return nil, box.ErrUnknownBoxType
	}

	createdBox := creator(h)
	if createdBox == nil {
		glog.Fatalf("create box type %s failed", h.Type.String())
	}

	switch h.Type.String() {
	case box.TypeHvcC:
		a.HvccConfig = createdBox.(*hvcc.HEVCConfigrationBox)
	case box.TypeLhvC:
		a.LhvcConfig = createdBox.(*hvcc.HEVCConfigrationBox)
	case box.TypeColr:
		a.Colr = createdBox.(*colr.Box)
	case box.TypeHfov:
		a.Hfov = createdBox.(*hfov.Box)
	case box.TypeVexu:
		a.Vexu = createdBox.(*vexu.Box)
	case box.TypeBtrt:
		a.Btrt = createdBox.(*btrt.Box)
	}

	return createdBox, nil
}

// ParsePayload parse payload which requires basic box already exist.
func (a *HEVCSampleEntry) ParsePayload(r io.Reader) error {
	if err := a.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", a.Type, err)
		return nil
	}

	// parse VisualSampleEntryFignore unkown box type sdtp sizeFignore unkown box type sdtp size
	if err := a.VisualSampleEntry.ParsePayload(r); err != nil {
		return err
	}

	var parsedBytes uint64
	for {
		readBytes, err := box.ParseBox(r, a, a.PayloadSize()-parsedBytes)
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

		if parsedBytes == a.PayloadSize() {
			break
		}
	}

	return nil
}

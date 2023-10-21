// Package hev1 represents HEVC Sample Entry.
package hev1

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/btrt"
	hvcc "github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/hvcC"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/vide"
)

// HEVCSampleEntry defined HEVC SampleEntry box.
type HEVCSampleEntry struct {
	vide.VisualSampleEntry

	HEVCConfig *hvcc.HEVCConfigrationBox `json:"hvcC"`
	Btrt       *btrt.Box                 `json:"btrt,omitempty"`

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
			box.TypehvcC: hvcc.New,
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
	case box.TypehvcC:
		a.HEVCConfig = createdBox.(*hvcc.HEVCConfigrationBox)
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
		boxHeader, err := box.ParseBox(r, a)
		if err != nil {
			if err == io.EOF {
				return err
			} else if err == box.ErrUnknownBoxType {
				// after ignore the box, continue to parse next
			} else {
				return err
			}
		}
		parsedBytes += boxHeader.BoxSize()

		if parsedBytes == a.PayloadSize() {
			break
		}
	}
	return nil
}

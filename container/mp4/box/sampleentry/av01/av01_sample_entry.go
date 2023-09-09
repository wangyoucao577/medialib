// Package av01 represents AV1 Sample Entry.
package av01

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry"
	av1c "github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/av1C"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/btrt"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/vide"
)

// AV1SampleEntry defined AV1SampleEntry box (https://aomediacodec.github.io/av1-isobmff).
type AV1SampleEntry struct {
	vide.VisualSampleEntry

	AV1Config *av1c.AV1ConfigrationBox `json:"av1C"`
	Btrt      *btrt.Box                `json:"btrt,omitempty"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &AV1SampleEntry{
		VisualSampleEntry: vide.VisualSampleEntry{
			SampleEntry: sampleentry.SampleEntry{
				Header: h,
			},
		},

		boxesCreator: map[string]box.NewFunc{
			box.TypeAv1C: av1c.New,
			box.TypeBtrt: btrt.New,
		},
	}
}

// CreateSubBox tries to create sub level box.
func (a *AV1SampleEntry) CreateSubBox(h box.Header) (box.Box, error) {
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
	case box.TypeAv1C:
		a.AV1Config = createdBox.(*av1c.AV1ConfigrationBox)
	case box.TypeBtrt:
		a.Btrt = createdBox.(*btrt.Box)
	}

	return createdBox, nil
}

// ParsePayload parse payload which requires basic box already exist.
func (a *AV1SampleEntry) ParsePayload(r io.Reader) error {
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

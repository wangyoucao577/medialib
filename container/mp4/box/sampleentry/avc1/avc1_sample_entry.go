// Package avc1 represents AVC Sample Entry.
package avc1

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry"
	avcc "github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/avcC"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/btrt"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/vide"
)

// AVCSampleEntry defined AVCSampleEntry box (ISO/IEC_14496-15 2012 5.3.4).
type AVCSampleEntry struct {
	vide.VisualSampleEntry

	AVCConfig *avcc.AVCConfigrationBox `json:"avcC"`
	Btrt      *btrt.Box                `json:"btrt,omitempty"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &AVCSampleEntry{
		VisualSampleEntry: vide.VisualSampleEntry{
			SampleEntry: sampleentry.SampleEntry{
				Header: h,
			},
		},

		boxesCreator: map[string]box.NewFunc{
			box.TypeAvcC: avcc.New,
			box.TypeBtrt: btrt.New,
		},
	}
}

// CreateSubBox tries to create sub level box.
func (a *AVCSampleEntry) CreateSubBox(h box.Header) (box.Box, error) {
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
	case box.TypeAvcC:
		a.AVCConfig = createdBox.(*avcc.AVCConfigrationBox)
	case box.TypeBtrt:
		a.Btrt = createdBox.(*btrt.Box)
	}

	return createdBox, nil
}

// ParsePayload parse payload which requires basic box already exist.
func (a *AVCSampleEntry) ParsePayload(r io.Reader) error {
	if err := a.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", a.Type, err)
		return nil
	}

	// parse VisualSampleEntry
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

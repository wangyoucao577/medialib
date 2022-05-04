// Package mp4a represents MP4 Visual Sample Entry.
package mp4a

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/esds"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/soun"
)

// MP4VisualSampleEntry defined MP4VisualSampleEntry box (ISO/IEC_14496-14 2003 5.6).
type MP4VisualSampleEntry struct {
	soun.AudioSampleEntry

	Esds *esds.Box `json:"esds"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &MP4VisualSampleEntry{
		AudioSampleEntry: soun.AudioSampleEntry{
			SampleEntry: sampleentry.SampleEntry{
				Header: h,
			},
		},

		boxesCreator: map[string]box.NewFunc{
			box.TypeEsds: esds.New,
		},
	}
}

// CreateSubBox tries to create sub level box.
func (m *MP4VisualSampleEntry) CreateSubBox(h box.Header) (box.Box, error) {
	creator, ok := m.boxesCreator[h.Type.String()]
	if !ok {
		glog.V(2).Infof("unknown box type %s, size %d payload %d", h.Type.String(), h.Size, h.PayloadSize())
		return nil, box.ErrUnknownBoxType
	}

	createdBox := creator(h)
	if createdBox == nil {
		glog.Fatalf("create box type %s failed", h.Type.String())
	}

	switch h.Type.String() {
	case box.TypeEsds:
		m.Esds = createdBox.(*esds.Box)
	}

	return createdBox, nil
}

// ParsePayload parse payload which requires basic box already exist.
func (m *MP4VisualSampleEntry) ParsePayload(r io.Reader) error {
	if err := m.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", m.Type, err)
		return nil
	}

	// parse VisualSampleEntry
	if err := m.AudioSampleEntry.ParsePayload(r); err != nil {
		return err
	}

	var parsedBytes uint64
	for {
		boxHeader, err := box.ParseBox(r, m)
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

		if parsedBytes == m.PayloadSize() {
			break
		}
	}
	return nil
}

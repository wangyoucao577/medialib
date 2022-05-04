// Package traf represents Track Fragment Box.
package traf

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/tfdt"
	"github.com/wangyoucao577/medialib/container/mp4/box/tfhd"
	"github.com/wangyoucao577/medialib/container/mp4/box/trun"
)

// Box represents a traf box.
type Box struct {
	box.Header `json:"header"`

	Tfhd *tfhd.Box  `json:"tfhd"`
	Trun []trun.Box `json:"trun"`
	Tfdt *tfdt.Box  `json:"tfdt,omitempty"`

	boxesCreator map[string]box.NewFunc
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,

		boxesCreator: map[string]box.NewFunc{
			box.TypeTfhd: tfhd.New,
			box.TypeTrun: trun.New,
			box.TypeTfdt: tfdt.New,
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
	case box.TypeTfhd:
		b.Tfhd = createdBox.(*tfhd.Box)
	case box.TypeTrun:
		b.Trun = append(b.Trun, *createdBox.(*trun.Box))
		createdBox = &b.Trun[len(b.Trun)-1]
	case box.TypeTfdt:
		b.Tfdt = createdBox.(*tfdt.Box)
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
		boxHeader, err := box.ParseBox(r, b)
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

		if parsedBytes == b.PayloadSize() {
			break
		}
	}

	return nil
}

// Package stsd represents Sample Description Box.
package stsd

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/hdlr"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/av01"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/avc1"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/hev1"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry/mp4a"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a stsd box.
type Box struct {
	box.FullHeader `json:"full_header"`

	EntryCount             uint32                      `json:"entry_connt"`
	AVC1SampleEntries      []avc1.AVCSampleEntry       `json:"avc1,omitempty"`
	HEV1SampleEntries      []hev1.HEVCSampleEntry      `json:"hev1,omitempty"`
	HVC1SampleEntries      []hev1.HEVCSampleEntry      `json:"hvc1,omitempty"`
	AV01SampleEntries      []av01.AV1SampleEntry       `json:"av01,omitempty"`
	MP4VisualSampleEntries []mp4a.MP4VisualSampleEntry `json:"mp4a,omitempty"`

	// passed from parent for later use
	hdlr *hdlr.Box `json:"-"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		FullHeader: box.FullHeader{
			Header: h,
		},

		boxesCreator: map[string]box.NewFunc{
			box.TypeAvc1: avc1.New,
			box.TypeHev1: hev1.New,
			box.TypeHvc1: hev1.New,
			box.TypeAv01: av01.New,
			box.TypeMp4a: mp4a.New,
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

	switch b.hdlr.HandlerType.String() {
	case box.TypeVide:
		switch h.Type.String() {
		case box.TypeAvc1:
			b.AVC1SampleEntries = append(b.AVC1SampleEntries, *createdBox.(*avc1.AVCSampleEntry))
			createdBox = &b.AVC1SampleEntries[len(b.AVC1SampleEntries)-1]
		case box.TypeHev1:
			b.HEV1SampleEntries = append(b.HEV1SampleEntries, *createdBox.(*hev1.HEVCSampleEntry))
			createdBox = &b.HEV1SampleEntries[len(b.HEV1SampleEntries)-1]
		case box.TypeHvc1:
			b.HVC1SampleEntries = append(b.HVC1SampleEntries, *createdBox.(*hev1.HEVCSampleEntry))
			createdBox = &b.HVC1SampleEntries[len(b.HVC1SampleEntries)-1]
		case box.TypeAv01:
			b.AV01SampleEntries = append(b.AV01SampleEntries, *createdBox.(*av01.AV1SampleEntry))
			createdBox = &b.AV01SampleEntries[len(b.AV01SampleEntries)-1]
		}
	case box.TypeSoun:
		switch h.Type.String() {
		case box.TypeMp4a:
			b.MP4VisualSampleEntries = append(b.MP4VisualSampleEntries, *createdBox.(*mp4a.MP4VisualSampleEntry))
			createdBox = &b.MP4VisualSampleEntries[len(b.MP4VisualSampleEntries)-1]
		}
	}

	return createdBox, nil
}

// SetHdlr passes hdlr box for later use.
func (b *Box) SetHdlr(h *hdlr.Box) {
	b.hdlr = h
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	// parse full header additional information first
	if err := b.FullHeader.ParseVersionFlag(r); err != nil {
		return err
	}

	// requires handler_type before parse
	if b.hdlr == nil {
		return fmt.Errorf("no handler_type found")
	}

	// start to parse payload
	var parsedBytes uint64

	data := make([]byte, 4)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.EntryCount = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	for i := 0; i < int(b.EntryCount); i++ {
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
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}

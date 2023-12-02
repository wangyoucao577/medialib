// Package trak represents Track Reference Box.
package trak

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/edts"
	"github.com/wangyoucao577/medialib/container/mp4/box/mdia"
	"github.com/wangyoucao577/medialib/container/mp4/box/meta"
	"github.com/wangyoucao577/medialib/container/mp4/box/tkhd"
	"github.com/wangyoucao577/medialib/container/mp4/box/uuid"
)

// Box represents a trak box.
type Box struct {
	box.FullHeader `json:"full_header"`

	Tkhd *tkhd.Box  `json:"tkhd,omitempty"`
	Mdia *mdia.Box  `json:"mdia,omitempty"`
	Uuid *uuid.Box  `json:"uuid,omitempty"`
	Edts []edts.Box `json:"edts,omitempty"`
	Meta []meta.Box `json:"meta,omitempty"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		FullHeader: box.FullHeader{
			Header: h,
		},

		boxesCreator: map[string]box.NewFunc{
			box.TypeTkhd: tkhd.New,
			box.TypeMdia: mdia.New,
			box.TypeUUID: uuid.New,
			box.TypeEdts: edts.New,
			box.TypeMeta: meta.New,
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
	case box.TypeTkhd:
		b.Tkhd = createdBox.(*tkhd.Box)
	case box.TypeMdia:
		b.Mdia = createdBox.(*mdia.Box)
	case box.TypeUUID:
		b.Uuid = createdBox.(*uuid.Box)
	case box.TypeEdts:
		b.Edts = append(b.Edts, *createdBox.(*edts.Box))
		createdBox = &b.Edts[len(b.Edts)-1] // reference to the last empty box
	case box.TypeMeta:
		b.Meta = append(b.Meta, *createdBox.(*meta.Box))
		createdBox = &b.Meta[len(b.Meta)-1] // reference to the last empty box
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

// Package moov represents Movie Box.
package moov

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/meta"
	"github.com/wangyoucao577/medialib/container/mp4/box/mvex"
	"github.com/wangyoucao577/medialib/container/mp4/box/mvhd"
	"github.com/wangyoucao577/medialib/container/mp4/box/trak"
	"github.com/wangyoucao577/medialib/container/mp4/box/udta"
)

// Box represents a mdat box.
type Box struct {
	box.Header `json:"header"`

	Mvhd *mvhd.Box  `json:"mvhd,omitempty"`
	Udta *udta.Box  `json:"udta,omitempty"`
	Trak []trak.Box `json:"trak,omitempty"`
	Mvex *mvex.Box  `json:"mvex,omitempty"`
	Meta []meta.Box `json:"meta,omitempty"`

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,

		boxesCreator: map[string]box.NewFunc{
			box.TypeMvhd: mvhd.New,
			box.TypeUdta: udta.New,
			box.TypeTrak: trak.New,
			box.TypeMvex: mvex.New,
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
	case box.TypeMvhd:
		b.Mvhd = createdBox.(*mvhd.Box)
	case box.TypeUdta:
		b.Udta = createdBox.(*udta.Box)
	case box.TypeTrak:
		b.Trak = append(b.Trak, *createdBox.(*trak.Box))
		createdBox = &b.Trak[len(b.Trak)-1] // reference to the last empty box
	case box.TypeMvex:
		b.Mvex = createdBox.(*mvex.Box)
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

	payloadSize := b.PayloadSize()

	if payloadSize == 0 {
		return fmt.Errorf("box %s is empty", b.Type)
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

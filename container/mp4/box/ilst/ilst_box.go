// Package ilst represents ilst type box.
// https://developer.apple.com/documentation/quicktime-file-format/metadata_item_list_atom
package ilst

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/data"
	"github.com/wangyoucao577/medialib/container/mp4/box/desc"
	"github.com/wangyoucao577/medialib/container/mp4/box/dottoo"
)

// Box represents a ilst box.
type Box struct {
	box.Header `json:"header"`

	EncodingTool *dottoo.Box `json:"encoding_tool,omitempty"`
	Desc         *desc.Box   `json:"desc,omitempty"`
	Data         []data.Box  `json:"unknown_type_data,omitempty"` // unspecified types

	boxesCreator map[string]box.NewFunc `json:"-"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,

		boxesCreator: map[string]box.NewFunc{
			box.TypeDottoo: dottoo.New,
			box.TypeDesc:   desc.New,
		},
	}
}

// CreateSubBox tries to create sub level box.
func (b *Box) CreateSubBox(h box.Header) (box.Box, error) {
	creator, ok := b.boxesCreator[h.Type.String()]
	if !ok {
		// glog.V(2).Infof("unknown box type %s, size %d payload %d", h.Type.String(), h.Size, h.PayloadSize())
		// return nil, box.ErrUnknownBoxType
		creator = data.New // use general data box for all unknown sub box of ilst
	}

	createdBox := creator(h)
	if createdBox == nil {
		glog.Fatalf("create box type %s failed", h.Type.String())
	}

	switch h.Type.String() {
	case box.TypeDottoo:
		b.EncodingTool = createdBox.(*dottoo.Box)
	case box.TypeDesc:
		b.Desc = createdBox.(*desc.Box)
	default:
		b.Data = append(b.Data, *createdBox.(*data.Box))
		createdBox = &b.Data[len(b.Data)-1]
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

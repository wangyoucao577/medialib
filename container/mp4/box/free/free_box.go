// Package free represents Free Space Box which may has type `free` or `skip`.
package free

import (
	"encoding/json"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// Box represents a ftyp box.
type Box struct {
	box.Header `json:"header"`

	Data []byte `json:"data"`
}

// MarshalJSON implements json.Marshaler interface.
func (b Box) MarshalJSON() ([]byte, error) {
	jsonBox := struct {
		box.Header `json:"header"`

		Data string `json:"data"`
	}{
		Header: b.Header,

		Data: string(b.Data),
	}

	return json.Marshal(jsonBox)
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		Header: h,
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	b.Data = make([]byte, b.PayloadSize())
	if err := util.ReadOrError(r, b.Data[:]); err != nil {
		return err
	}

	return nil
}

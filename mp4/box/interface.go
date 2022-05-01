package box

import (
	"errors"
	"io"
)

// errors
var (
	ErrNotImplemented = errors.New("not implemented")
	ErrUnknownBoxType = errors.New("unkown box type")
)

// NewFunc defines generic new function to create box.
type NewFunc func(Header) Box

// Box defines interfaces for boxes.
type Box interface {

	// Parse payload. It requires BasicBox(Header) has been set to the subset Box.
	ParsePayload(r io.Reader) error
}

// ParentBox defines functions if a box possible to have sub/child box.
type ParentBox interface {
	// CreateSubBox creates directly included box, such as create `mvhd` in `moov`, or create `moov` on top level.
	//   return ErrNotImplemented is the box doesn't have any sub box.
	CreateSubBox(Header) (Box, error)
}

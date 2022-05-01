package esds

import (
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// SLConfigDescriptor represents SLConfigDescriptor.
type SLConfigDescriptor struct {
	Descriptor Descriptor `json:"descriptor"`

	Predefined uint8 `json:"pre_defined"`
}

func (s *SLConfigDescriptor) parse(r io.Reader) (uint64, error) {
	var parsedBytes uint64

	// parse descriptor header
	if bytes, err := s.Descriptor.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		s.Predefined = data[0]
		parsedBytes += 1
	}

	//TODO: parse following data if exists

	return parsedBytes, nil
}

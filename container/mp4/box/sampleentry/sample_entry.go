// Package sampleentry represents sample entry.
package sampleentry

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// SampleEntry is the appropriate sample entry.
type SampleEntry struct {
	box.Header `json:"header"`

	// 6 bytes reserved here
	DataReferenceIndex uint16 `json:"data_reference_index"`
}

// ParseData parses payload data of sample entry, which requires header already exists.
// header's payload size will minus these data that has been parsed.
func (s *SampleEntry) ParseData(r io.Reader) error {
	if err := s.Header.Validate(); err != nil {
		return err
	}

	var parsedBytes uint64

	// ignore reserved 6 bytes in here
	if err := util.ReadOrError(r, make([]byte, 6)); err != nil {
		return err
	} else {
		parsedBytes += 6
	}

	data := make([]byte, 4)

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return err
	} else {
		s.DataReferenceIndex = binary.BigEndian.Uint16(data[:2])
		parsedBytes += 2
	}

	s.PayloadSizeMinus(int(parsedBytes))
	return nil
}

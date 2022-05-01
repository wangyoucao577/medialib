// Package soun represents Audio Sample Entry.
package soun

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/mp4/box/sampleentry"
	"github.com/wangyoucao577/medialib/util"
)

// AudioSampleEntry represents audio sample entry.
type AudioSampleEntry struct {
	sampleentry.SampleEntry

	// 8 bytes  reserved here
	ChannelCount uint16 `json:"channel_cound"`
	SampleSize   uint16 `json:"sample_size"`
	// 2 bytes pre_defined here
	// 2 bytes reserved here
	SampleRate uint32 `json:"sample_rate"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &AudioSampleEntry{
		SampleEntry: sampleentry.SampleEntry{
			Header: h,
		},
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (a *AudioSampleEntry) ParsePayload(r io.Reader) error {

	// parse sample entry data
	if err := a.SampleEntry.ParseData(r); err != nil {
		return err
	}

	var parsedBytes uint64
	data := make([]byte, 4)

	// ignore 8 bytes  reserved here
	if err := util.ReadOrError(r, make([]byte, 8)); err != nil {
		return err
	} else {
		parsedBytes += 8
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		a.ChannelCount = binary.BigEndian.Uint16(data[:2])
		a.SampleSize = binary.BigEndian.Uint16(data[2:])
		parsedBytes += 4
	}

	// ignore 4 bytes reserved + pre_defined here
	if err := util.ReadOrError(r, make([]byte, 4)); err != nil {
		return err
	} else {
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		a.SampleRate = (binary.BigEndian.Uint32(data) >> 16)
		parsedBytes += 4
	}

	a.Header.PayloadSizeMinus(int(parsedBytes))

	return nil
}

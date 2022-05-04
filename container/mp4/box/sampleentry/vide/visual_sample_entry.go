// Package vide represents Visual Sample Entry.
package vide

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/sampleentry"
	"github.com/wangyoucao577/medialib/util"
)

// VisualSampleEntry represents video sample entry.
type VisualSampleEntry struct {
	sampleentry.SampleEntry

	// 2 bytes  pre_defined here
	// 2 bytes  reserved here
	// 12 bytes pre_defined here
	Width           uint16 `json:"width"`
	Height          uint16 `json:"height"`
	Horizresolution uint32 `json:"horizresolution"`
	Vertresolution  uint32 `json:"vertresolution"`
	// 4 bytes reserved here
	FrameCount     uint16 `json:"frame_count"`
	Compressorname string `json:"compressorname"` // 32 bytes in file
	Depth          uint16 `json:"depth"`
	// 2 bytes pre_defined here
	// other optional boxes from derived specifications
}

// // MarshalJSON implements json.Marshaler interface.
// func (v VisualSampleEntry) MarshalJSON() ([]byte, error) {
// 	jsonEntry := struct {
// 		sampleentry.SampleEntry

// 		Width           uint16 `json:"width"`
// 		Height          uint16 `json:"height"`
// 		Horizresolution uint32 `json:"horizresolution"`
// 		Vertresolution  uint32 `json:"vertresolution"`

// 		FrameCount     uint16 `json:"frame_count"`
// 		Compressorname string `json:"compressorname"`
// 		Depth          uint16 `json:"depth"`
// 	}{
// 		SampleEntry: v.SampleEntry,

// 		Width:           v.Width,
// 		Height:          v.Height,
// 		Horizresolution: v.Horizresolution,
// 		Vertresolution:  v.Vertresolution,

// 		FrameCount:     v.FrameCount,
// 		Compressorname: string(v.Compressorname[:]),
// 		Depth:          v.Depth,
// 	}

// 	return json.Marshal(jsonEntry)
// }

// New creates a new Box.
func New(h box.Header) box.Box {
	return &VisualSampleEntry{
		SampleEntry: sampleentry.SampleEntry{
			Header: h,
		},
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (v *VisualSampleEntry) ParsePayload(r io.Reader) error {

	// parse sample entry data
	if err := v.SampleEntry.ParseData(r); err != nil {
		return err
	}

	var parsedBytes uint
	data := make([]byte, 4)

	// ignore reserved and pre_defined 16 bytes in here
	if err := util.ReadOrError(r, make([]byte, 16)); err != nil {
		return err
	} else {
		parsedBytes += 16
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		v.Width = binary.BigEndian.Uint16(data[:2])
		v.Height = binary.BigEndian.Uint16(data[2:])
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		v.Horizresolution = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		v.Vertresolution = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	// ignore reserved 4 bytes in here
	if err := util.ReadOrError(r, make([]byte, 4)); err != nil {
		return err
	} else {
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return err
	} else {
		v.FrameCount = binary.BigEndian.Uint16(data[:2])
		parsedBytes += 2
	}

	const compressorNameDataLen = 32 // compressor name has fixed 32 bytes in file
	compressorNameData := make([]byte, compressorNameDataLen)
	if err := util.ReadOrError(r, compressorNameData); err != nil {
		return err
	} else {
		compressorNameData = bytes.TrimRight(compressorNameData, "\u0000") // trim last 0 to avoid `\u0000` in encoded json
		v.Compressorname = string(compressorNameData)
		parsedBytes += uint(compressorNameDataLen)
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return err
	} else {
		v.Depth = binary.BigEndian.Uint16(data[:2])
		parsedBytes += 2
	}

	// ignore pre_defined 2 bytes in here
	if err := util.ReadOrError(r, make([]byte, 2)); err != nil {
		return err
	} else {
		parsedBytes += 2
	}

	v.Header.PayloadSizeMinus(int(parsedBytes))

	return nil
}

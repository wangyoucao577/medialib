// Package es represents MPEG-4 AVC Elementary Stream.
// It contains "Video elementary stream only" which also named "mp4 es".
// The structure was defined in ISO/IEC-14496-15 5.2.3.
package es

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// LengthNALU represents a length and nalu composition.
type LengthNALU struct {
	Length uint32       `json:"length"`
	NALU   nalu.NALUnit `json:"nalu"`
}

// ElementaryStream represents AVC Elementary Stream.
type ElementaryStream struct {
	LengthNALU []LengthNALU `json:"length_nalu"`

	LengthSize uint32 `json:"length_size"`

	// cache for slice parsing
	sps *sps.SequenceParameterSetData `json:"-"`
	pps *pps.PictureParameterSet      `json:"-"`
}

// SetLengthSize sets length size before every nalu.
// It's mandantory that should be set before `Parse``.
func (e *ElementaryStream) SetLengthSize(l uint32) {
	e.LengthSize = l
}

// Parse parses bytes to AVC Elementary Stream, return parsed bytes or error.
// It's allowed to call multiple times since data maybe splitted in storage.
func (e *ElementaryStream) Parse(r io.Reader, size int) (uint64, error) {
	if e.LengthSize == 0 {
		return 0, fmt.Errorf("length size not set")
	}

	var parsedBytes uint64
	for parsedBytes < uint64(size) {
		ln := LengthNALU{
			NALU: nalu.NALUnit{SequenceParameterSetData: e.sps, PictureParameterSet: e.pps},
		}

		// parse nalu length
		data := make([]byte, 4)
		if err := util.ReadOrError(r, data[4-e.LengthSize:]); err != nil {
			return parsedBytes, err
		} else {
			if e.LengthSize == 4 || e.LengthSize == 3 {
				ln.Length = binary.BigEndian.Uint32(data)
			} else if e.LengthSize == 2 {
				ln.Length = uint32(binary.BigEndian.Uint16(data))
			} else if e.LengthSize == 1 {
				ln.Length = uint32(data[3])
			} else {
				return parsedBytes, fmt.Errorf("invalid length size: %d", e.LengthSize)
			}
			parsedBytes += uint64(e.LengthSize)
		}

		if bytes, err := ln.NALU.Parse(r, int(ln.Length)); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += bytes
		}
		if ln.NALU.NALUnitType == nalu.TypeSPS {
			e.sps = ln.NALU.SequenceParameterSetData
		}
		if ln.NALU.NALUnitType == nalu.TypePPS {
			e.pps = ln.NALU.PictureParameterSet
		}

		e.LengthNALU = append(e.LengthNALU, ln)
	}

	return parsedBytes, nil
}

// JSON marshals elementary stream to JSON representation
func (e *ElementaryStream) JSON() ([]byte, error) {
	return json.Marshal(e)
}

// JSONIndent marshals elementary stream to JSON representation with customized indent.
func (e *ElementaryStream) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(e, prefix, indent)
}

// YAML formats elementary stream to YAML representation.
func (e *ElementaryStream) YAML() ([]byte, error) {
	j, err := json.Marshal(e)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CSV formats boxes to CSV representation, which isn't supported at the moment.
func (e *ElementaryStream) CSV() ([]byte, error) {
	return nil, fmt.Errorf("csv representation does not support yet")
}

// Dump dumps raw data into io.Writer.
func (e *ElementaryStream) Dump(w io.Writer) (int, error) {
	if e.LengthSize == 0 || e.LengthSize > 4 {
		return 0, fmt.Errorf("invalid elementary stream")
	}

	var writedBytes int

	for i := range e.LengthNALU {
		// data := []byte{0x00, 0x00, 0x00, 0x01}
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data, e.LengthNALU[i].Length)
		data = data[4-e.LengthSize:]
		if n, err := w.Write(data); err != nil {
			return writedBytes, err
		} else if n != len(data) {
			return writedBytes, fmt.Errorf("write bytes unmatch, expect(%d) != actual(%d)", len(data), n)
		} else {
			writedBytes += n
		}

		rsbp := e.LengthNALU[i].NALU.Raw()
		if n, err := w.Write(rsbp); err != nil {
			return writedBytes, err
		} else if n != len(rsbp) {
			return writedBytes, fmt.Errorf("write bytes unmatch, expect(%d) != actual(%d)", len(data), n)
		} else {
			writedBytes += n
		}
	}

	return writedBytes, nil
}

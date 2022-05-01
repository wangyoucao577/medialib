// Package annexbes represents Annex B defined AVC Elementary byte stream,
// which was defined in ISO/IEC-14496-19 Annex B.
package annexbes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

var (
	startCode3Bytes = []byte{0x00, 0x00, 0x01}
	startCode4Bytes = []byte{0x00, 0x00, 0x00, 0x01}
)

// ElementaryStream represents AVC Elementary Stream.
type ElementaryStream struct {
	NALU []nalu.NALUnit `json:"nalu"`
}

// Parse parses bytes to AVC AnnexB format Elementary Stream, return parsed bytes or error.
// The size could be 0 that indicates parse until nothing to read, otherwise read max size.
// It's allowed to call multiple times since data maybe splitted in storage.
func (e *ElementaryStream) Parse(r io.Reader, size int) (uint64, error) {

	var parsedBytes uint64
	startCodeData := []byte{}
	naluData := []byte{}
	var sps *sps.SequenceParameterSetData
	var pps *pps.PictureParameterSet

	for {
		if size > 0 && parsedBytes >= uint64(size) { // valid size
			break
		}

		if nextByte, err := util.ReadByteOrError(r); err != nil {
			if err == io.EOF {
				break // break to parse last nalu
			}
			return parsedBytes, err
		} else {
			startCodeData = append(startCodeData, nextByte)
			naluData = append(naluData, nextByte)
			parsedBytes += 1
		}

		if len(startCodeData) < 3 { // data not enough to parse
			continue
		}
		if len(startCodeData) == 3 &&
			!bytes.Equal(startCodeData, startCode3Bytes) {
			continue // read one more byte
		}
		if len(startCodeData) == 4 &&
			!bytes.Equal(startCodeData[:3], startCode3Bytes) &&
			!bytes.Equal(startCodeData, startCode4Bytes) {
			startCodeData = startCodeData[1:] // ignore 1 byte
			continue
		}

		var threeMatchedInFour bool
		if len(startCodeData) == 4 && bytes.Equal(startCodeData[:3], startCode3Bytes) {
			threeMatchedInFour = true // [0x0, 0x0, 0x1, xx] matched, need special handle
		}

		naluData = bytes.TrimSuffix(naluData, startCodeData) // clear start code

		// parse NALU here
		if len(naluData) > 0 {
			n := nalu.NALUnit{SequenceParameterSetData: sps, PictureParameterSet: pps}
			if _, err := n.Parse(bytes.NewReader(naluData), len(naluData)); err != nil {
				return parsedBytes, err
			}
			if n.NALUnitType == nalu.TypeSPS {
				sps = n.SequenceParameterSetData
			}
			if n.NALUnitType == nalu.TypePPS {
				pps = n.PictureParameterSet
			}
			e.NALU = append(e.NALU, n)
		}

		naluData = naluData[:0] // clear
		if threeMatchedInFour {
			startCodeData = startCodeData[3:]             // clear to ignore start code 3 bytes, but keep the last one
			naluData = append(naluData, startCodeData...) // also keep the last byte
		} else {
			startCodeData = startCodeData[:0] // clear to ignore start code
		}
	}

	// parse last NALU
	if len(naluData) > 0 {
		n := nalu.NALUnit{SequenceParameterSetData: sps, PictureParameterSet: pps}
		if _, err := n.Parse(bytes.NewReader(naluData), len(naluData)); err != nil {
			return parsedBytes, err
		}
		e.NALU = append(e.NALU, n)
	}

	if size > 0 && parsedBytes != uint64(size) {
		glog.Warningf("expect parse %d bytes but actually parsed %d bytes", size, parsedBytes)
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
	if len(e.NALU) == 0 {
		return 0, fmt.Errorf("empty elementary stream")
	}

	var writedBytes int

	for i := range e.NALU {
		data := startCode4Bytes // Annex B start code
		if n, err := w.Write(data); err != nil {
			return writedBytes, err
		} else if n != len(data) {
			return writedBytes, fmt.Errorf("write bytes unmatch, expect(%d) != actual(%d)", len(data), n)
		} else {
			writedBytes += n
		}

		rsbp := e.NALU[i].Raw()
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

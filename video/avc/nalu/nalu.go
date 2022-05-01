// Package nalu represents AVC NAL(Network Abstract Layer) Units.
package nalu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/nalu/aud"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sei"
	"github.com/wangyoucao577/medialib/video/avc/nalu/slice"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// NALUnit represents AVC NAL Unit that defined in ISO/IEC-14496-10 7.3.1.
type NALUnit struct {
	RawBytes []byte `json:"-"` // store raw bytes

	ForbiddenZeroBit uint8 `json:"forbidden_zero_bit"` // 1 bit, shoule be 0 always
	NALRefIdc        uint8 `json:"nal_ref_idc"`        // 2 bits
	NALUnitType      uint8 `json:"nal_unit_type"`      // 5 bits

	nalUnitHeaderSvcExtension []byte `json:"-"` // TODO: parse nal_unit_header_svc_extension

	RBSP []byte `json:"-"` // Raw byte sequence payloads

	// parsed RBRP if available
	SEIMessage               *sei.SEIMessage                      `json:"sei_message,omitempty"`
	AccessUnitDelimiter      *aud.AccessUnitDelimiter             `json:"access_unit_delimiter,omitempty"`
	SequenceParameterSetData *sps.SequenceParameterSetData        `json:"seq_parameter_set_data,omitempty"`
	PictureParameterSet      *pps.PictureParameterSet             `json:"picture_parameter_set,omitempty"`
	IDR                      []slice.LayerWithoutPartitioningRbsp `json:"idr,omitempty"`
	NonIDR                   []slice.LayerWithoutPartitioningRbsp `json:"non-idr,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (n *NALUnit) MarshalJSON() ([]byte, error) {
	var nj = struct {
		RawBytes []byte `json:"raw_bytes,omitempty"`

		ForbiddenZeroBit       uint8  `json:"forbidden_zero_bit"` // 1 bit, shoule be 0 always
		NALRefIdc              uint8  `json:"nal_ref_idc"`        // 2 bits
		NALUnitType            uint8  `json:"nal_unit_type"`      // 5 bits
		NALUnitTypeDescription string `json:"nal_unit_type_description"`

		nalUnitHeaderSvcExtension []byte `json:"-"` // TODO: parse nal_unit_header_svc_extension

		// raw bytes and raw bytes sequence payloads
		RBSP []byte `json:"rbsp,omitempty"` // Raw byte sequence payloads

		// parsed RBRP data
		SEIMessage               *sei.SEIMessage                      `json:"sei_message,omitempty"`
		AccessUnitDelimiter      *aud.AccessUnitDelimiter             `json:"access_unit_delimiter,omitempty"`
		SequenceParameterSetData *sps.SequenceParameterSetData        `json:"seq_parameter_set_data,omitempty"`
		PictureParameterSet      *pps.PictureParameterSet             `json:"picture_parameter_set,omitempty"`
		IDR                      []slice.LayerWithoutPartitioningRbsp `json:"idr,omitempty"`
		NonIDR                   []slice.LayerWithoutPartitioningRbsp `json:"non-idr,omitempty"`
	}{
		// RawBytes:               n.RawBytes, // set by type

		ForbiddenZeroBit:       n.ForbiddenZeroBit,
		NALRefIdc:              n.NALRefIdc,
		NALUnitType:            n.NALUnitType,
		NALUnitTypeDescription: TypeDescription(int(n.NALUnitType)),

		nalUnitHeaderSvcExtension: n.nalUnitHeaderSvcExtension,

		// RBSP: b.RBSP, // set by type

		SEIMessage:               n.SEIMessage,
		AccessUnitDelimiter:      n.AccessUnitDelimiter,
		SequenceParameterSetData: n.SequenceParameterSetData,
		PictureParameterSet:      n.PictureParameterSet,
		IDR:                      n.IDR,
		NonIDR:                   n.NonIDR,
	}

	switch n.NALUnitType {
	case TypeSEI:
		fallthrough
	case TypeAccessUnitDelimiter:
		fallthrough
	case TypeSPS:
		fallthrough
	case TypePPS:
		nj.RawBytes = n.RawBytes
		nj.RBSP = n.RBSP
	}

	switch n.NALUnitType {
	case TypeSPS:
	case TypePPS: // clear sps data if PPS type, since they're set for data parsing
		nj.SequenceParameterSetData = nil
	default: // clear sps and pps data if NOT SPS/PPS type, since they're set for data parsing
		nj.SequenceParameterSetData = nil
		nj.PictureParameterSet = nil
	}

	return json.Marshal(nj)
}

// Parse parses bytes to AVC NAL Unit, return parsed bytes or error.
// The NAL Unit syntax defined in ISO/IEC-14496-10 7.3.1.
func (n *NALUnit) Parse(r io.Reader, size int) (uint64, error) {
	var parsedBytes uint64

	// parse nalu length
	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		n.RawBytes = append(n.RawBytes, data...)
		n.ForbiddenZeroBit = (data[0] >> 7) & 0x1
		n.NALRefIdc = (data[0] >> 5) & 0x3
		n.NALUnitType = data[0] & 0x1F
		parsedBytes += 1
	}

	if n.ForbiddenZeroBit != 0 {
		return parsedBytes, fmt.Errorf("nalu forbidden_zero_bit should be 0")
	}
	if !IsValidNALUType(int(n.NALUnitType)) {
		return parsedBytes, fmt.Errorf("unknown nal_unit_type %d", n.NALUnitType)
	}

	nalUnitHeaderBytes := 1

	if n.NALUnitType == TypePrefix || n.NALUnitType == TypeSliceExtersion {
		glog.Warningf("nalu type %d svc_extension_flag and nal_unit_header_svc_extension parsing TODO", n.NALUnitType)
		//TODO: parse payload
		n.nalUnitHeaderSvcExtension = make([]byte, 2)
		if err := util.ReadOrError(r, n.nalUnitHeaderSvcExtension); err != nil {
			return parsedBytes, err
		} else {
			n.RawBytes = append(n.RawBytes, n.nalUnitHeaderSvcExtension...)
			parsedBytes += 2
			nalUnitHeaderBytes += 2
		}
	}

	n.RBSP = make([]byte, size-nalUnitHeaderBytes)
	if err := util.ReadOrError(r, n.RBSP); err != nil {
		return parsedBytes, err
	} else {
		n.RawBytes = append(n.RawBytes, n.RBSP...)
		parsedBytes += uint64(size - nalUnitHeaderBytes)
	}
	n.RBSP = getRBSP(n.RBSP)

	// Parse RBSP
	parser := n.prepareRBRPParser()
	if parser != nil {
		if _, err := parser.Parse(bytes.NewReader(n.RBSP), len(n.RBSP)); err != nil {
			if err != slice.ErrEmptyParameterSet {
				return parsedBytes, fmt.Errorf("parse nalu type %d rbrp failed, err %v", n.NALUnitType, err)
			} else {
				glog.Warningf("parse nalu type %d rbrp failed, ignore it, err %v", n.NALUnitType, err)
			}
		}
	} else {
		glog.Warningf("unknown nalu type %d, ignored", n.NALUnitType)
	}

	return parsedBytes, nil
}

func (n *NALUnit) prepareRBRPParser() NALUParser {
	switch n.NALUnitType {
	case TypeSEI:
		n.SEIMessage = &sei.SEIMessage{}
		return n.SEIMessage
	case TypeAccessUnitDelimiter:
		n.AccessUnitDelimiter = &aud.AccessUnitDelimiter{}
		return n.AccessUnitDelimiter
	case TypeSPS:
		n.SequenceParameterSetData = &sps.SequenceParameterSetData{}
		return n.SequenceParameterSetData
	case TypePPS:
		n.PictureParameterSet = &pps.PictureParameterSet{}
		n.PictureParameterSet.SetSPS(n.SequenceParameterSetData)
		return n.PictureParameterSet
	case TypeIDR:
		n.IDR = append(n.IDR, slice.LayerWithoutPartitioningRbsp{})
		newSlice := &n.IDR[len(n.IDR)-1]
		newSlice.SetSequenceHeaders(n.SequenceParameterSetData, n.PictureParameterSet)
		return newSlice
	case TypeNonIDR:
		n.NonIDR = append(n.NonIDR, slice.LayerWithoutPartitioningRbsp{})
		newSlice := &n.NonIDR[len(n.NonIDR)-1]
		newSlice.SetSequenceHeaders(n.SequenceParameterSetData, n.PictureParameterSet)
		return newSlice

		// TODO: others
	}
	return nil
}

// Raw translates to raw bytes data.
func (n *NALUnit) Raw() []byte {
	return n.RawBytes
}

// raw RBSP -> RBSP, remove emulation_prevention_three_byte 0x03
func getRBSP(rbspBytes []byte) []byte {
	numBytesOfRBSP := len(rbspBytes)

	rbsp := []byte{}
	for i := 0; i < numBytesOfRBSP; i++ {
		if i+2 < numBytesOfRBSP &&
			rbspBytes[i] == 0x00 &&
			rbspBytes[i+1] == 0x00 &&
			rbspBytes[i+2] == 0x03 {
			rbsp = append(rbsp, rbspBytes[i], rbspBytes[i+1])
			i += 2
			// ignore emulation_prevention_three_byte, equal to 0x03
		} else {
			rbsp = append(rbsp, rbspBytes[i])
		}
	}
	return rbsp
}

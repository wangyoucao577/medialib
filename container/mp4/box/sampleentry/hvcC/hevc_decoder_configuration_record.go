package hvcc

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// LengthNALU represents Length and NALU composition.
type LengthNALU struct {
	NALUnitLength uint16 `json:"nalUnitLength"`
	NALUnit       []byte `json:"nalUnit,omitempty"`
}

// Array represents array in HEVC Decoder configuration record.
type Array struct {
	ArrayCompleteness uint8 `json:"array_completeness"` // 1 bit
	// 1 bit reserved
	NALUnitType uint8        `json:"NAL_unit_type"` // 6 bits
	NumNalus    uint16       `json:"numNalus"`
	LengthNALUs []LengthNALU `json:"length_nalu,omitempty"`
}

// HEVCDecoderConfigurationRecord defines HEVC Decoder configuration record.
type HEVCDecoderConfigurationRecord struct {
	ConfigurationVersion             uint8  `json:"configuration_version"`
	GeneralProfileSpace              uint8  `json:"general_profile_space"` // 2 bits
	GeneralTierFlag                  uint8  `json:"general_tier_flag"`     // 1 bit
	GeneralProfileIdc                uint8  `json:"general_profile_idc"`   // 5 bits
	GeneralProfileCompatibilityFlags uint32 `json:"general_profile_compatibility_flags"`
	GeneralConstraintIndicatorFlags  uint64 `json:"general_constraint_indicator_flags"` // 48 bits
	GeneralLevelIdc                  uint8  `json:"general_level_idc"`
	// 4 bits reserved
	MinSpatialSegmentationIdc uint16 `json:"min_spatial_segmentation_idc"` // 12 bits
	// 6 bits reserved
	ParallelismType uint8 `json:"parallelismType"` // 2 bits
	// 6 bits reserved
	ChromaFormatIdc uint8 `json:"chroma_format_idc"` // 2 bits
	// 5 bits reserved
	BitDepthLumaMinus8 uint8 `json:"bit_depth_luma_minus8"` // 3 bits
	// 5 bits reserved
	BitDepthChromaMinus8 uint8   `json:"bit_depth_chroma_minus8"` // 3 bits
	AvgFrameRate         uint16  `json:"avgFrameRate"`
	ConstantFrameRate    uint8   `json:"constantFrameRate"`  // 2 bits
	NumTemporalLayers    uint8   `json:"numTemporalLayers"`  // 3 bits
	TemporalIdNested     uint8   `json:"temporalIdNested"`   // 1 bit
	LengthSizeMinusOne   uint8   `json:"lengthSizeMinusOne"` // 2 bits
	NumOfArrays          uint8   `json:"numOfArrays"`
	Arrays               []Array `json:"arrays,omitempty"`
}

// Parse parses HEVCDecoderConfigurationRecord.
func (h *HEVCDecoderConfigurationRecord) Parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	data := make([]byte, 8)

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return parsedBytes, err
	} else {
		h.ConfigurationVersion = data[0]
		h.GeneralProfileSpace = (data[1] >> 6) & 0x3
		h.GeneralTierFlag = (data[1] >> 5) & 1
		h.GeneralProfileSpace = data[1] & 0x1F

		parsedBytes += 2
	}

	if err := util.ReadOrError(r, data[:4]); err != nil {
		return parsedBytes, err
	} else {
		h.GeneralProfileCompatibilityFlags = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data[2:]); err != nil {
		return parsedBytes, err
	} else {
		data[0] = 0x0
		data[1] = 0x0
		h.GeneralConstraintIndicatorFlags = binary.BigEndian.Uint64(data)
		parsedBytes += 6
	}

	if err := util.ReadOrError(r, data[:7]); err != nil {
		return parsedBytes, err
	} else {
		h.GeneralLevelIdc = data[0]
		data[1] &= 0xF
		h.MinSpatialSegmentationIdc = binary.BigEndian.Uint16(data[1:3])
		h.ParallelismType = data[3] & 0x3
		h.ChromaFormatIdc = data[4] & 0x3
		h.BitDepthLumaMinus8 = data[5] & 0x7
		h.BitDepthChromaMinus8 = data[6] & 0x7
		parsedBytes += 7
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return parsedBytes, err
	} else {
		h.AvgFrameRate = binary.BigEndian.Uint16(data[:2])
		parsedBytes += 2
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return parsedBytes, err
	} else {
		h.ConstantFrameRate = (data[0] >> 6) & 0x3
		h.NumTemporalLayers = (data[0] >> 3) & 0x7
		h.TemporalIdNested = (data[0] >> 2) & 0x1
		h.LengthSizeMinusOne = data[0] & 0x3
		h.NumOfArrays = data[1]
		parsedBytes += 2
	}

	if h.NumOfArrays == 0 {
		return parsedBytes, nil
	}

	h.Arrays = make([]Array, h.NumOfArrays)
	for i := 0; i < int(h.NumOfArrays); i++ {
		if err := util.ReadOrError(r, data[:3]); err != nil {
			return parsedBytes, err
		} else {
			h.Arrays[i].ArrayCompleteness = (data[0] >> 7) & 0x1
			h.Arrays[i].NALUnitType = data[0] & 0x3F
			h.Arrays[i].NumNalus = binary.BigEndian.Uint16(data[1:3])
			parsedBytes += 3
		}

		for j := 0; j < int(h.Arrays[i].NumNalus); j++ {
			lenNALU := LengthNALU{}
			if err := util.ReadOrError(r, data[:2]); err != nil {
				return parsedBytes, err
			} else {
				lenNALU.NALUnitLength = binary.BigEndian.Uint16(data[:2])
				parsedBytes += 2
			}

			if lenNALU.NALUnitLength == 0 {
				continue
			}

			lenNALU.NALUnit = make([]byte, lenNALU.NALUnitLength)
			if err := util.ReadOrError(r, lenNALU.NALUnit); err != nil {
				return parsedBytes, err
			} else {
				parsedBytes += uint64(lenNALU.NALUnitLength)
			}

			h.Arrays[i].LengthNALUs = append(h.Arrays[i].LengthNALUs, lenNALU)
		}

	}

	return parsedBytes, nil
}

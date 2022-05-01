package sps

import (
	"encoding/binary"

	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
)

// VUIParameters represents vui_parameters defined in ISO/IEC-14496-10 Annex E.1.1
type VUIParameters struct {
	AspectRatioInfoPresentFlag         uint8                     `json:"aspect_ratio_info_present_flag"` // 1 bit
	AspectRatioIdc                     *uint8                    `json:"aspect_ratio_idc,omitempty"`     // 8 bits
	SarWidth                           *uint16                   `json:"sar_width,omitempty"`
	SarHeight                          *uint16                   `json:"sar_height,omitempty"`
	OverscanInfoPresentFlag            uint8                     `json:"overscan_info_present_flag"`                // 1 bit
	OverscanAppropriateFlag            *uint8                    `json:"overscan_appropriate_flag,omitempty"`       // 1 bit
	VideoSignalTypePresentFlag         uint8                     `json:"video_signal_type_present_flag"`            // 1 bit
	VideoFormat                        *uint8                    `json:"video_format,omitempty"`                    // 3 bits
	VideoFullRangeFlag                 *uint8                    `json:"video_full_range_flag,omitempty"`           // 1 bit
	ColourDescriptionPresentFlag       *uint8                    `json:"colour_description_present_flag,omitempty"` // 1 bit
	ColourPrimaries                    *uint8                    `json:"colour_primaries,omitempty"`
	TransferCharacteristics            *uint8                    `json:"transfer_characteristics,omitempty"`
	MatrixCoefficients                 *uint8                    `json:"matrix_coefficients,omitempty"`
	ChromaLocInfoPresentFlag           uint8                     `json:"chroma_loc_info_present_flag"` // 1 bit
	ChromaSampleLocTypeTopField        *expgolombcoding.Unsigned `json:"chroma_sample_loc_type_top_field,omitempty"`
	ChromaSampleLocTypeBottomField     *expgolombcoding.Unsigned `json:"chroma_sample_loc_type_bottom_field,omitempty"`
	TimingInfoPresentFlag              uint8                     `json:"timing_info_present_flag"`
	NumUnitsInTick                     *uint32                   `json:"num_units_in_tick,omitempty"`
	TimeScale                          *uint32                   `json:"time_scale,omitempty"`
	FixedFrameRateFlag                 *uint8                    `json:"fixed_frame_rate_flag,omitempty"` // 1 bit
	NalHrdParametersPresentFlag        uint8                     `json:"nal_hrd_parameters_present_flag"` // 1 bit
	NalHrdParmaeters                   *HrdParameters            `json:"nal_hrd_parameters,omitempty"`
	VclHrdParametersPresentFlag        uint8                     `json:"vcl_hrd_parameters_present_flag"` // 1 bit
	VclHrdParmaeters                   *HrdParameters            `json:"vcl_hrd_parameters,omitempty"`
	LowDelayHrdFlag                    *uint8                    `json:"low_delay_hrd_flag,omitempty"`                      // 1 bit
	PicStructPresentFlag               uint8                     `json:"pic_struct_present_flag"`                           // 1 bit
	BitstreamRestrictionFlag           uint8                     `json:"bitstream_restriction_flag"`                        // 1 bit
	MotionVectorsOverPicBoundariesFlag *uint8                    `json:"motion_vectors_over_pic_boundaries_flag,omitempty"` // 1 bit
	MaxBytesPerPicDenom                *expgolombcoding.Unsigned `json:"max_bytes_per_pic_denom,omitempty"`
	MaxBitsPerMbDenom                  *expgolombcoding.Unsigned `json:"max_bits_per_mb_denom,omitempty"`
	Log2MaxMvLengthHorizontal          *expgolombcoding.Unsigned `json:"log2_max_mv_length_horizontal,omitempty"`
	Log2MaxMvLengthVertical            *expgolombcoding.Unsigned `json:"log2_max_mv_length_vertical,omitempty"`
	MaxNumReorderFrames                *expgolombcoding.Unsigned `json:"max_num_reorder_frames,omitempty"`
	MaxDecFrameFuffering               *expgolombcoding.Unsigned `json:"max_dec_frame_buffering,omitempty"`
}

// return parsed bits
func (v *VUIParameters) parse(br *bitreader.Reader) (uint64, error) {
	var parsedBits uint64

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.AspectRatioInfoPresentFlag = nextBit
		parsedBits++
	}

	if v.AspectRatioInfoPresentFlag != 0 {
		if nextByte, err := br.ReadByte(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.AspectRatioIdc = &nextByte
			parsedBits += bitsPerByte
		}

		if *v.AspectRatioIdc == 255 { // Extented_SAR = 255
			if nextBytes, err := br.ReadBytes(4); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				width := binary.BigEndian.Uint16(nextBytes[:2])
				v.SarWidth = &width
				height := binary.BigEndian.Uint16(nextBytes[2:])
				v.SarHeight = &height
				parsedBits += 4 * bitsPerByte
			}
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.OverscanInfoPresentFlag = nextBit
		parsedBits++
	}

	if v.OverscanInfoPresentFlag != 0 {
		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.OverscanAppropriateFlag = &nextBit
			parsedBits++
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.VideoSignalTypePresentFlag = nextBit
		parsedBits++
	}

	if v.VideoSignalTypePresentFlag != 0 {
		if nextBits, err := br.ReadBits(3); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.VideoFormat = &nextBits[0]
			parsedBits += 3
		}

		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.VideoFullRangeFlag = &nextBit
			parsedBits++
		}

		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.ColourDescriptionPresentFlag = &nextBit
			parsedBits++
		}

		if *v.ColourDescriptionPresentFlag != 0 {
			if nextBytes, err := br.ReadBytes(3); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				v.ColourPrimaries = &nextBytes[0]
				v.TransferCharacteristics = &nextBytes[1]
				v.MatrixCoefficients = &nextBytes[2]
				parsedBits += 3 * bitsPerByte
			}
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.ChromaLocInfoPresentFlag = nextBit
		parsedBits++
	}

	if v.ChromaLocInfoPresentFlag != 0 {
		expUnsigned := &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.ChromaSampleLocTypeTopField = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.ChromaSampleLocTypeBottomField = expUnsigned
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.TimingInfoPresentFlag = nextBit
		parsedBits++
	}

	if v.TimingInfoPresentFlag != 0 {
		if nextBytes, err := br.ReadBytes(8); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			numUnitsInTick := binary.BigEndian.Uint32(nextBytes[:4])
			timeScale := binary.BigEndian.Uint32(nextBytes[4:])
			v.NumUnitsInTick = &numUnitsInTick
			v.TimeScale = &timeScale
			parsedBits += 8 * bitsPerByte
		}

		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.FixedFrameRateFlag = &nextBit
			parsedBits++
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.NalHrdParametersPresentFlag = nextBit
		parsedBits++
	}

	if v.NalHrdParametersPresentFlag != 0 {
		v.NalHrdParmaeters = &HrdParameters{}
		if costBits, err := v.NalHrdParmaeters.parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.VclHrdParametersPresentFlag = nextBit
		parsedBits++
	}

	if v.VclHrdParametersPresentFlag != 0 {
		v.VclHrdParmaeters = &HrdParameters{}
		if costBits, err := v.VclHrdParmaeters.parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
	}

	if v.NalHrdParametersPresentFlag != 0 || v.VclHrdParametersPresentFlag != 0 {
		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.LowDelayHrdFlag = &nextBit
			parsedBits++
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.PicStructPresentFlag = nextBit
		parsedBits++
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		v.BitstreamRestrictionFlag = nextBit
		parsedBits++
	}

	if v.BitstreamRestrictionFlag != 0 {
		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			v.MotionVectorsOverPicBoundariesFlag = &nextBit
			parsedBits++
		}

		expUnsigned := &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.MaxBytesPerPicDenom = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.MaxBitsPerMbDenom = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.Log2MaxMvLengthHorizontal = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.Log2MaxMvLengthVertical = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.MaxNumReorderFrames = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		v.MaxDecFrameFuffering = expUnsigned
	}

	return parsedBits, nil
}

// Package sps defined AVC Sequence Parameter Sets information.
package sps

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
)

// SequenceParameterSetData represents SequenceParameterSetData defined in ISO/IEC-14496-10 7.3.2.
type SequenceParameterSetData struct {
	ProfileIdc         uint8  `json:"profile_idc"`
	ProfileIdcName     string `json:"profile_idc_name"`     // NOT in byte stream, only store for better intuitive
	ConstraintSet0Flag uint8  `json:"constraint_set0_flag"` // 1 bit
	ConstraintSet1Flag uint8  `json:"constraint_set1_flag"` // 1 bit
	ConstraintSet2Flag uint8  `json:"constraint_set2_flag"` // 1 bit
	ConstraintSet3Flag uint8  `json:"constraint_set3_flag"` // 1 bit
	ConstraintSet4Flag uint8  `json:"constraint_set4_flag"` // 1 bit
	ConstraintSet5Flag uint8  `json:"constraint_set5_flag"` // 1 bit
	// 2 bytes reserved here
	LevelIdc                        uint8                      `json:"level_idc"`
	SeqParameterSetID               expgolombcoding.Unsigned   `json:"seq_parameter_set_id"`                           // Exp-Golomb-coded
	ChromaFormatIdc                 *expgolombcoding.Unsigned  `json:"chroma_format_idc,omitempty"`                    // Exp-Golomb-coded
	SeparateColourPlaneFlag         *uint8                     `json:"separate_colour_plane_flag,omitempty"`           // 1 bit
	BitDepthLumaMinus8              *expgolombcoding.Unsigned  `json:"bit_depth_luma_minus8,omitempty"`                // Exp-Golomb-coded
	BitDepthChromaMinus8            *expgolombcoding.Unsigned  `json:"bit_depth_chroma_minus8,omitempty"`              // Exp-Golomb-coded
	QpprimeYZeroTransformBypassFlag *uint8                     `json:"qpprime_y_zero_transform_bypass_flag,omitempty"` // 1 bit
	SeqScalingMatrixPresentFlag     *uint8                     `json:"seq_scaling_matrix_present_flag,omitempty"`      // 1 bit
	SeqScalingListPresentFlag       []uint8                    `json:"seq_scaling_list_present_flag,omitempty"`        // 1 bit per flag
	DeltaScale                      [][]expgolombcoding.Signed `json:"delta_scale,omitempty"`
	ScalingList4x4                  [][]int                    `json:"scaling_list_4x4,omitempty"`
	ScalingList8x8                  [][]int                    `json:"scaling_list_8x8,omitempty"`
	Log2MaxFrameNumMinus4           expgolombcoding.Unsigned   `json:"log2_max_frame_num_minus4"`                       // Exp-Golomb-coded
	PicOrderCntType                 expgolombcoding.Unsigned   `json:"pic_order_cnt_type"`                              // Exp-Golomb-coded
	Log2MaxPicOrderCntLsbMinus4     *expgolombcoding.Unsigned  `json:"log2_max_pic_order_cnt_lsb_minus4,omitempty"`     // Exp-Golomb-coded
	DeltaPicOrderAlwaysZeroFlag     *uint8                     `json:"delta_pic_order_always_zero_flag,omitempty"`      // 1 bit
	OffsetForNonRefPic              *expgolombcoding.Signed    `json:"offset_for_non_ref_pic,omitempty"`                // Exp-Golomb-coded
	OffsetForTopToBottomField       *expgolombcoding.Signed    `json:"offset_for_top_to_bottom_field,omitempty"`        // Exp-Golomb-coded
	NumRefFramesInPicOrderCntCycle  *expgolombcoding.Unsigned  `json:"num_ref_frames_in_pic_order_cnt_cycle,omitempty"` // Exp-Golomb-coded
	OffsetForRefFrame               []expgolombcoding.Signed   `json:"offset_for_ref_frame,omitempty"`                  // Exp-Golomb-coded
	MaxNumRefFrames                 expgolombcoding.Unsigned   `json:"max_num_ref_frames"`                              // Exp-Golomb-coded
	GapsInFrameNumValueAllowedFlag  uint8                      `json:"gaps_in_frame_num_value_allowed_flag"`            // 1 bit
	PicWidthInMbsMinus1             expgolombcoding.Unsigned   `json:"pic_width_in_mbs_minus1"`                         // Exp-Golomb-coded
	PicHeightInMapUnitsMinus1       expgolombcoding.Unsigned   `json:"pic_height_in_map_units_minus1"`                  // Exp-Golomb-coded
	FrameMbsOnlyFlag                uint8                      `json:"frame_mbs_only_flag"`                             // 1 bit
	MbAdaptiveFrameFieldFlag        *uint8                     `json:"mb_adaptive_frame_field_flag,omitempty"`          // 1 bit
	Direct8x8InferenceFlag          uint8                      `json:"direct_8x8_inference_flag"`                       // 1 bit
	FrameCroppingFlag               uint8                      `json:"frame_cropping_flag"`                             // 1 bit
	FrameCropLeftOffset             *expgolombcoding.Unsigned  `json:"frame_crop_left_offset,omitempty"`                // Exp-Golomb-coded
	FrameCropRightOffset            *expgolombcoding.Unsigned  `json:"frame_crop_right_offset,omitempty"`               // Exp-Golomb-coded
	FrameCropTopOffset              *expgolombcoding.Unsigned  `json:"frame_crop_top_offset,omitempty"`                 // Exp-Golomb-coded
	FrameCropBottomOffset           *expgolombcoding.Unsigned  `json:"frame_crop_bottom_offset,omitempty"`              // Exp-Golomb-coded
	VuiParametersPresentFlag        uint8                      `json:"vui_parameters_present_flag"`                     // 1 bit
	VUIParameters                   *VUIParameters             `json:"vui_parameters,omitempty"`
}

// Parse parses bytes to AVC SPS NAL Unit, return parsed bytes or error.
func (s *SequenceParameterSetData) Parse(r io.Reader, size int) (uint64, error) {
	var parsedBits uint64
	br := bitreader.New(r) // start bit-level parsing here

	if nextByte, err := br.ReadByte(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.ProfileIdc = nextByte
		s.ProfileIdcName = ProfileName(s.ProfileIdc)
		parsedBits += bitsPerByte
	}

	if nextByte, err := br.ReadByte(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.ConstraintSet0Flag = (nextByte >> 7) & 0x1
		s.ConstraintSet1Flag = (nextByte >> 6) & 0x1
		s.ConstraintSet2Flag = (nextByte >> 5) & 0x1
		s.ConstraintSet3Flag = (nextByte >> 4) & 0x1
		s.ConstraintSet4Flag = (nextByte >> 3) & 0x1
		s.ConstraintSet5Flag = (nextByte >> 2) & 0x1
		parsedBits += bitsPerByte
	}

	if nextByte, err := br.ReadByte(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.LevelIdc = nextByte
		parsedBits += bitsPerByte
	}

	if costBits, err := s.SeqParameterSetID.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	// ISO/IEC-14496-10 7.3.2.1.1
	if s.ProfileIdc == 100 || s.ProfileIdc == 110 || s.ProfileIdc == 122 ||
		s.ProfileIdc == 244 || s.ProfileIdc == 44 || s.ProfileIdc == 83 ||
		s.ProfileIdc == 86 || s.ProfileIdc == 118 || s.ProfileIdc == 128 {

		expUnsigned := &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.ChromaFormatIdc = expUnsigned

		if s.ChromaFormatIdc.Value() == 3 {
			if nextBit, err := br.ReadBit(); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				s.SeparateColourPlaneFlag = &nextBit
				parsedBits++
			}
		}

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.BitDepthLumaMinus8 = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.BitDepthChromaMinus8 = expUnsigned

		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits, err
		} else {
			s.QpprimeYZeroTransformBypassFlag = &nextBit
			parsedBits++
		}

		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			s.SeqScalingMatrixPresentFlag = &nextBit
			parsedBits++
		}

		if s.SeqScalingMatrixPresentFlag != nil && *s.SeqScalingMatrixPresentFlag != 0 {

			scalingListPresentFlagLen := 12
			if s.ChromaFormatIdc.Value() != 3 {
				scalingListPresentFlagLen = 8
			}

			for i := 0; i < scalingListPresentFlagLen; i++ {
				if nextBit, err := br.ReadBit(); err != nil {
					return parsedBits / bitsPerByte, err
				} else {
					s.SeqScalingListPresentFlag = append(s.SeqScalingListPresentFlag, nextBit)
					parsedBits++
				}

				if s.SeqScalingListPresentFlag[i] != 0 {
					if i < 6 {
						sp := scalingListParser{sizeOfScalingList: 16}
						if costBits, err := sp.parse(br); err != nil {
							return parsedBits / bitsPerByte, err
						} else {
							parsedBits += costBits
						}
						s.ScalingList4x4 = append(s.ScalingList4x4, sp.scalingList)
						s.DeltaScale = append(s.DeltaScale, sp.deltaScale)
					} else {
						sp := scalingListParser{sizeOfScalingList: 64}
						if costBits, err := sp.parse(br); err != nil {
							return parsedBits / bitsPerByte, err
						} else {
							parsedBits += costBits
						}
						s.ScalingList8x8 = append(s.ScalingList8x8, sp.scalingList)
						s.DeltaScale = append(s.DeltaScale, sp.deltaScale)
					}
				}
				glog.Warningf("Scaling List parsing maybe wrong, check the results first!!!")
			}
		}
	}

	if costBits, err := s.Log2MaxFrameNumMinus4.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	if costBits, err := s.PicOrderCntType.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	if s.PicOrderCntType.Value() == 0 {
		expUnsigned := &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.Log2MaxPicOrderCntLsbMinus4 = expUnsigned
	} else if s.PicOrderCntType.Value() == 1 {
		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			s.DeltaPicOrderAlwaysZeroFlag = &nextBit
			parsedBits++
		}

		expSigned := &expgolombcoding.Signed{}
		if costBits, err := expSigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.OffsetForNonRefPic = expSigned

		expSigned = &expgolombcoding.Signed{}
		if costBits, err := expSigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.OffsetForTopToBottomField = expSigned

		expUnsigned := &expgolombcoding.Unsigned{}
		if costBits, err := expSigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.NumRefFramesInPicOrderCntCycle = expUnsigned

		for i := 0; i < int(s.NumRefFramesInPicOrderCntCycle.Value()); i++ {
			expSigned := expgolombcoding.Signed{}
			if costBits, err := expSigned.Parse(br); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				s.OffsetForRefFrame = append(s.OffsetForRefFrame, expSigned)
				parsedBits += costBits
			}
		}
	}

	if costBits, err := s.MaxNumRefFrames.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.GapsInFrameNumValueAllowedFlag = nextBit
		parsedBits++
	}

	if costBits, err := s.PicWidthInMbsMinus1.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	if costBits, err := s.PicHeightInMapUnitsMinus1.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.FrameMbsOnlyFlag = nextBit
		parsedBits++
	}

	if s.FrameMbsOnlyFlag == 0 {
		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			s.MbAdaptiveFrameFieldFlag = &nextBit
			parsedBits++
		}
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.Direct8x8InferenceFlag = nextBit
		parsedBits++
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.FrameCroppingFlag = nextBit
		parsedBits++
	}

	if s.FrameCroppingFlag != 0 {
		expUnsigned := &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.FrameCropLeftOffset = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.FrameCropRightOffset = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.FrameCropTopOffset = expUnsigned

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		s.FrameCropBottomOffset = expUnsigned
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		s.VuiParametersPresentFlag = nextBit
		parsedBits++
	}

	if s.VuiParametersPresentFlag != 0 {
		s.VUIParameters = &VUIParameters{}
		if costBits, err := s.VUIParameters.parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
	}

	if br.CachedBitsCount() > 0 {
		ignoreBits := uint(br.CachedBitsCount())
		if _, err := br.ReadBits(ignoreBits); err != nil { // ignore rbsp_stop_one_bit and several rbsp_alignment_zero_bit
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += uint64(ignoreBits)
		}
	}

	// bits to bytes
	parsedBytes := parsedBits / bitsPerByte
	if parsedBits%bitsPerByte != 0 {
		glog.Warningf("parsed bits doesn't align in 8 bits, total %d bits", parsedBits)
		parsedBytes += 1
	}

	if int(parsedBytes) != size {
		glog.Warningf("parsed bytes != expect size : %d!=%d", parsedBytes, size)
	}
	return parsedBytes, nil
}

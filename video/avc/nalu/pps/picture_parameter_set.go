// Package pps defined AVC Picture Parameter Sets information.
package pps

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// PictureParameterSet represents pic_parameter_set_rbsp defined in ISO/IEC-14496-10 7.3.2.2.
type PictureParameterSet struct {
	PicParameterSetId                     expgolombcoding.Unsigned   `json:"pic_parameter_set_id"`
	SeqParameterSetId                     expgolombcoding.Unsigned   `json:"seq_parameter_set_id"`
	EntropyCodingModeFlag                 uint8                      `json:"entropy_coding_mode_flag"`
	BottomFieldPicOrderInFramePresentFlag uint8                      `json:"bottom_field_pic_order_in_frame_present_flag"`
	NumSliceGroupsMinus1                  expgolombcoding.Unsigned   `json:"num_slice_groups_minus1"`
	SliceGroupMapType                     *expgolombcoding.Unsigned  `json:"slice_group_map_type,omitempty"`
	RunLengthMinus1                       []expgolombcoding.Unsigned `json:"run_length_minus1,omitempty"`
	TopLeft                               []expgolombcoding.Unsigned `json:"top_left,omitempty"`
	BottomRight                           []expgolombcoding.Unsigned `json:"bottom_right,omitempty"`
	SliceGroupChangeDirectionFlag         *uint8                     `json:"slice_group_change_direction_flag,omitempty"`
	SliceGroupChangeRateMinus1            *expgolombcoding.Unsigned  `json:"slice_group_change_rate_minus1,omitempty"`
	PicSizeInMapUnitsMinus1               *expgolombcoding.Unsigned  `json:"pic_size_in_map_units_minus1,omitempty"`
	SliceGroupId                          []uint8                    `json:"slice_group_id,omitempty"`
	NumRefIdxL0DefaultActiveMinus1        expgolombcoding.Unsigned   `json:"num_ref_idx_l0_default_active_minus1"`
	NumRefIdxL1DefaultActiveMinus1        expgolombcoding.Unsigned   `json:"num_ref_idx_l1_default_active_minus1"`
	WeightedPredFlag                      uint8                      `json:"weighted_pred_flag"`
	WeightedBipredIdc                     uint8                      `json:"weighted_bipred_idc"`
	PicInitQpMinus26                      expgolombcoding.Signed     `json:"pic_init_qp_minus26"`
	PicInitQsMinus26                      expgolombcoding.Signed     `json:"pic_init_qs_minus26"`
	ChromaQpIndexOffset                   expgolombcoding.Signed     `json:"chroma_qp_index_offset"`
	DeblockingFilterControlPresentFlag    uint8                      `json:"deblocking_filter_control_present_flag"`
	ConstrainedIntraPredFlag              uint8                      `json:"constrained_intra_pred_flag"`
	RedundantPicCntPresentFlag            uint8                      `json:"redundant_pic_cnt_present_flag"`
	Transform8x8ModeFlag                  *uint8                     `json:"transform_8x8_mode_flag,omitempty"`
	PicScalingMatrixPresentFlag           *uint8                     `json:"pic_scaling_matrix_present_flag,omitempty"`
	PicScalingListPresentFlag             []uint8                    `json:"pic_scaling_list_present_flag,omitempty"`
	SecondChromaQpIndexOffset             *expgolombcoding.Signed    `json:"second_chroma_qp_index_offset,omitempty"`

	// internal fields
	sps *sps.SequenceParameterSetData `json:"-"`
}

const bitsPerByte = 8

// SetSPS sets SequenceParameterSetData for parsing.
func (p *PictureParameterSet) SetSPS(sps *sps.SequenceParameterSetData) {
	p.sps = sps
}

// Parse parses bytes to AVC SPS NAL Unit, return parsed bytes or error.
func (p *PictureParameterSet) Parse(r io.Reader, size int) (uint64, error) {
	var parsedBits uint64
	br := bitreader.New(r) // start bit-level parsing here

	expUnsigned := &expgolombcoding.Unsigned{}
	if costBits, err := expUnsigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.PicParameterSetId = *expUnsigned

	expUnsigned = &expgolombcoding.Unsigned{}
	if costBits, err := expUnsigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.SeqParameterSetId = *expUnsigned

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.EntropyCodingModeFlag = nextBit
		parsedBits++
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.BottomFieldPicOrderInFramePresentFlag = nextBit
		parsedBits++
	}

	expUnsigned = &expgolombcoding.Unsigned{}
	if costBits, err := expUnsigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.NumSliceGroupsMinus1 = *expUnsigned

	if p.NumSliceGroupsMinus1.Value() > 0 {
		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		p.SliceGroupMapType = expUnsigned

		if p.SliceGroupMapType.Value() == 0 {
			for i := 0; uint64(i) < p.NumSliceGroupsMinus1.Value(); i++ {
				expUnsigned = &expgolombcoding.Unsigned{}
				if costBits, err := expUnsigned.Parse(br); err != nil {
					return parsedBits / bitsPerByte, err
				} else {
					parsedBits += costBits
				}
				p.RunLengthMinus1 = append(p.RunLengthMinus1, *expUnsigned)
			}
		} else if p.SliceGroupMapType.Value() == 2 {
			for i := 0; uint64(i) < p.NumSliceGroupsMinus1.Value(); i++ {
				expUnsigned = &expgolombcoding.Unsigned{}
				if costBits, err := expUnsigned.Parse(br); err != nil {
					return parsedBits / bitsPerByte, err
				} else {
					parsedBits += costBits
				}
				p.TopLeft = append(p.TopLeft, *expUnsigned)

				expUnsigned = &expgolombcoding.Unsigned{}
				if costBits, err := expUnsigned.Parse(br); err != nil {
					return parsedBits / bitsPerByte, err
				} else {
					parsedBits += costBits
				}
				p.BottomRight = append(p.BottomRight, *expUnsigned)
			}
		} else if p.SliceGroupMapType.Value() == 3 ||
			p.SliceGroupMapType.Value() == 4 ||
			p.SliceGroupMapType.Value() == 5 {

			if nextBit, err := br.ReadBit(); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				p.SliceGroupChangeDirectionFlag = &nextBit
				parsedBits++
			}

			expUnsigned = &expgolombcoding.Unsigned{}
			if costBits, err := expUnsigned.Parse(br); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				parsedBits += costBits
			}
			p.SliceGroupChangeRateMinus1 = expUnsigned
		} else if p.SliceGroupMapType.Value() == 6 {
			expUnsigned = &expgolombcoding.Unsigned{}
			if costBits, err := expUnsigned.Parse(br); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				parsedBits += costBits
			}
			p.PicSizeInMapUnitsMinus1 = expUnsigned

			for i := 0; uint64(i) < p.PicSizeInMapUnitsMinus1.Value(); i++ {
				if nextBit, err := br.ReadBit(); err != nil {
					return parsedBits / bitsPerByte, err
				} else {
					p.SliceGroupId = append(p.SliceGroupId, nextBit)
					parsedBits++
				}
			}
		}
	}

	expUnsigned = &expgolombcoding.Unsigned{}
	if costBits, err := expUnsigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.NumRefIdxL0DefaultActiveMinus1 = *expUnsigned

	expUnsigned = &expgolombcoding.Unsigned{}
	if costBits, err := expUnsigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.NumRefIdxL1DefaultActiveMinus1 = *expUnsigned

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.WeightedPredFlag = nextBit
		parsedBits++
	}

	if nextBits, err := br.ReadBits(2); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.WeightedBipredIdc = nextBits[0]
		parsedBits += 2
	}

	expSigned := &expgolombcoding.Signed{}
	if costBits, err := expSigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.PicInitQpMinus26 = *expSigned

	expSigned = &expgolombcoding.Signed{}
	if costBits, err := expSigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.PicInitQsMinus26 = *expSigned

	expSigned = &expgolombcoding.Signed{}
	if costBits, err := expSigned.Parse(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}
	p.ChromaQpIndexOffset = *expSigned

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.DeblockingFilterControlPresentFlag = nextBit
		parsedBits++
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.ConstrainedIntraPredFlag = nextBit
		parsedBits++
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.RedundantPicCntPresentFlag = nextBit
		parsedBits++
	}

	// return if no more rbsp data
	if parsedBits >= uint64(size)*parsedBits {
		return parsedBits / bitsPerByte, nil
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.Transform8x8ModeFlag = &nextBit
		parsedBits++
	}

	if nextBit, err := br.ReadBit(); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		p.PicScalingMatrixPresentFlag = &nextBit
		parsedBits++
	}

	if *p.PicScalingMatrixPresentFlag != 0 {
		loopCount := 6
		if p.sps.ChromaFormatIdc.Value() != 3 {
			loopCount = 2
		}
		loopCount *= int(*p.Transform8x8ModeFlag)
		loopCount += 6
		for i := 0; i < loopCount; i++ {
			if nextBit, err := br.ReadBit(); err != nil {
				return parsedBits / bitsPerByte, err
			} else {
				p.PicScalingListPresentFlag = append(p.PicScalingListPresentFlag, nextBit)
				parsedBits++
			}

			//TODO: scaling list parse
			glog.Warningf("pps scaling list parse doesn't support yet, refer to SPS implementation and do it here")
		}

		expSigned = &expgolombcoding.Signed{}
		if costBits, err := expSigned.Parse(br); err != nil {
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += costBits
		}
		p.SecondChromaQpIndexOffset = expSigned
	}
	if br.CachedBitsCount() > 0 {
		ignoreBits := uint(br.CachedBitsCount())
		if _, err := br.ReadBits(ignoreBits); err != nil { // ignore rbsp_stop_one_bit and several rbsp_alignment_zero_bit
			return parsedBits / bitsPerByte, err
		} else {
			parsedBits += uint64(ignoreBits)
		}
	}

	if parsedBits != uint64(size)*bitsPerByte {
		glog.Infof("parsed bits %d but expect bits %d", parsedBits, size*bitsPerByte)
	}
	return parsedBits / bitsPerByte, nil
}

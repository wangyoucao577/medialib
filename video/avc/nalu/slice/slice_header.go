// Package slice reprensents AVC Slice structures.
package slice

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
)

// Header represents slice header
type Header struct {
	FirstMBInSlice    expgolombcoding.Unsigned `json:"first_mb_in_slice"`
	SliceType         expgolombcoding.Unsigned `json:"slice_type"`
	PicParameterSetID expgolombcoding.Unsigned `json:"pic_parameter_set_id"`
	ColourPlaneId     *uint8                   `json:"colour_plane_id,omitempty"` // 2 bits
	FrameNum          uint64                   `json:"frame_num"`
	FieldPicFlag      *uint8                   `json:"field_pic_flag,omitempty"`
	BottomPicFlag     *uint8                   `json:"bottom_pic_flag,omitempty"`
	// IdrPicId                    *expgolombcoding.Unsigned  `json:"idr_pic_id,omitempty"`
	// PicOrderCntLsb              *uint64                    `json:"pic_order_cnt_lsb,omitempty"`
	// DeltaPicOrderCntBottom      *expgolombcoding.Signed    `json:"delta_pic_order_cnt_bottom,omitempty"`
	// DeltaPicOrderCnt            [2]*expgolombcoding.Signed `json:"delta_pic_order_cnt"`
	// RedundantPicCnt             *expgolombcoding.Unsigned  `json:"redundant_pic_cnt,omitempty"`
	// DirectSpatialMvPredFlag     *uint8                     `json:"direct_spatial_mv_pred_flag,omitempty"`
	// NumRefIdxActiveOverrideFlag *uint8                     `json:"num_ref_idx_active_override_flag,omitempty"`
	// NumRefIdxL0ActiveMinus1     *expgolombcoding.Unsigned  `json:"num_ref_idx_l0_active_minus1,omitempty"`
	// NumRefIdxL1ActiveMinus1     *expgolombcoding.Unsigned  `json:"num_ref_idx_l1_active_minus1,omitempty"`
	//TODO: ref_pic_list_mvc_modification or ref_pic_list_modification
}

// MarshalJSON implements json.Marshaler.
func (h *Header) MarshalJSON() ([]byte, error) {
	var hj = struct {
		FirstMBInSlice expgolombcoding.Unsigned `json:"first_mb_in_slice"`
		SliceType      expgolombcoding.Unsigned `json:"slice_type"`
		SliceTypeName  string                   `json:"slice_type_name"`

		PicParameterSetID expgolombcoding.Unsigned `json:"pic_parameter_set_id"`
		ColourPlaneId     *uint8                   `json:"colour_plane_id,omitempty"` // 2 bits
		FrameNum          uint64                   `json:"frame_num"`
		FieldPicFlag      *uint8                   `json:"field_pic_flag,omitempty"`
		BottomPicFlag     *uint8                   `json:"bottom_pic_flag,omitempty"`
	}{
		FirstMBInSlice: h.FirstMBInSlice,
		SliceType:      h.SliceType,
		SliceTypeName:  Type(int(h.SliceType.Value())),

		PicParameterSetID: h.PicParameterSetID,
		ColourPlaneId:     h.ColourPlaneId,
		FrameNum:          h.FrameNum,
		FieldPicFlag:      h.FieldPicFlag,
		BottomPicFlag:     h.BottomPicFlag,
	}

	return json.Marshal(hj)
}

// return parsed bits
func (l *LayerWithoutPartitioningRbsp) parseHeader(br *bitreader.Reader) (uint64, error) {
	h := &l.Header

	var parsedBits uint64

	if costBits, err := h.FirstMBInSlice.Parse(br); err != nil {
		return parsedBits, err
	} else {
		parsedBits += costBits
	}

	if costBits, err := h.SliceType.Parse(br); err != nil {
		return parsedBits, err
	} else {
		parsedBits += costBits
	}

	if costBits, err := h.PicParameterSetID.Parse(br); err != nil {
		return parsedBits, err
	} else {
		parsedBits += costBits
	}

	if l.sps == nil {
		return parsedBits, ErrEmptyParameterSet
	}
	if l.sps.SeparateColourPlaneFlag != nil && *l.sps.SeparateColourPlaneFlag == 1 {
		if nextBits, err := br.ReadBits(2); err != nil {
			return parsedBits, err
		} else {
			h.ColourPlaneId = &nextBits[0]
			parsedBits += 2
		}
	}

	frameNumBits := l.sps.Log2MaxFrameNumMinus4.Value() + 4
	if nextBits, err := br.ReadBits(uint(frameNumBits)); err != nil {
		return parsedBits, err
	} else {
		if len(nextBits) == 1 {
			h.FrameNum = uint64(nextBits[0])
		} else if len(nextBits) == 2 {
			h.FrameNum = uint64(binary.BigEndian.Uint16(nextBits))
		} else if len(nextBits) == 3 {
			nextBits = append([]byte{0x0}, nextBits...)
			h.FrameNum = uint64(binary.BigEndian.Uint32(nextBits))
		} else if len(nextBits) < 8 {
			nextBits = append(make([]byte, 8-len(nextBits)), nextBits...)
			h.FrameNum = uint64(binary.BigEndian.Uint64(nextBits))
		} else if len(nextBits) == 8 {
			h.FrameNum = uint64(binary.BigEndian.Uint64(nextBits))
		} else {
			return parsedBits, fmt.Errorf("can not parse frame num bits %v", nextBits)
		}

		parsedBits += frameNumBits
	}

	if l.sps.FrameMbsOnlyFlag == 0 {
		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits, err
		} else {
			h.FieldPicFlag = &nextBit
			parsedBits += 1
		}

		if *h.FieldPicFlag != 0 {
			if nextBit, err := br.ReadBit(); err != nil {
				return parsedBits, err
			} else {
				h.BottomPicFlag = &nextBit
				parsedBits += 1
			}
		}
	}

	//TODO: parse next
	glog.V(3).Infof("slice header idr_pic_id parsing TODO")

	return parsedBits, nil
}

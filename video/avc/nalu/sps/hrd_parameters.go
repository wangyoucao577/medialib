package sps

import (
	"fmt"

	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
)

// HrdParameters represents hrd_parameters defined in ISO/IEC-14496-10 Annex E.1.2
type HrdParameters struct {
	CpbCntMinus1                       expgolombcoding.Unsigned   `json:"cpb_cnt_minus1"`
	BitRateScale                       uint8                      `json:"bit_rate_scale"` // 4 bits
	CpbSizeScale                       uint8                      `json:"cpb_size_scale"` // 4 bits
	BitRateValueMinus1                 []expgolombcoding.Unsigned `json:"bit_rate_value_minus1,omitempty"`
	CpbSizeValueMinus1                 []expgolombcoding.Unsigned `json:"cpb_size_value_minus1,omitempty"`
	CbrFlag                            []uint8                    `json:"cbr_flag,omitempty"`                      // 1 bit per flag
	InitialCpbRemovalDelayLengthMinus1 uint8                      `json:"initial_cpb_removal_delay_length_minus1"` // 5 bits
	CpbRemovalDelayLengthMinus1        uint8                      `json:"cpb_removal_delay_length_minus1"`         // 5 bits
	DpbOutputDelayLengthMinus1         uint8                      `json:"dpb_output_delay_length_minus1"`          // 5 bits
	TimeOffsetLength                   uint8                      `json:"time_offset_length"`                      // 5 bits
}

// return parsed bits
func (h *HrdParameters) parse(br *bitreader.Reader) (uint64, error) {
	//TODO:
	return 0, fmt.Errorf("hrd_parameters parsing doesn't impemented yet")
}

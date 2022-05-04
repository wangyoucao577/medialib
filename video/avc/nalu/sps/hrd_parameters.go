package sps

import (
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
	var parsedBits uint64

	expUnsigned := &expgolombcoding.Unsigned{}
	if costBits, err := expUnsigned.Parse(br); err != nil {
		return parsedBits, err
	} else {
		parsedBits += costBits
	}
	h.CpbCntMinus1 = *expUnsigned

	if nextByte, err := br.ReadByte(); err != nil {
		return parsedBits, err
	} else {
		h.BitRateScale = (nextByte >> 4) & 0xF
		h.CpbSizeScale = nextByte & 0xF
		parsedBits += 8
	}

	for i := 0; i <= int(h.CpbCntMinus1.Value()); i++ {

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits, err
		} else {
			parsedBits += costBits
		}
		h.BitRateValueMinus1 = append(h.BitRateValueMinus1, *expUnsigned)

		expUnsigned = &expgolombcoding.Unsigned{}
		if costBits, err := expUnsigned.Parse(br); err != nil {
			return parsedBits, err
		} else {
			parsedBits += costBits
		}
		h.CpbSizeValueMinus1 = append(h.CpbSizeValueMinus1, *expUnsigned)

		if nextBit, err := br.ReadBit(); err != nil {
			return parsedBits, err
		} else {
			h.CbrFlag = append(h.CbrFlag, nextBit)
			parsedBits += 1
		}
	}

	if nextBits, err := br.ReadBits(5); err != nil {
		return parsedBits, err
	} else {
		h.InitialCpbRemovalDelayLengthMinus1 = nextBits[0] & 0x1F
		parsedBits += 5
	}

	if nextBits, err := br.ReadBits(5); err != nil {
		return parsedBits, err
	} else {
		h.CpbRemovalDelayLengthMinus1 = nextBits[0] & 0x1F
		parsedBits += 5
	}

	if nextBits, err := br.ReadBits(5); err != nil {
		return parsedBits, err
	} else {
		h.DpbOutputDelayLengthMinus1 = nextBits[0] & 0x1F
		parsedBits += 5
	}

	if nextBits, err := br.ReadBits(5); err != nil {
		return parsedBits, err
	} else {
		h.TimeOffsetLength = nextBits[0] & 0x1F
		parsedBits += 5
	}

	return parsedBits, nil
}

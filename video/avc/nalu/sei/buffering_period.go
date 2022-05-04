package sei

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// BufferingPeriod represents AVC SEI buffering_period.
type BufferingPeriod struct {
	SeqParameterSetID            expgolombcoding.Unsigned `json:"seq_parameter_set_id"`
	InitialCpbRemovalDelay       []uint                   `json:"initial_cpb_removal_delay,omitempty"`
	InitialCpbRemovalDelayOffset []uint                   `json:"initial_cpb_removal_delay_offset,omitempty"`

	// store for some internal parsing
	sps *sps.SequenceParameterSetData `json:"-"`
	pps *pps.PictureParameterSet      `json:"-"`
}

// setSequenceHeaders sets both SPS and PPS for parsing.
func (b *BufferingPeriod) setSequenceHeaders(sps *sps.SequenceParameterSetData, pps *pps.PictureParameterSet) {
	b.sps = sps
	b.pps = pps
}

func (b *BufferingPeriod) parse(r io.Reader, payloadSize int) (uint64, error) {

	if b.sps == nil || b.sps.VuiParametersPresentFlag == 0 || b.sps.VUIParameters == nil {
		return 0, fmt.Errorf("invalid sps %v", b.sps)
	}

	var parsedBits uint64
	br := bitreader.New(r)

	b.SeqParameterSetID = expgolombcoding.Unsigned{}
	if costBits, err := b.SeqParameterSetID.Parse(br); err != nil {
		return parsedBits / 8, err
	} else {
		parsedBits += costBits
	}

	var nalHrdBpPresentFlag bool
	if b.sps.VUIParameters.NalHrdParametersPresentFlag == 1 {
		nalHrdBpPresentFlag = true
	}
	var vclHrdBpPresentFlag bool
	if b.sps.VUIParameters.VclHrdParametersPresentFlag == 1 {
		vclHrdBpPresentFlag = true
	}

	var cpbCntMinus1 int
	var initialCpbSizeMinus1 int
	if nalHrdBpPresentFlag {
		cpbCntMinus1 = int(b.sps.VUIParameters.NalHrdParmaeters.CpbCntMinus1.Value())
		initialCpbSizeMinus1 = int(b.sps.VUIParameters.NalHrdParmaeters.InitialCpbRemovalDelayLengthMinus1)
	}
	if vclHrdBpPresentFlag {
		cpbCntMinus1 = int(b.sps.VUIParameters.VclHrdParmaeters.CpbCntMinus1.Value())
		initialCpbSizeMinus1 = int(b.sps.VUIParameters.VclHrdParmaeters.InitialCpbRemovalDelayLengthMinus1)
	}

	for i := 0; i <= cpbCntMinus1; i++ {
		cpbDelaybits := initialCpbSizeMinus1 + 1
		if cpbDelaybits > 32 { // > 4 bytes
			return parsedBits / 8, fmt.Errorf("unsupported cpb removal delay bits %d", cpbDelaybits)
		}

		if nextBits, err := br.ReadBits(uint(cpbDelaybits)); err != nil {
			return parsedBits / 8, err
		} else {
			if len(nextBits) == 1 {
				b.InitialCpbRemovalDelay = append(b.InitialCpbRemovalDelay, uint(nextBits[0]))
			} else if len(nextBits) == 2 {
				b.InitialCpbRemovalDelay = append(b.InitialCpbRemovalDelay, uint(binary.BigEndian.Uint16(nextBits)))
			} else if len(nextBits) == 3 {
				b.InitialCpbRemovalDelay = append(b.InitialCpbRemovalDelay, uint(binary.BigEndian.Uint32([]byte{0x00, nextBits[0], nextBits[1], nextBits[2]})))
			} else if len(nextBits) == 4 {
				b.InitialCpbRemovalDelay = append(b.InitialCpbRemovalDelay, uint(binary.BigEndian.Uint32(nextBits)))
			}
			parsedBits += uint64(cpbDelaybits)
		}

		if nextBits, err := br.ReadBits(uint(cpbDelaybits)); err != nil {
			return parsedBits / 8, err
		} else {
			if len(nextBits) == 1 {
				b.InitialCpbRemovalDelayOffset = append(b.InitialCpbRemovalDelayOffset, uint(nextBits[0]))
			} else if len(nextBits) == 2 {
				b.InitialCpbRemovalDelayOffset = append(b.InitialCpbRemovalDelayOffset, uint(binary.BigEndian.Uint16(nextBits)))
			} else if len(nextBits) == 3 {
				b.InitialCpbRemovalDelayOffset = append(b.InitialCpbRemovalDelayOffset, uint(binary.BigEndian.Uint32([]byte{0x00, nextBits[0], nextBits[1], nextBits[2]})))
			} else if len(nextBits) == 4 {
				b.InitialCpbRemovalDelayOffset = append(b.InitialCpbRemovalDelayOffset, uint(binary.BigEndian.Uint32(nextBits)))
			}
			parsedBits += uint64(cpbDelaybits)
		}
	}

	return parsedBits / 8, nil
}

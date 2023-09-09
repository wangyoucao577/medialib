package av1c

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// AV1CodecConfigurationRecord defines AV1 Codec configuration record.
type AV1CodecConfigurationRecord struct {
	Marker  uint8 `json:"marker"`  // 1bit
	Version uint8 `json:"version"` // 7bits

	SeqProfile   uint8 `json:"seq_profile"`     // 3bits
	SeqLevelIdx0 uint8 `json:"seq_level_idx_0"` // 5bits

	SeqTier0             uint8 `json:"seq_tier_0"`             // 1bit
	HighBitdepth         uint8 `json:"high_bitdepth"`          // 1bit
	TwelveBit            uint8 `json:"twelve_bit"`             // 1bit
	Monochrome           uint8 `json:"monochrome"`             // 1bit
	ChromaSubsamplingX   uint8 `json:"chroma_subsampling_x"`   // 1bit
	ChromaSubsamplingY   uint8 `json:"chroma_subsampling_y"`   // 1bit
	ChromaSamplePosition uint8 `json:"chroma_sample_position"` // 2bits

	// then 3 bits reserved here
	InitialPresentationDelayPresent  uint8  `json:"initial_presentation_delay_present"`             // 1bit
	InitialPresentationDelayMinusOne *uint8 `json:"initial_presentation_delay_minus_one,omitempty"` // 4bits or not exist

	ConfigOBUs []uint8 `json:"configOBUs,omitempty"`
}

// Parse parses AV1CodecConfigurationRecord.
func (a *AV1CodecConfigurationRecord) Parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	data := make([]byte, 4)

	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		a.Marker = (data[0] >> 7) & 0x1
		a.Version = data[0] & 0x7F

		a.SeqProfile = (data[1] >> 5) & 0x7
		a.SeqLevelIdx0 = data[1] & 0x1F

		a.SeqTier0 = (data[2] >> 7) & 0x1
		a.HighBitdepth = (data[2] >> 6) & 0x1
		a.TwelveBit = (data[2] >> 5) & 0x1
		a.Monochrome = (data[2] >> 4) & 0x1
		a.ChromaSubsamplingX = (data[2] >> 3) & 0x1
		a.ChromaSubsamplingY = (data[2] >> 2) & 0x1
		a.ChromaSamplePosition = data[2] & 0x3

		a.InitialPresentationDelayPresent = (data[3] >> 4) & 0x1
		if a.InitialPresentationDelayPresent > 0 {
			num := (data[3] & 0xF)
			a.InitialPresentationDelayMinusOne = &num
		}

		parsedBytes += 4
	}

	// TODO: configOBUs
	glog.Warning("TODO: parse configOBUs of av1C")

	return parsedBytes, nil
}

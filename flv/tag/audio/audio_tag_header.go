package audio

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// TagHeader reprensets Audio Tag Header.
type TagHeader struct {
	// Format of SoundData. The following values are defined:
	// 0 = Linear PCM, platform endian
	// 1 = ADPCM
	// 2 = MP3
	// 3 = Linear PCM, little endian 4 = Nellymoser 16 kHz mono 5 = Nellymoser 8 kHz mono 6 = Nellymoser
	// 7 = G.711 A-law logarithmic PCM
	// 8 = G.711 mu-law logarithmic PCM
	// 9 = reserved
	// 10 = AAC
	// 11 = Speex
	// 14 = MP3 8 kHz
	// 15 = Device-specific sound
	// Formats 7, 8, 14, and 15 are reserved.
	// AAC is supported in Flash Player 9,0,115,0 and higher. Speex is supported in Flash Player 10 and higher.
	SoundFormat uint8 `json:"SoundFormat"` // 4 bits

	// Sampling rate. The following values are defined:
	// 0 = 5.5 kHz
	// 1 = 11 kHz
	// 2 = 22 kHz
	// 3 = 44 kHz
	SoundRate uint8 `json:"SoundRate"` // 2 bits

	// Size of each audio sample. This parameter only pertains to uncompressed formats. Compressed formats always decode to 16 bits internally.
	// 0 = 8-bit samples
	// 1 = 16-bit samples
	SoundSize uint8 `json:"SoundSize"` // 1 bit

	// Mono or stereo sound
	// 0 = Mono sound
	// 1 = Stereo sound
	SoundType uint8 `json:"SoundType"` // 1 bit

	// 	The following values are defined:
	// 0 = AAC sequence header
	// 1 = AAC raw
	AACPacketType *uint8 `json:"AACPacketType"` // 8 bits if exist
}

// MarshalJSON implements json.Marshaler.
func (t *TagHeader) MarshalJSON() ([]byte, error) {
	var tj = struct {
		SoundRate            uint8  `json:"SoundRate"`
		SoundRateDescription string `json:"SoundRateDescription"`

		SoundSize            uint8  `json:"SoundSize"`
		SoundSizeDescription string `json:"SoundSizeDescription"`

		SoundType            uint8  `json:"SoundType"`
		SoundTypeDescription string `json:"SoundTypeDescription"`

		AACPacketType            *uint8 `json:"AACPacketType"`
		AACPacketTypeDescription string `json:"AACPacketTypeDescription,omitempty"`
	}{
		SoundRate:            t.SoundRate,
		SoundRateDescription: SoundFormatDescription(int(t.SoundFormat)),

		SoundSize:            t.SoundSize,
		SoundSizeDescription: SoundSizeDescription(int(t.SoundSize)),

		SoundType:            t.SoundType,
		SoundTypeDescription: SoundTypeDescription(int(t.SoundType)),

		AACPacketType:            t.AACPacketType,
		AACPacketTypeDescription: AACPacketTypeDescription(int(*t.AACPacketType)),
	}

	return json.Marshal(tj)
}

func (t *TagHeader) parse(r io.Reader) (uint64, error) {
	var parsedBytes uint64

	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		t.SoundFormat = (data[0] >> 4) & 0xF
		t.SoundRate = (data[0] >> 2) & 0x3
		t.SoundSize = (data[0] >> 1) & 0x1
		t.SoundType = data[0] & 0x1
		parsedBytes += 1
	}

	if t.SoundFormat == SoundFormatAAC {
		if err := util.ReadOrError(r, data); err != nil {
			return parsedBytes, err
		} else {
			t.AACPacketType = &data[0]
			parsedBytes += 1
		}
	} else {
		return parsedBytes, fmt.Errorf("sound format %d(%s) doesn't support yet", t.SoundFormat, SoundFormatDescription(int(t.SoundFormat)))
	}

	return parsedBytes, nil
}

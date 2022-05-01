// Package audio represents Audio Tag.
package audio

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv/tag"
	"github.com/wangyoucao577/medialib/util"
)

// TagHeader reprensets Audio Tag Header.
type TagHeader struct {
	// Format of SoundData. The following values are defined: 0 = Linear PCM, platform endian
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

	// 	Sampling rate. The following values are defined: 0 = 5.5 kHz
	// 1 = 11 kHz
	// 2 = 22 kHz
	// 3 = 44 kHz
	SoundRate uint8 `json:"SoundRate"` // 2 bits

	// 	Size of each audio sample. This parameter only pertains to uncompressed formats. Compressed formats always decode to 16 bits internally.
	// 0 = 8-bit samples
	// 1 = 16-bit samples
	SoundSize uint8 `json:"SoundSize"` // 1 bit

	// 	Mono or stereo sound 0 = Mono sound
	// 1 = Stereo sound
	SoundType uint8 `json:"SoundType"` // 1 bit

	// 	The following values are defined: 0 = AAC sequence header
	// 1 = AAC raw
	AACPacketType *uint8 `json:"AACPacketType"` // 8 bits if exist
}

// Tag represents audio tag.
type Tag struct {
	Header         tag.Header `json:"TagHeader"`
	AudioTagHeader TagHeader  `json:"AudioTagHeader"`
	Body           *TagBody   `json:"AudioTagBody"`
}

// GetTagHeader returns tag header.
func (t *Tag) GetTagHeader() tag.Header {
	return t.Header
}

// ParsePayload parses AudioTagHeader and TayBody data with preset tag.Header.
func (t *Tag) ParsePayload(r io.Reader) error {
	if err := t.Header.Validate(); err != nil {
		return err
	}

	//TODO: parse payload
	glog.Warningf("tag type %d doesn't implemented yet, ignore payload size %d", t.Header.TagType, t.Header.DataSize)
	if err := util.ReadOrError(r, make([]byte, t.Header.DataSize)); err != nil {
		return err
	}
	return nil
}

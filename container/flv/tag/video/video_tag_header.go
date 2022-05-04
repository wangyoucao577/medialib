package video

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// TagHeader reprensets Audio Tag Header.
type TagHeader struct {

	// Type of video frame. The following values are defined:
	// 1 = key frame (for AVC, a seekable frame)
	// 2 = inter frame (for AVC, a non-seekable frame)
	// 3 = disposable inter frame (H.263 only)
	// 4 = generated key frame (reserved for server use only)
	// 5 = video info/command frame
	FrameType uint8 `json:"FrameType"` // 4 bits

	// 	Codec Identifier. The following values are defined:
	// 2 = Sorenson H.263
	// 3 = Screen video
	// 4 = On2 VP6
	// 5 = On2 VP6 with alpha channel 6 = Screen video version 2
	// 7 = AVC
	CodecID uint8 `json:"CodecID"` // 4 bits

	// The following values are defined:
	// 0 = AVC sequence header
	// 1 = AVC NALU
	// 2 = AVC end of sequence (lower level NALU sequence ender is not required or supported)
	AVCPacketType *uint8 `json:"AVCPacketType,omitempty"`

	// IF AVCPacketType == 1
	//   Composition time offset
	// ELSE
	//   0
	// See ISO 14496-12, 8.15.3 for an explanation of composition times.
	// The offset in an FLV file is always in milliseconds.
	CompositionTime *int32 `json:"CompositionTime,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (t *TagHeader) MarshalJSON() ([]byte, error) {
	var tj = struct {
		FrameType            uint8  `json:"FrameType"`
		FrameTypeDescription string `json:"FrameTypeDescription"`

		CodecID            uint8  `json:"CodecID"`
		CodecIDDescription string `json:"CodecIDDescription"`

		AVCPacketType            *uint8 `json:"AVCPacketType,omitempty"`
		AVCPacketTypeDescription string `json:"AVCPacketTypeDescription,omitempty"`

		CompositionTime *int32 `json:"CompositionTime,omitempty"`
	}{
		FrameType:            t.FrameType,
		FrameTypeDescription: FrameTypeDescription(int(t.FrameType)),

		CodecID:            t.CodecID,
		CodecIDDescription: CodecIDDescription(int(t.CodecID)),

		AVCPacketType:            t.AVCPacketType,
		AVCPacketTypeDescription: AVCPacketTypeDescription(int(*t.AVCPacketType)),

		CompositionTime: t.CompositionTime,
	}

	return json.Marshal(tj)
}

func (t *TagHeader) parse(r io.Reader) (uint64, error) {
	var parsedBytes uint64

	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += 1
	}

	t.FrameType = (data[0] >> 4) & 0xF
	t.CodecID = data[0] & 0xF

	if t.CodecID == CodecIDAVC {
		data = make([]byte, 4)
		if err := util.ReadOrError(r, data); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += 4
		}

		t.AVCPacketType = &data[0]
		cts := int32(binary.BigEndian.Uint32([]byte{0x00, data[1], data[2], data[3]}))
		t.CompositionTime = &cts
	}

	return parsedBytes, nil
}

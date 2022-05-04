package chunk

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// MessageHeader represents RTMP chunk message header.
type MessageHeader struct {
	Timestamp uint32  `json:"timestamp"`                   // first 24 bits except type == 3
	Length    *uint32 `json:"message_length,omitempty"`    // next 24 bits if (type == 0 || type == 1)
	TypeID    *uint8  `json:"message_type_id,omitempty"`   // next 8 bits if (type == 0 || type == 1)
	StreamID  *uint32 `json:"message_stream_id,omitempty"` // next 32 bits if type == 0
}

// Serialize serializes message header to binary format.
func (m *MessageHeader) Serialize(fmt uint8) []byte {
	if fmt >= MessageHeaderFmt3 { // only 0, 1, 2 need to serialize
		return nil
	}

	data := []byte{}

	timestampBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(timestampBytes, m.Timestamp)
	data = append(data, timestampBytes[1:]...)

	if fmt <= MessageHeaderFmt1 {
		if m.Length != nil {
			lenData := make([]byte, 4)
			binary.BigEndian.PutUint32(timestampBytes, *m.Length)
			data = append(data, lenData[1:]...)
		}

		if m.TypeID != nil {
			data = append(data, *m.TypeID)
		}
	}
	if fmt == MessageHeaderFmt0 {
		if m.StreamID != nil {
			streamIDData := make([]byte, 4)
			binary.BigEndian.PutUint32(streamIDData, *m.StreamID)
			data = append(data, streamIDData...)
		}
	}

	return data
}

// Parse parses message header from binary format.
func (m *MessageHeader) Parse(r io.Reader, fmt uint8) (uint64, error) {
	if fmt >= MessageHeaderFmt3 { // only 0, 1, 2 need to parse
		return 0, nil
	}

	var parsedBytes uint

	data := make([]byte, 4)
	if err := util.ReadOrError(r, data[1:]); err != nil {
		return uint64(parsedBytes), err
	} else {
		m.Timestamp = binary.BigEndian.Uint32(data)
		parsedBytes += 3
	}

	if fmt <= MessageHeaderFmt1 {
		if err := util.ReadOrError(r, data[1:]); err != nil {
			return uint64(parsedBytes), err
		} else {
			length := binary.BigEndian.Uint32(data)
			m.Length = &length
			parsedBytes += 3
		}

		if nextByte, err := util.ReadByteOrError(r); err != nil {
			return uint64(parsedBytes), err
		} else {
			m.TypeID = &nextByte
			parsedBytes += 1
		}
	}

	if fmt == MessageHeaderFmt0 {
		if err := util.ReadOrError(r, data); err != nil {
			return uint64(parsedBytes), err
		} else {
			streamID := binary.BigEndian.Uint32(data)
			m.StreamID = &streamID
			parsedBytes += 4
		}
	}

	return uint64(parsedBytes), nil
}

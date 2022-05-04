package chunk

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// Message represents RTMP chunk message.
type Message struct {
	BasicHeader       BasicHeader    `json:"basic_header"`
	MessageHeader     *MessageHeader `json:"message_header,omitempty"`
	ExtendedTimestamp *uint32        `json:"extended_timestamp,omitempty"`

	// TODO: data
}

// Serialize serializes chunk message to binary format.
func (m *Message) Serialize() []byte {
	data := []byte{}

	bh := m.BasicHeader.Serialize()
	data = append(data, bh...)

	if m.MessageHeader != nil {
		mh := m.MessageHeader.Serialize(m.BasicHeader.Fmt)
		data = append(data, mh...)
	}
	if m.ExtendedTimestamp != nil {
		timestampData := make([]byte, 4)
		binary.BigEndian.PutUint32(timestampData, *m.ExtendedTimestamp)
		data = append(data, timestampData[1:]...)
	}

	// TODO: data

	return data
}

// ParsePayload parses chunk message payload from binary format.
// It assumes the BasicHeader has been parsed already.
func (m *Message) ParsePayload(r io.Reader) error {
	var parsedBytes uint64

	if m.BasicHeader.Fmt < MessageHeaderFmt3 {
		m.MessageHeader = &MessageHeader{}
		if bytes, err := m.MessageHeader.Parse(r, m.BasicHeader.Fmt); err != nil {
			return err
		} else {
			parsedBytes += bytes
		}
	}

	if m.MessageHeader != nil &&
		m.MessageHeader.Timestamp == 0xFFFFFF {
		data := make([]byte, 4)
		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			timestamp := binary.BigEndian.Uint32(data)
			m.ExtendedTimestamp = &timestamp
			parsedBytes += 4
		}
	}

	// TODO: data
	glog.Warningf("parsed %d bytes, still need to parse data", parsedBytes)

	return nil
}

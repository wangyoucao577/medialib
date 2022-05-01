package sei

import (
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// UserDataUnregistered represents AVC SEI UserDataUnregistered.
type UserDataUnregistered struct {
	UUID                []byte `json:"uuid"` // uuid_iso_iec_11578, fixed 16 bytes
	UserDataPayloadByte []byte `json:"user_data_payload_byte"`
}

func (u *UserDataUnregistered) parse(r io.Reader, payloadSize int) (uint64, error) {
	var parsedBytes uint64

	data := make([]byte, 16)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		u.UUID = data
		parsedBytes += 16
	}

	u.UserDataPayloadByte = make([]byte, uint64(payloadSize)-parsedBytes)
	if err := util.ReadOrError(r, u.UserDataPayloadByte); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += uint64(len(u.UserDataPayloadByte))
	}

	return parsedBytes, nil
}

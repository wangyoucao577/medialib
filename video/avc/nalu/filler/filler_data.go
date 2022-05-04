// Package filler represents Filler Data NAL Units.
package filler

import (
	"encoding/json"
	"io"

	"github.com/wangyoucao577/medialib/util"
)

// Data represents Filler data.
type Data struct {
	size int `json:"-"` // filler data size, not in byte stream but calculated
}

// MarshalJSON implements json.Marshaler.
func (d *Data) MarshalJSON() ([]byte, error) {
	var dj = struct {
		Size int `json:"size"`
	}{
		Size: d.size,
	}

	return json.Marshal(dj)
}

// Parse parses bytes to AVC NAL Unit, return parsed bytes or error.
// The NAL Unit syntax defined in ISO/IEC-14496-10 7.3.1.
func (d *Data) Parse(r io.Reader, size int) (uint64, error) {
	d.size = size

	// read to ignore
	if err := util.ReadOrError(r, make([]byte, size)); err != nil {
		return 0, err
	}

	return uint64(size), nil
}

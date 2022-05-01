// Package aud represents AVC access unit delimiter NAL unit, defined in ISO/IEC-14496-10 7.3.2.4.
package aud

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// AccessUnitDelimiter represents AccessUnitDelimiter NALU structure.
type AccessUnitDelimiter struct {
	PrimaryPicType uint8 `json:"primary_pic_type"` // 3 bits
}

// Parse parses bytes to AVC AccessUnitDelimiter NAL Unit, return parsed bytes or error.
func (a *AccessUnitDelimiter) Parse(r io.Reader, size int) (uint64, error) {
	var parsedBytes uint64

	// parse nalu length
	data := make([]byte, 1)
	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		a.PrimaryPicType = (data[0] >> 5) & 0x7
		parsedBytes += 1
	}

	// ignore rbsp_trailing_bits()
	glog.V(3).Infof("nalu AccessUnitDelimiter ignore %d bits rbsp_trailing_bits", 5+(size-int(parsedBytes))*8)

	return parsedBytes, nil
}

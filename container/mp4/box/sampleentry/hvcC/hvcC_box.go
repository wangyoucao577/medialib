// Package hvcc reprensents hvcC, i.e., HEVC Configraiton box.
package hvcc

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// HEVCConfigrationBox defines HEVC Configraiton box.
type HEVCConfigrationBox struct {
	box.Header `json:"header"`

	HEVCConfig HEVCDecoderConfigurationRecord `json:"hevc_config"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &HEVCConfigrationBox{
		Header: h,
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (h *HEVCConfigrationBox) ParsePayload(r io.Reader) error {
	if err := h.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", h.Type, err)
		return nil
	}

	var parsedBytes uint64
	if bytes, err := h.HEVCConfig.Parse(r); err != nil {
		return err
	} else {
		parsedBytes += bytes
	}

	if parsedBytes != h.PayloadSize() {
		if parsedBytes > h.PayloadSize() { // parse wrong
			return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", h.Type, parsedBytes, h.PayloadSize())
		} else {
			//TODO: these can be removed if we have supported full AVCDecoderConfigurationRecord parsing
			remainBytes := h.PayloadSize() - parsedBytes
			glog.Warningf("sample entry box type %s still has %d bytes hasn't been parsed yet, ignore them", h.Type, remainBytes)

			// ignore pre_defined 2 bytes in here
			if err := util.ReadOrError(r, make([]byte, remainBytes)); err != nil {
				return err
			}
		}
	}

	return nil
}

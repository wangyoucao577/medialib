// Package avcc reprensents avcC, i.e., AVC Configraiton box.
package avcc

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// AVCConfigrationBox defines AVC Configraiton box.
type AVCConfigrationBox struct {
	box.Header `json:"header"`

	AVCConfig AVCDecoderConfigurationRecord `json:"avc_config"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &AVCConfigrationBox{
		Header: h,
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (a *AVCConfigrationBox) ParsePayload(r io.Reader) error {
	if err := a.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", a.Type, err)
		return nil
	}

	var parsedBytes uint64
	if bytes, err := a.AVCConfig.Parse(r); err != nil {
		return err
	} else {
		parsedBytes += bytes
	}

	if parsedBytes != a.PayloadSize() {
		if parsedBytes > a.PayloadSize() { // parse wrong
			return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", a.Type, parsedBytes, a.PayloadSize())
		} else {
			//TODO: these can be removed if we have supported full AVCDecoderConfigurationRecord parsing
			remainBytes := a.PayloadSize() - parsedBytes
			glog.Warningf("sample entry box type %s still has %d bytes hasn't been parsed yet, ignore them", a.Type, remainBytes)

			// ignore pre_defined 2 bytes in here
			if err := util.ReadOrError(r, make([]byte, remainBytes)); err != nil {
				return err
			}
		}
	}

	return nil
}

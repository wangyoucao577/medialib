package sei

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// PicTiming represents AVC SEI pic_timing.
type PicTiming struct {
	CpbRemovalDelay *uint8 `json:"cpb_removal_delay,omitempty"`
	DpbOutputDelay  *uint8 `json:"dpb_output_delay,omitempty"`
	PicStruct       *uint8 `json:"pic_struct,omitempty"`

	// store for some internal parsing
	sps *sps.SequenceParameterSetData `json:"-"`
	pps *pps.PictureParameterSet      `json:"-"`
}

// setSequenceHeaders sets both SPS and PPS for parsing.
func (p *PicTiming) setSequenceHeaders(sps *sps.SequenceParameterSetData, pps *pps.PictureParameterSet) {
	p.sps = sps
	p.pps = pps
}

func (p *PicTiming) parse(r io.Reader, payloadSize int) (uint64, error) {
	var parsedBytes uint64

	if p.sps == nil || p.sps.VuiParametersPresentFlag == 0 || p.sps.VUIParameters == nil {
		return parsedBytes, fmt.Errorf("invalid sps %v", p.sps)
	}

	// TODO:cpb, dpb parsing
	// var cpbDpbDelaysPresentFlag bool
	// if p.sps.VUIParameters.NalHrdParametersPresentFlag == 1 || p.sps.VUIParameters.VclHrdParametersPresentFlag == 1 {
	// 	cpbDpbDelaysPresentFlag = true
	// }
	// if cpbDpbDelaysPresentFlag {

	// }

	glog.V(3).Infof("sei pic_timing payload bytes %d parsing TODO", payloadSize)
	//TODO: parse payload
	if err := util.ReadOrError(r, make([]byte, payloadSize)); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += uint64(payloadSize)
	}

	return parsedBytes, nil
}

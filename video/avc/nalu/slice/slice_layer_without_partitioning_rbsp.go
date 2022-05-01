package slice

import (
	"io"

	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/video/avc/nalu/pps"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// LayerWithoutPartitioningRbsp represents slice_layer_without_partitioning_rbsp defined in ISO/IEC-14496-10 7.3.2.8.
type LayerWithoutPartitioningRbsp struct {
	Header Header `json:"slice_header"`
	// Data   Data   `json:"slice_data"` // TODO:

	sps *sps.SequenceParameterSetData `json:"-"`
	pps *pps.PictureParameterSet      `json:"-"`
}

// SetSequenceHeaders sets both SPS and PPS for parsing.
func (l *LayerWithoutPartitioningRbsp) SetSequenceHeaders(sps *sps.SequenceParameterSetData, pps *pps.PictureParameterSet) {
	l.sps = sps
	l.pps = pps
}

const bitsPerByte = 8

// Parse parses bytes to AVC SliceLayerWithoutPartitioningRbsp NAL Unit, return parsed bytes or error.
func (l *LayerWithoutPartitioningRbsp) Parse(r io.Reader, size int) (uint64, error) {
	var parsedBits uint64
	br := bitreader.New(r) // start bit-level parsing here

	if costBits, err := l.parseHeader(br); err != nil {
		return parsedBits / bitsPerByte, err
	} else {
		parsedBits += costBits
	}

	//TODO: slice data

	return parsedBits / bitsPerByte, nil
}

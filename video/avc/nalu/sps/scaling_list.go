package sps

import (
	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util/bitreader"
	"github.com/wangyoucao577/medialib/util/expgolombcoding"
)

const bitsPerByte = 8

type scalingListParser struct {
	sizeOfScalingList int

	scalingList []int                    // result
	deltaScale  []expgolombcoding.Signed // parsed data
}

// return parsed bits
func (s *scalingListParser) parse(r *bitreader.Reader) (uint64, error) {
	var parsedBits uint64

	lastScale := 8
	nextScale := 8

	s.scalingList = make([]int, s.sizeOfScalingList)
	for j := 0; j < s.sizeOfScalingList; j++ {

		if nextScale != 0 {
			deltaScale := expgolombcoding.Signed{}
			if costBits, err := deltaScale.Parse(r); err != nil {
				return parsedBits, err
			} else {
				parsedBits += costBits
			}
			s.deltaScale = append(s.deltaScale, deltaScale)

			nextScale = (lastScale + int(deltaScale.Value()) + 256) % 256

			//TODO: use default flags: useDefaultScalingMatrixFlag = ( j = = 0 && nextScale = = 0 )
			glog.V(3).Infof("TODO useDefaultScalingMatrixFlag")
		}

		if nextScale == 0 {
			s.scalingList[j] = lastScale
		} else {
			s.scalingList[j] = nextScale
		}
		lastScale = s.scalingList[j]
	}

	return 0, nil
}

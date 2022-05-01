package avcc

import (
	"encoding/binary"
	"io"

	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
	"github.com/wangyoucao577/medialib/video/avc/nalu/sps"
)

// LengthParameterSetNALU represents Length,AVC SPS/PPS/SPSExt NALU, Length, ... composition.
type LengthParameterSetNALU struct {
	Length  uint16       `json:"length"`
	NALUnit nalu.NALUnit `json:"nal_unit"`
}

// AVCDecoderConfigurationRecord defines AVC Decoder configuration record.
type AVCDecoderConfigurationRecord struct {
	ConfigurationVersion     uint8  `json:"configuration_version"`
	AVCProfileIndication     uint8  `json:"avc_profile_indication"`
	AVCProfileIndicationName string `json:"avc_profile_indication_name"` // NOT in byte stream, only store for better intuitive
	ProfileCompatibility     uint8  `json:"profile_compatibility"`
	AVCLevelIndication       uint8  `json:"avc_level_indication"`

	// 6 bits reserved here
	LengthSizeMinusOne uint8 `json:"length_size_minus_one"` // 2 bits in file

	// sps
	// 3 bits reserved here
	NumOfSequenceParameterSets uint8                    `json:"num_of_sequence_parameter_sets"` // 5 bits in file
	LengthSPSNALU              []LengthParameterSetNALU `json:"sequence_parameter_set,omitempty"`

	// pps
	NumOfPictureParameterSets uint8                    `json:"num_of_picture_parameter_sets"`
	LengthPPSNALU             []LengthParameterSetNALU `json:"picture_parameter_set,omitempty"`

	// 6 bits reserved here
	ChromaFormat uint8 `json:"chroma_format"` // 2 bits in file
	// 5 bits reserved here
	BitDepthLumaMinus8 uint8 `json:"bit_depth_luma_minus8"` // 3 bits in file
	// 5 bits reserved here
	BitDepthChromaMinus8 uint8 `json:"bit_depth_chroma_minus8"` // 3 bits in file

	// sps extensions
	NumOfSequenceParameterSetExt uint8                    `json:"num_of_sequence_parameter_set_ext"`
	LengthSPSExtNALU             []LengthParameterSetNALU `json:"sequence_parameter_set_ext,omitempty"`
}

// LengthSize returns NALU prefix length size.
func (a *AVCDecoderConfigurationRecord) LengthSize() uint32 {
	return uint32(a.LengthSizeMinusOne) + 1
}

// Parse parses AVCDecoderConfigurationRecord.
func (a *AVCDecoderConfigurationRecord) Parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64

	data := make([]byte, 4)

	if err := util.ReadOrError(r, data); err != nil {
		return parsedBytes, err
	} else {
		a.ConfigurationVersion = data[0]
		a.AVCProfileIndication = data[1]
		a.AVCProfileIndicationName = sps.ProfileName(a.AVCProfileIndication)
		a.ProfileCompatibility = data[2]
		a.AVCLevelIndication = data[3]

		parsedBytes += 4
	}

	if err := util.ReadOrError(r, data[:2]); err != nil {
		return parsedBytes, err
	} else {
		a.LengthSizeMinusOne = data[0] & 0x3
		a.NumOfSequenceParameterSets = data[1] & 0x1F
		parsedBytes += 2
	}

	a.LengthSPSNALU = make([]LengthParameterSetNALU, a.NumOfSequenceParameterSets)
	for i := 0; i < int(a.NumOfSequenceParameterSets); i++ {
		var len uint16
		if err := util.ReadOrError(r, data[:2]); err != nil {
			return parsedBytes, err
		} else {
			len = binary.BigEndian.Uint16(data[:2])
			a.LengthSPSNALU[i].Length = len
			parsedBytes += 2
		}

		if bytes, err := a.LengthSPSNALU[i].NALUnit.Parse(r, int(len)); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += uint64(bytes)
		}
	}

	if err := util.ReadOrError(r, data[:1]); err != nil {
		return parsedBytes, err
	} else {
		a.NumOfPictureParameterSets = data[0]
		parsedBytes += 1
	}

	a.LengthPPSNALU = make([]LengthParameterSetNALU, a.NumOfPictureParameterSets)
	for i := 0; i < int(a.NumOfPictureParameterSets); i++ {
		var len uint16
		if err := util.ReadOrError(r, data[:2]); err != nil {
			return parsedBytes, err
		} else {
			len = binary.BigEndian.Uint16(data[:2])
			a.LengthPPSNALU[i].Length = len
			parsedBytes += 2
		}

		if bytes, err := a.LengthPPSNALU[i].NALUnit.Parse(r, int(len)); err != nil {
			return parsedBytes, err
		} else {
			parsedBytes += uint64(bytes)
		}
	}

	// ISO/IEC-14496-15 5.2.4.1
	// But seems not work, try it out if need
	// if a.AVCProfileIndication == 100 || a.AVCProfileIndication == 110 ||
	// 	a.AVCProfileIndication == 122 || a.AVCProfileIndication == 144 {

	// 	if err := util.ReadOrError(r, data); err != nil {
	// 		return parsedBytes, err
	// 	} else {
	// 		a.ChromaFormat = data[0] & 0x3
	// 		a.BitDepthLumaMinus8 = data[1] & 0x7
	// 		a.BitDepthChromaMinus8 = data[2] & 0x7
	// 		a.NumOfSequenceParameterSetExt = data[3]

	// 		parsedBytes += 4
	// 	}

	// 	a.LengthSPSExtNALU = make([]LengthParameterSetNALU, a.NumOfSequenceParameterSetExt)
	// 	for i := 0; i < int(a.NumOfSequenceParameterSetExt); i++ {
	// 		var len uint16
	// 		if err := util.ReadOrError(r, data[:2]); err != nil {
	// 			return parsedBytes, err
	// 		} else {
	// 			len = binary.BigEndian.Uint16(data[:2])
	// 			a.LengthSPSExtNALU[i].Length = len
	// 			parsedBytes += 2
	// 		}

	// 		if bytes, err := a.LengthSPSExtNALU[i].NALUnit.Parse(r, int(len)); err != nil {
	// 			return parsedBytes, err
	// 		} else {
	// 			parsedBytes += uint64(bytes)
	// 		}
	// 	}
	// }

	return parsedBytes, nil
}

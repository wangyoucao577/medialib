package esds

import (
	"encoding/binary"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// ESDescriptor represents ES_Descriptor.
type ESDescriptor struct {
	Descriptor Descriptor `json:"descriptor"`

	ESID                    uint16                  `json:"es_id"`
	StreamDependenceFlag    uint8                   `json:"stream_dependence_flag"` // 1 bit
	URLFlag                 uint8                   `json:"url_flag"`               // 1 bit
	OCRStreamFlag           uint8                   `json:"ocr_stream_flag"`        // 1 bit
	StreamPriority          uint8                   `json:"stream_priority"`        // 5 bits
	DependsOnESID           uint16                  `json:"depends_on_es_id,omitempty"`
	URLLength               uint8                   `json:"url_length,omitempty"`
	URLstring               string                  `json:"url_string,omitempty"` // len(URLstring) == URLLength
	OCR_ES_Id               uint16                  `json:"ocr_es_id,omitempty"`
	DecoderConfigDescriptor DecoderConfigDescriptor `json:"decoder_config_descriptor"`
	SLConfigDescriptor      *SLConfigDescriptor     `json:"sl_config_descriptor,omitempty"`
}

func (e *ESDescriptor) parse(r io.Reader) (uint64, error) {

	var parsedBytes uint64
	var parsedHeaderBytes uint64

	data := make([]byte, 4)

	// parse descriptor header
	if bytes, err := e.Descriptor.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
		parsedHeaderBytes += bytes
	}

	if err := util.ReadOrError(r, data[:3]); err != nil {
		return parsedBytes, err
	} else {
		e.ESID = binary.BigEndian.Uint16(data[:2])
		e.StreamDependenceFlag = (data[2] >> 7) & 0x1
		e.URLFlag = (data[2] >> 6) & 0x1
		e.OCRStreamFlag = (data[2] >> 5) & 0x1
		e.StreamPriority = data[2] & 0x1F
		parsedBytes += 3
	}

	if e.StreamDependenceFlag == 0x1 {
		if err := util.ReadOrError(r, data[:2]); err != nil {
			return parsedBytes, err
		} else {
			e.DependsOnESID = binary.BigEndian.Uint16(data[:2])
			parsedBytes += 2
		}
	}

	if e.URLFlag == 0x1 {
		if err := util.ReadOrError(r, data[:1]); err != nil {
			return parsedBytes, err
		} else {
			e.URLLength = uint8(data[0])
			parsedBytes += 1
		}

		if e.URLLength > 0 {
			u := make([]byte, e.URLLength)
			if err := util.ReadOrError(r, u); err != nil {
				return parsedBytes, err
			} else {
				e.URLstring = string(u)
				parsedBytes += uint64(e.URLLength)
			}
		}
	}

	if e.OCRStreamFlag == 0x1 {
		if err := util.ReadOrError(r, data[:2]); err != nil {
			return parsedBytes, err
		} else {
			e.OCR_ES_Id = binary.BigEndian.Uint16(data[:2])
			parsedBytes += 2
		}
	}

	if bytes, err := e.DecoderConfigDescriptor.parse(r); err != nil {
		return parsedBytes, err
	} else {
		parsedBytes += bytes
	}

	if uint64(e.Descriptor.Size) > parsedBytes-parsedHeaderBytes {
		slDesc := SLConfigDescriptor{}
		if bytes, err := slDesc.parse(r); err != nil {
			return parsedBytes, err
		} else {
			e.SLConfigDescriptor = &slDesc
			parsedBytes += bytes
		}
	}

	if uint64(e.Descriptor.Size) != parsedBytes-parsedHeaderBytes {
		glog.Warningf("descriptor %s(%d) still has %d bytes need to parse, parsed payload bytes != payload size: %d != %d", classTagName(e.Descriptor.Tag), e.Descriptor.Tag, uint64(e.Descriptor.Size)-uint64(parsedBytes-parsedHeaderBytes), parsedBytes-parsedHeaderBytes, e.Descriptor.Size)
	}

	return parsedBytes, nil
}

// Package flv reprensents The FLV File Format, defines in Adobe Flash Video File Format Specification Version 10.1 AnnexE.
package flv

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv/tag"
	"github.com/wangyoucao577/medialib/flv/tag/audio"
	"github.com/wangyoucao577/medialib/flv/tag/script"
	"github.com/wangyoucao577/medialib/flv/tag/video"
	"github.com/wangyoucao577/medialib/util"
)

// FLV represents the FLV file format.
type FLV struct {
	Header Header `json:"header"`

	// PreviousTagSize[0] always == 0
	// PreviousTagSize[i+1] = len(Tags[i])
	PreviousTagSize []uint32  `json:"PreviousTagSize"`
	Tags            []tag.Tag `json:"Tags"`
}

// Parse parses FLV data.
func (f *FLV) Parse(r io.Reader) error {
	if err := f.Header.Parse(r); err != nil {
		return err
	}

	tagSizeData := make([]byte, 4) // fixed 4 bytes
	for {
		// parse previous tag size
		if err := util.ReadOrError(r, tagSizeData); err != nil {
			return err
		}
		size := binary.BigEndian.Uint32(tagSizeData)
		f.PreviousTagSize = append(f.PreviousTagSize, size)

		// parse tag header
		tagHeader := tag.Header{}
		if err := tagHeader.Parse(r); err != nil {
			return err
		}

		// parse tag data
		var t tag.Tag
		if tagHeader.TagType == tag.TypeAudio {
			t = &audio.Tag{Header: tagHeader}
		} else if tagHeader.TagType == tag.TypeVideo {
			t = &video.Tag{Header: tagHeader}
		} else if tagHeader.TagType == tag.TypeSriptData {
			t = &script.Tag{Header: tagHeader}
		}
		if tagHeader.Filter == 1 {

			//TODO: parse payload
			glog.Warningf("unsupported tag type %d with filter=1, ignore payload size %d", tagHeader.TagType, tagHeader.DataSize)
			if err := util.ReadOrError(r, make([]byte, tagHeader.DataSize)); err != nil {
				return err
			}
			continue
		}

		if err := t.ParsePayload(r); err != nil {
			return err
		}
		f.Tags = append(f.Tags, t)
	}
}

// JSON marshals FLV to JSON representation
func (f FLV) JSON() ([]byte, error) {
	return json.Marshal(f)
}

// JSONIndent marshals FLV to JSON representation with customized indent.
func (f FLV) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(f, prefix, indent)
}

// YAML formats FLV to YAML representation.
func (f FLV) YAML() ([]byte, error) {
	j, err := json.Marshal(f)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CSV formats FLV to CSV representation, which isn't supported at the moment.
func (f FLV) CSV() ([]byte, error) {
	return nil, fmt.Errorf("csv representation does not support yet")
}

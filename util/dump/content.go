package dump

import "fmt"

// ContentType represents content type to ouput.
type ContentType string

// Contents to output
const (
	// after parse
	ContentTypeBoxes       = "boxes"         // mp4/fmp4 boxes parsing data
	ContentTypeTags        = "tags"          // FLV header and tags parsing data
	ContentTypeES          = "es"            // AVC/HEVC Elementary Stream Parsing data
	ContentTypeRawES       = "raw_es"        // AVC/HEVC Elementary Stream Raw data (mp4 video elementary stream only, no sps/pps)
	ContentTypeRawAnnexBES = "raw_es_annexb" // AVC/HEVC Elementary Stream Raw data (AnnexB byte format, video elementary stream and parameter set elementary stream)

	// no parse needed
	ContentTypeBoxTypes  = "box_types"  // Supported boxes
	ContentTypeNALUTypes = "nalu_types" // NALU types
)

var conentDescriptions = map[ContentType]string{
	ContentTypeBoxes:       "parsed mp4/fmp4 boxes",
	ContentTypeTags:        "parsed flv header and tags",
	ContentTypeES:          "parsed avc/hevc elementary stream",
	ContentTypeRawES:       "extracted raw data of avc/hevc elementary stream, mp4 video elementary stream only, no sps/pps",
	ContentTypeRawAnnexBES: "extracted raw data of avc/hevc elementary stream described by AnnexB byte format, including video elementary stream and parameter set elementary stream",

	ContentTypeBoxTypes:  "supported box types",
	ContentTypeNALUTypes: "supported nal unit types",
}

// String implements Stringer.
func (c ContentType) String() string {
	return string(c)
}

// Description reprensents description of content type.
func (c ContentType) Description() string {
	d, ok := conentDescriptions[c]
	if !ok {
		return "unknown"
	}
	return d
}

// FixedLenString returns fixed-length string for better shown in helper, current max 16 bytes, append space at the end.
func (c ContentType) FixedLenString(length int) string {
	const maxFixedLen = 16
	if length > maxFixedLen {
		length = maxFixedLen
	}

	if len(c) >= length {
		return c.String()
	}

	s := c.String()
	for i := len(c); i < length; i++ {
		s += " "
	}
	return s
}

// GetContentType returns ContentType if valid, otherwise error
func GetConentType(s string) (ContentType, error) {
	_, ok := conentDescriptions[ContentType(s)]
	if !ok {
		return "", fmt.Errorf("invalid content type %s", s)
	}
	return ContentType(s), nil
}

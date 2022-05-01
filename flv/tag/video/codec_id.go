package video

// Codec IDs
const (
	CodecIDSorensonH263           = 2
	CodecIDScreenVideo            = 3
	CodecIDOn2VP6                 = 4
	CodecIDOn2VP6WithAlphaChannel = 5
	CodecIDScreenVideoVersion2    = 6
	CodecIDAVC                    = 7
)

var codecIDDescriptions = map[int]string{
	CodecIDSorensonH263:           "Sorenson H.263",
	CodecIDScreenVideo:            "Screen video",
	CodecIDOn2VP6:                 "On2 VP6",
	CodecIDOn2VP6WithAlphaChannel: "On2 VP6 with alpha channel",
	CodecIDScreenVideoVersion2:    "Screen video version 2",
	CodecIDAVC:                    "AVC",
}

// CodecIDDescription returns description of codec ID.
func CodecIDDescription(t int) string {
	d, ok := codecIDDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

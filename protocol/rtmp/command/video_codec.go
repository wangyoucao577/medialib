package command

// Values for the video codecs property
const (
	VideoCodecSupportVIDUnused    = 0x0001
	VideoCodecSupportVIDJPEG      = 0x0002
	VideoCodecSupportVIDSorenson  = 0x0004
	VideoCodecSupportVIDHomebrew  = 0x0008
	VideoCodecSupportVIDVP6       = 0x0010
	VideoCodecSupportVIDVP6Alpha  = 0x0020
	VideoCodecSupportVIDHomebrewV = 0x0040
	VideoCodecSupportVIDH264      = 0x0080
	VideoCodecSupportVIDAll       = 0x00FF
)

var videoCodecDescriptions = map[int]string{
	VideoCodecSupportVIDUnused:    "Obsolete value",
	VideoCodecSupportVIDJPEG:      "Obsolete value",
	VideoCodecSupportVIDSorenson:  "Sorenson Flash video",
	VideoCodecSupportVIDHomebrew:  "V1 screen sharing",
	VideoCodecSupportVIDVP6:       "On2 video (Flash 8+)",
	VideoCodecSupportVIDVP6Alpha:  "On2 video with alpha channel",
	VideoCodecSupportVIDHomebrewV: "Screen sharing version 2(Flash 8+)",
	VideoCodecSupportVIDH264:      "H264 video",
	VideoCodecSupportVIDAll:       "All RTMP-supported video codecs",
}

// VideoCodecDescription returns description of supported video codec.
func VideoCodecDescription(t int) string {
	d, ok := audioCodecDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

// IsVideoCodecValid checks whether it's a valid video codec.
func IsVideoCodecValid(t int) bool {
	_, ok := videoCodecDescriptions[t]
	return ok
}

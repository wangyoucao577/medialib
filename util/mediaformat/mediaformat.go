package mediaformat

// Media Formats
const (
	MP4  = "mp4"
	M4S  = "m4s"
	FMP4 = "fmp4"
	MOV  = "mov"

	FLV = "flv"

	H264 = "h264"
)

// AsExtension returns extension representation of the format, e.g. return '.mp4' for format 'mp4'.
func AsExtension(format string) string {
	return "." + format
}

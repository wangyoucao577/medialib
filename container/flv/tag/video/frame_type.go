package video

// Frame types
const (
	FrameTypeKey                = 1
	FrameTypeInner              = 2
	FrameTypeDisposableInner    = 3
	FrameTypeGeneratedKey       = 4
	FrameTypeVideoInfoOrCommand = 5
)

var frameTypeDescriptions = map[int]string{
	FrameTypeKey:                "key frame (for AVC, a seekable frame)",
	FrameTypeInner:              "inter frame (for AVC, a non-seekable frame)",
	FrameTypeDisposableInner:    "disposable inter frame (H.263 only)",
	FrameTypeGeneratedKey:       "generated key frame (reserved for server use only)",
	FrameTypeVideoInfoOrCommand: "video info/command frame",
}

// FrameTypeDescription returns description of frame type.
func FrameTypeDescription(t int) string {
	d, ok := frameTypeDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

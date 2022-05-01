package tag

// FLV tag types
const (
	TypeAudio     = 8
	TypeVideo     = 9
	TypeSriptData = 18
)

var typeDescriptions = map[int]string{
	TypeAudio:     "Audio",
	TypeVideo:     "Video",
	TypeSriptData: "Script Data",
}

// TypeDescription returns description of tag type.
func TypeDescription(t int) string {
	d, ok := typeDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

// IsTypeValid checks whether it's a valid tag type.
func IsTypeValid(t int) bool {
	_, ok := typeDescriptions[t]
	return ok
}

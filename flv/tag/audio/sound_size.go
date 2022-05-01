package audio

// Size of each audio sample.
// This parameter only pertains to uncompressed formats.
// Compressed formats always decode to 16 bits internally.
const (
	SoundSize8bit  = 0 // 0 = 8-bit samples
	SoundSize16bit = 1 // 1 = 16-bit samples
)

var soundSizeDescriptions = map[int]string{
	SoundSize8bit:  "8-bit samples",
	SoundSize16bit: "16-bit samples",
}

// SoundSizeDescription returns description of sound size.
func SoundSizeDescription(t int) string {
	d, ok := soundSizeDescriptions[t]
	if !ok {
		return "unknown"
	}
	return d
}

package audio

// Mono or stereo sound
const (
	SoundTypeMono   = 0 // 0 = Mono sound
	SoundTypeStereo = 1 // 1 = Stereo sound
)

var soundTypeDescriptions = map[int]string{
	SoundTypeMono:   "Mono sound",
	SoundTypeStereo: "Stereo sound",
}

// SoundTypeDescription returns description of sound type.
func SoundTypeDescription(t int) string {
	d, ok := soundTypeDescriptions[t]
	if !ok {
		return "unknown"
	}
	return d
}

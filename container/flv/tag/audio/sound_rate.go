package audio

// Sampling rate. The following values are defined:
const (
	SoundRate5p5 = 0 // 0 = 5.5 kHz
	SoundRate11  = 1 // 1 = 11 kHz
	SoundRate22  = 2 // 2 = 22 kHz
	SoundRate44  = 3 // 3 = 44 kHz
)

var soundRateDescriptions = map[int]string{
	SoundRate5p5: "5.5 kHz",
	SoundRate11:  "11 kHz",
	SoundRate22:  "22 kHz",
	SoundRate44:  "44 kHz",
}

// SoundRateDescription returns description of sound rate.
func SoundRateDescription(t int) string {
	d, ok := soundRateDescriptions[t]
	if !ok {
		return "unknown"
	}
	return d
}

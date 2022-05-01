package audio

// Formats of sound data
// Formats 7, 8, 14, and 15 are reserved.
// AAC is supported in Flash Player 9,0,115,0 and higher. Speex is supported in Flash Player 10 and higher.
const (
	SoundFormatLinearPCMPlatformEndian = 0
	SoundFormatADPCM                   = 1
	SoundFormatMP3                     = 2
	SoundFormatLinearPCMLittleEndian   = 3

	SoundFormatG711ALawLogarithmicPCM  = 7
	SoundFormatG711MuLawLogarithmicPCM = 8

	SoundFormatAAC   = 10
	SoundFormatSpeex = 11

	SoundFormatMP38KHz        = 14
	SoundFormatDeviceSpecific = 15
)

var soundFommatDescriptions = map[int]string{
	SoundFormatLinearPCMPlatformEndian: "Linear PCM, platform endian",
	SoundFormatADPCM:                   "ADPCM",
	SoundFormatMP3:                     "MP3",
	SoundFormatLinearPCMLittleEndian:   "Linear PCM, little endian 4 = Nellymoser 16 kHz mono 5 = Nellymoser 8 kHz mono 6 = Nellymoser",

	SoundFormatG711ALawLogarithmicPCM:  "G.711 A-law logarithmic PCM, reserved.",
	SoundFormatG711MuLawLogarithmicPCM: "G.711 mu-law logarithmic PCM, reserved.",

	SoundFormatAAC:   "AAC",
	SoundFormatSpeex: "Speex",

	SoundFormatMP38KHz:        "MP3 8 kHz, reserved.",
	SoundFormatDeviceSpecific: "Device-specific sound, reserved.",
}

// SoundFormatDescription returns description of sound format.
func SoundFormatDescription(t int) string {
	d, ok := soundFommatDescriptions[t]
	if !ok {
		return "unknown"
	}
	return d
}

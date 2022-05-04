package command

// Values for the audio codecs property
const (
	AudioCodecSupportSndNone    = 0x0001
	AudioCodecSupportSndADPCM   = 0x0002
	AudioCodecSupportSndMP3     = 0x0003
	AudioCodecSupportSndIntel   = 0x0008
	AudioCodecSupportSndUnused  = 0x0010
	AudioCodecSupportSndNelly8  = 0x0020
	AudioCodecSupportSndNelly   = 0x0040
	AudioCodecSupportSndG711A   = 0x0080
	AudioCodecSupportSndG711U   = 0x0100
	AudioCodecSupportSndNelly16 = 0x0200
	AudioCodecSupportSndAAC     = 0x0400
	AudioCodecSupportSndSpeex   = 0x0800
	AudioCodecSupportSndAll     = 0x0FFF
)

var audioCodecDescriptions = map[int]string{
	AudioCodecSupportSndNone:    "Raw sound, no compression",
	AudioCodecSupportSndADPCM:   "ADPCM compression",
	AudioCodecSupportSndMP3:     "mp3 compression",
	AudioCodecSupportSndIntel:   "Not used",
	AudioCodecSupportSndUnused:  "Not used",
	AudioCodecSupportSndNelly8:  "NellyMoser at 8-kHz compression",
	AudioCodecSupportSndNelly:   "NellyMoser compression(5, 11, 22, and 44kHz)",
	AudioCodecSupportSndG711A:   "G711A sound compression(Flash Media Server only)",
	AudioCodecSupportSndG711U:   "G711U sound compression(Flash Media Server only)",
	AudioCodecSupportSndNelly16: "NellyMouser at 16-kHz compression",
	AudioCodecSupportSndAAC:     "Advanced audio coding (AAC) codec",
	AudioCodecSupportSndSpeex:   "Speex Audio",
	AudioCodecSupportSndAll:     "All RTMP-supported audio codecs",
}

// AudioCodecDescription returns description of supported audio codec.
func AudioCodecDescription(t int) string {
	d, ok := audioCodecDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

// IsAudioCodecValid checks whether it's a valid audio codec.
func IsAudioCodecValid(t int) bool {
	_, ok := audioCodecDescriptions[t]
	return ok
}

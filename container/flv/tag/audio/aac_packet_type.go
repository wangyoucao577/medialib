package audio

// The following values are defined:
const (
	AACPacketTypeSequenceHeader = 0 // 0 = AAC sequence header
	AACPacketTypeRaw            = 1 // 1 = AAC raw
)

var aacPacketTypeDescriptions = map[int]string{
	AACPacketTypeSequenceHeader: "AAC sequence header",
	AACPacketTypeRaw:            "AAC raw",
}

// AACPacketTypeDescription returns description of AAC Packet Type.
func AACPacketTypeDescription(t int) string {
	d, ok := aacPacketTypeDescriptions[t]
	if !ok {
		return "unknown"
	}
	return d
}

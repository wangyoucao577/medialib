package video

// AVC Packet Types
const (
	AVCPacketTypeSequenceHeader = 0
	AVCPacketTypeNALU           = 1
	AVCPacketTypeEOS            = 2
)

var avcPacketTypeDescriptions = map[int]string{
	AVCPacketTypeSequenceHeader: "AVC sequence header",
	AVCPacketTypeNALU:           "AVC NALU",
	AVCPacketTypeEOS:            "AVC end of sequence (lower level NALU sequence ender is not required or supported)",
}

// AVCPacketTypeDescription returns description of frame AVC Packet Type.
func AVCPacketTypeDescription(t int) string {
	d, ok := avcPacketTypeDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

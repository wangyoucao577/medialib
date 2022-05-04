package chunk

const (
	MessageTypeIDSetPacketSize     = 0x01
	MessageTypeIDAbort             = 0x02
	MessageTypeIDAcknowledge       = 0x03
	MessageTypeIDControl           = 0x04
	MessageTypeIDServerBandwidth   = 0x05
	MessageTypeIDClientBandwidth   = 0x06
	MessageTypeIDVirtualControl    = 0x07
	MessageTypeIDAudioPacket       = 0x08
	MessageTypeIDVideoPacket       = 0x09
	MessageTypeIDDataExtended      = 0x0F
	MessageTypeIDContainerExtended = 0x10
	MessageTypeIDCommandExtended   = 0x11
	MessageTypeIDData              = 0x12
	MessageTypeIDContainer         = 0x13
	MessageTypeIDCommand           = 0x14
	MessageTypeIDUDP               = 0x15
	MessageTypeIDAggregate         = 0x16
	MessageTypeIDPresent           = 0x17
)

var messageTypeIDDescriptions = map[int]string{
	MessageTypeIDSetPacketSize:     "Set Packet Size Message",
	MessageTypeIDAbort:             "Abort",
	MessageTypeIDAcknowledge:       "Acknowledge",
	MessageTypeIDControl:           "Control Message",
	MessageTypeIDServerBandwidth:   "Server Bandwidth",
	MessageTypeIDClientBandwidth:   "Client Bandwidth",
	MessageTypeIDVirtualControl:    "Virtual Control",
	MessageTypeIDAudioPacket:       "Audio Packet",
	MessageTypeIDVideoPacket:       "Video Packet",
	MessageTypeIDDataExtended:      "Data Extended",
	MessageTypeIDContainerExtended: "Container Extended",
	MessageTypeIDCommandExtended:   "Command Extended (An AMF3 type command)",
	MessageTypeIDData:              "Data (Invoke (onMetaData info is sent as such))",
	MessageTypeIDContainer:         "Container",
	MessageTypeIDCommand:           "Command (An AMF0 type command)",
	MessageTypeIDUDP:               "UDP",
	MessageTypeIDAggregate:         "Aggregate",
	MessageTypeIDPresent:           "Present",
}

// MessageTypeIDDescription returns description of message type.
func MessageTypeIDDescription(t int) string {
	d, ok := messageTypeIDDescriptions[t]
	if !ok {
		return ""
	}
	return d
}

// IsMessageTypeIDValid checks whether it's a valid message type.
func IsMessageTypeIDValid(t int) bool {
	_, ok := messageTypeIDDescriptions[t]
	return ok
}

package sei

// SEI payload types
const (
	PayloadTypeBufferingPeriod      = 0
	PayloadTypePicTiming            = 1
	PayloadTypeUserDataUnregistered = 5

	//TODO: other types
)

var payloadTypeDescriptions = map[int]string{
	PayloadTypeBufferingPeriod:      "buffering_period",
	PayloadTypePicTiming:            "pic_timing",
	PayloadTypeUserDataUnregistered: "user_data_unregistered",
}

// PayloadTypeDescription represents sei payload type description.
func PayloadTypeDescription(t int) string {
	n, ok := payloadTypeDescriptions[t]
	if !ok {
		return "unknown"
	}
	return n
}

// IsValidPayloadType checks whether input sei payload Type is valid or not.
func IsValidPayloadType(t int) bool {
	_, ok := payloadTypeDescriptions[t]
	return ok
}

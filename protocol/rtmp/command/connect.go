package command

import (
	"github.com/wangyoucao577/medialib/util/amf"
	"github.com/wangyoucao577/medialib/util/amf/amf0"
)

// ConnectCommandObject reprensents Connect Command Object.
type ConnectCommandObject struct {
	App            string `json:"app"`
	Flashver       string `json:"flashver"`
	SwfUrl         string `json:"swfUrl"`
	TcUrl          string `json:"tcUrl"`
	Fpad           bool   `json:"fpad"`
	AudioCodecs    int    `json:"audioCodecs"`
	VideoCodecs    int    `json:"videoCodecs"`
	PageUrl        string `json:"pageUrl"`
	ObjectEncoding int    `json:"objectEncoding"` // AMF encoding method
}

// Connect represents Connect Command Message.
type Connect struct {
	CommandName   string               `json:"CommandName"`
	TranscationID int                  `json:"TranscationID"`
	CommandObject ConnectCommandObject `json:"CommandObject"`
}

// ConnectResult represents server response for connect.
type ConnectResult struct {
	CommandName   string                 `json:"CommandName"`
	TranscationID int                    `json:"TranscationID"`
	Properties    map[string]interface{} `json:"Properties"`
	Information   map[string]interface{} `json:"Information"`
}

// NewConnect creates new Connect command message.
func NewConnect() Connect {
	return Connect{
		CommandName:   "connect",
		TranscationID: 1,
		CommandObject: ConnectCommandObject{
			Flashver:    "LNX 9,0,124,2",
			Fpad:        false,
			AudioCodecs: AudioCodecSupportSndAAC,
			VideoCodecs: VideoCodecSupportVIDH264,
		},
	}
}

// Serialize serializes Connect command to byte stream.
func (c *Connect) Serialize() ([]byte, error) {
	var data []byte

	// command name
	vt := &amf0.ValueType{
		TypeMarker: amf0.TypeMarkerString,
		Value: amf0.StringPayload{
			Length: uint16(len(c.CommandName)),
			Str:    c.CommandName,
		},
	}
	if d, err := vt.Encode(); err != nil {
		return data, err
	} else {
		data = append(data, d...)
	}

	// transcation ID
	vt = &amf0.ValueType{
		TypeMarker: amf0.TypeMarkerNumber,
		Value:      c.TranscationID,
	}
	if d, err := vt.Encode(); err != nil {
		return data, err
	} else {
		data = append(data, d...)
	}

	// command object
	vtObjectProperty := []amf0.ObjectProperty{
		{String: amf0.StringPayload{
			Length: uint16(len("app")),
			Str:    "app",
		},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerString,
				Value: amf0.StringPayload{
					Length: uint16(len(c.CommandObject.App)),
					Str:    c.CommandObject.App,
				},
			},
		},
		{
			String: amf0.StringPayload{
				Length: uint16(len("flashver")),
				Str:    "flashver",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerString,
				Value: amf0.StringPayload{
					Length: uint16(len(c.CommandObject.Flashver)),
					Str:    c.CommandObject.Flashver,
				},
			},
		},
	}

	if len(c.CommandObject.SwfUrl) > 0 {
		vtObjectProperty = append(vtObjectProperty, amf0.ObjectProperty{
			String: amf0.StringPayload{
				Length: uint16(len("swfUrl")),
				Str:    "swfUrl",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerString,
				Value: amf0.StringPayload{
					Length: uint16(len(c.CommandObject.SwfUrl)),
					Str:    c.CommandObject.SwfUrl,
				},
			},
		})
	}

	vtObjectProperty = append(vtObjectProperty,
		amf0.ObjectProperty{
			String: amf0.StringPayload{
				Length: uint16(len("tcUrl")),
				Str:    "tcUrl",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerString,
				Value: amf0.StringPayload{
					Length: uint16(len(c.CommandObject.TcUrl)),
					Str:    c.CommandObject.TcUrl,
				},
			},
		},
		amf0.ObjectProperty{
			String: amf0.StringPayload{
				Length: uint16(len("fpad")),
				Str:    "fpad",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerBoolean,
				Value:      c.CommandObject.Fpad,
			},
		},
		amf0.ObjectProperty{
			String: amf0.StringPayload{
				Length: uint16(len("audioCodecs")),
				Str:    "audioCodecs",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerNumber,
				Value:      AudioCodecSupportSndAAC,
			},
		},
		amf0.ObjectProperty{
			String: amf0.StringPayload{
				Length: uint16(len("videoCodecs")),
				Str:    "videoCodecs",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerNumber,
				Value:      VideoCodecSupportVIDH264,
			},
		},
	)

	if len(c.CommandObject.PageUrl) > 0 {
		vtObjectProperty = append(vtObjectProperty, amf0.ObjectProperty{
			String: amf0.StringPayload{
				Length: uint16(len("pageUrl")),
				Str:    "pageUrl",
			},
			ValueType: amf0.ValueType{
				TypeMarker: amf0.TypeMarkerString,
				Value: amf0.StringPayload{
					Length: uint16(len(c.CommandObject.PageUrl)),
					Str:    c.CommandObject.PageUrl,
				},
			},
		})
	}

	vtObjectProperty = append(vtObjectProperty, amf0.ObjectProperty{
		String: amf0.StringPayload{
			Length: uint16(len("objectEncoding")),
			Str:    "objectEncoding",
		},
		ValueType: amf0.ValueType{
			TypeMarker: amf0.TypeMarkerNumber,
			Value:      amf.Version0,
		},
	})

	vt = &amf0.ValueType{
		TypeMarker: amf0.TypeMarkerObject,
		Value: amf0.ObjectPayload{
			ObjectProperty: vtObjectProperty,
		},
	}

	if d, err := vt.Encode(); err != nil {
		return data, err
	} else {
		data = append(data, d...)
	}

	return data, nil
}

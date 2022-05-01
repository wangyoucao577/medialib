package audio

// AACAudioData reprensents AACAudioData.
type AACAudioData struct {
	AudioSpecificConfig []byte `json:"AudioSpecificConfig,omitempty"`
	RawAACFrameData     []byte `json:"RawAACFrameData"`
}

// TagBody represents audio tag payload.
type TagBody struct {
	AACAudioData *AACAudioData `json:"AACAudioData"`
}

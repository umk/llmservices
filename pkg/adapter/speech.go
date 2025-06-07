package adapter

type SpeechMessage struct {
	Content      string  `json:"content,omitempty"`
	Instructions *string `json:"instructions,omitempty"`
}

type SpeechParams struct {
	Model string   `json:"model"`
	Voice string   `json:"voice"`
	Speed *float64 `json:"speed,omitempty"`
}

type Speech struct {
	Audio []byte `json:"audio"`
}

package adapter

type SpeechMessage struct {
	Content      string  `json:"content,omitempty"`      // The text content to be synthesized to speech
	Instructions *string `json:"instructions,omitempty"` // Optional instructions for speech synthesis
}

type SpeechParams struct {
	Model string   `json:"model"`           // The speech model to be used for synthesis
	Voice string   `json:"voice"`           // The voice identifier to use for speech
	Speed *float64 `json:"speed,omitempty"` // Optional speech rate modifier (1.0 is normal speed)
}

type Speech struct {
	Audio Audio `json:"audio"` // The synthesized audio output
}

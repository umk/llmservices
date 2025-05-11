package openai

import (
	"github.com/openai/openai-go"
)

type Adapter struct{ openai.Client }

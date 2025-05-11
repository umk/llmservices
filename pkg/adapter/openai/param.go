package openai

import "github.com/openai/openai-go/packages/param"

func getOpt[V comparable](value *V) param.Opt[V] {
	if value == nil {
		return param.NullOpt[V]()
	}
	return param.NewOpt(*value)
}

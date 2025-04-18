package fuzzy

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func Decode(bodyBytes []byte, req any) ([]byte, error) {
	RegisterPointerFuzzyDecoders()
	extra.RegisterFuzzyDecoders() // 启用模糊解码

	if err := jsoniter.Unmarshal(bodyBytes, &req); err != nil {
		return nil, err
	}

	fuzzyDecodeBytes, err := jsoniter.Marshal(req)
	if err != nil {
		return nil, err
	}

	return fuzzyDecodeBytes, nil
}

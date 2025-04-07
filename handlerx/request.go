package handlerx

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/zeromicro/go-zero/core/logx"
)

func WeaklyDecodeRequest(r *http.Request, req any) error {
	if r.Body == nil {
		return nil
	}

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		return nil
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err = r.Body.Close(); err != nil {
		return err
	}

	bodyBytes, err = weaklyDecodeRequest(bodyBytes, req)

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	r.ContentLength = int64(len(bodyBytes))

	logx.Debugf("new request body bytes: %s", bodyBytes)

	return nil
}

func weaklyDecodeRequest(bodyBytes []byte, req any) ([]byte, error) {
	RegisterPointerFuzzyDecoders()
	extra.RegisterFuzzyDecoders() // 启用模糊解码

	if err := jsoniter.Unmarshal(bodyBytes, &req); err != nil {
		return nil, err
	}

	weaklyDecodeBytes, err := jsoniter.Marshal(req)
	if err != nil {
		return nil, err
	}

	return weaklyDecodeBytes, nil
}

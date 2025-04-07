package handlerx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
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
	var (
		rawData map[string]any
		err     error
	)
	if err = json.Unmarshal(bodyBytes, &rawData); err != nil {
		return nil, err
	}

	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // 允许弱类型转换
		Result:           &req,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err = decoder.Decode(rawData); err != nil {
		return nil, err
	}

	bodyBytes, err = json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

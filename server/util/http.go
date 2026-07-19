package util

import (
	"io"
	"net/http"
	"project/config"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func HTTPGet(url string) ([]byte, error) {

	rsp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(rsp.Body, config.MaxBodySize))
	if err != nil {
		return nil, err
	}

	return body, nil
}

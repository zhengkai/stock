// Package gold ...
package gold

import (
	"io"
	"net/http"
	"project/config"
	"time"

	jp "github.com/buger/jsonparser"
)

const apiURL = `https://data-asg.goldprice.org/dbXRates/USD`
const referer = `https://goldprice.org/`

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func fetch() (float64, error) {

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set(`Referer`, referer)

	rsp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(rsp.Body, config.MaxBodySize))
	if err != nil {
		return 0, err
	}

	num, err := jp.GetFloat(body, "items", "[0]", "xauPrice")
	if err != nil {
		return 0, err
	}
	return num, nil
}

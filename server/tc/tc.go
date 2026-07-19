// Package tc 腾讯 api
package tc

import (
	"bytes"
	"fmt"
	"regexp"
	"project/util"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const stockURLBase = `https://qt.gtimg.cn/q=`

var regexpTCLine = regexp.MustCompile(`^v_\w+="([^"]+)";$`)

func stockURL(code string) (string, error) {

	region, err := util.StockRegion(code)
	if err != nil {
		return ``, err
	}

	url := stockURLBase + region + code

	ab, err := util.HTTPGet(url)
	if err != nil || len(ab) < 300 {
		return ``, fmt.Errorf(`fetch %s fail: %v`, url, err)
	}

	ab = bytes.TrimSpace(ab)
	u8, err := simplifiedchinese.GBK.NewDecoder().Bytes(ab)
	if err != nil {
		return ``, fmt.Errorf(`decode GBK fail: %v`, err)
	}

	m := regexpTCLine.FindSubmatch(u8)
	if len(m) != 2 {
		return ``, fmt.Errorf(`invalid response: %s`, string(u8))
	}

	return string(m[1]), nil
}

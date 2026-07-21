package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseCents(s string) (int64, error) {
	integer, decimal, ok := strings.Cut(s, ".")
	if !ok || integer == "" || len(decimal) != 2 {
		return 0, fmt.Errorf("金额格式错误: %q", s)
	}

	combo := integer + decimal

	for _, c := range strings.TrimPrefix(combo, "-") {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("金额格式错误: %q", s)
		}
	}

	return strconv.ParseInt(integer+decimal, 10, 64)
}

type QuoteParser struct {
	Err error
}

func (p *QuoteParser) Num(field, v string) *int64 {

	if p.Err != nil {
		return nil
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		p.Err = fmt.Errorf("%s 解析失败: %w", field, err)
		return nil
	}
	return &i
}

func (p *QuoteParser) BP(field, v string) *int32 {

	if p.Err != nil {
		return nil
	}

	i, err := ParseCents(v)
	if err != nil {
		p.Err = fmt.Errorf("%s 解析失败: %w", field, err)
		return nil
	}
	i *= 100
	return new(int32(i))
}

func (p *QuoteParser) Cents(field, v string) *int64 {
	if p.Err != nil {
		return nil
	}

	i, err := ParseCents(v)
	if err != nil {
		p.Err = fmt.Errorf("%s 解析失败: %w", field, err)
		return nil
	}
	return &i
}

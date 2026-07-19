package tc

import (
	"fmt"
	"project/pb"
	"project/util"
	"strings"
)

var ErrInvalidQuote = fmt.Errorf(`invalid quote string`)

func stockQuote(q string) (*pb.Quote, error) {

	ql := strings.Split(q, `~`)
	if len(ql) < 50 {
		return nil, ErrInvalidQuote
	}

	p := &util.QuoteParser{}
	b := pb.Quote_builder{
		Price:    p.Cents(`Price`, ql[3]),
		PreClose: p.Cents(`PreClose`, ql[4]),
		Open:     p.Cents(`Open`, ql[5]),
		High:     p.Cents(`High`, ql[33]),
		Low:      p.Cents(`Low`, ql[34]),

		ChangeBp:       p.BP(`ChangeBp`, ql[31]),
		TurnoverRateBp: p.BP(`TurnoverRateBp`, ql[38]),
		AmplitudeBp:    p.BP(`AmplitudeBp`, ql[43]),

		Volume:   p.Num(`Volume`, ql[36]),
		Turnover: p.Num(`Turnover`, ql[37]),
	}
	if p.Err != nil {
		return nil, p.Err
	}

	return b.Build(), nil
}

func Stock(code string) (*pb.Quote, error) {
	q, err := stockURL(code)
	if err != nil {
		return nil, err
	}

	return stockQuote(q)
}

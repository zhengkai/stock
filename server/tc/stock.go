package tc

import (
	"fmt"
	"project/metrics"
	"project/pb"
	"project/util"
	"strings"
)

var ErrInvalidQuote = fmt.Errorf(`invalid quote string`)

func stockQuote(q string, code string) (*pb.Quote, error) {

	ql := strings.Split(q, `~`)
	if len(ql) < 60 {
		metrics.StockFetchFail(code)
		return nil, ErrInvalidQuote
	}

	// fmt.Println(util.JSON(ql))

	p := &util.QuoteParser{}
	b := pb.Quote_builder{
		Price:    p.Cents(`Price`, ql[3]),
		PreClose: p.Cents(`PreClose`, ql[4]),
		Open:     p.Cents(`Open`, ql[5]),
		High:     p.Cents(`High`, ql[33]),
		Low:      p.Cents(`Low`, ql[34]),
		Change:   p.Cents(`Change`, ql[31]),

		ChangeBp:       p.BP(`ChangeBp`, ql[32]),
		TurnoverRateBp: p.BP(`TurnoverRateBp`, ql[38]),
		AmplitudeBp:    p.BP(`AmplitudeBp`, ql[43]),

		Volume:   p.Num(`Volume`, ql[36]),
		Turnover: p.Num(`Turnover`, ql[37]),

		LimitUp:   p.Cents(`LimitUp`, ql[47]),
		LimitDown: p.Cents(`LimitDown`, ql[48]),

		Ts: new(util.TS()),
	}
	if p.Err != nil {
		fmt.Println(`quote parse error:`, p.Err)
		metrics.StockFetchFail(code)
		return nil, p.Err
	}

	return b.Build(), nil
}

func Stock(code string) (*pb.Quote, error) {
	metrics.StockFetch(code)
	q, err := stockURL(code)
	if err != nil {
		metrics.StockFetchFail(code)
		return nil, err
	}
	return stockQuote(q, code)
}

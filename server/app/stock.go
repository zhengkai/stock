package app

import (
	"project/config"
	"project/metrics"
	"project/pb"
	"project/tc"
	"project/util"
)

func Stock(code string) *pb.Quote {

	metrics.StockQuery(code)

	f := util.NewFileF(`stock/%s.pb`, code)
	if d := stockCache(f); d != nil {
		// fmt.Println(`stock cache hit`, code, d.GetPrice())
		return d
	}
	d, err := tc.Stock(code)
	if err != nil {
		return nil
	}

	f.WriteProto(d)
	return d
}

func stockCache(f *util.File) *pb.Quote {

	d := &pb.Quote{}
	err := f.ReadProto(d)
	if err != nil {
		return nil
	}

	ts := util.TS() - config.StockCacheExpire
	if util.IsWorkTime(false) && ts > d.GetTs() {
		return nil
	}

	return d
}

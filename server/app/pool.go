package app

import (
	"project/pb"
	"sort"
	"sync"
)

type pool struct {
	mux sync.Mutex
	m   map[string]*pb.AppStock
}

func (po *pool) SetConfig(cfg *pb.Config) {
	po.mux.Lock()
	defer po.mux.Unlock()

	for _, v := range cfg.GetAlert() {
		code := v.GetCode()
		c, ok := po.m[code]
		if !ok {
			c = &pb.AppStock{}
			c.SetCode(code)
			po.m[code] = c
		}
		c.SetAlert(v)
	}
}

func (po *pool) setQuote(code string, q *pb.Quote) {
	po.mux.Lock()
	defer po.mux.Unlock()
	v, ok := po.m[code]
	if !ok {
		return
	}
	v.SetQuote(q)
}

func (po *pool) GetPB() *pb.AppPool {
	po.mux.Lock()
	defer po.mux.Unlock()

	r := &pb.AppPool_builder{
		Stock: make([]*pb.AppStock, 0, len(po.m)),
	}
	for _, v := range po.m {
		r.Stock = append(r.Stock, v)
	}

	sort.Slice(r.Stock, func(i, j int) bool {
		return r.Stock[i].GetCode() < r.Stock[j].GetCode()
	})

	return r.Build()
}

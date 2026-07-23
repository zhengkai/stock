package metrics

var (
	stockQuery     = newCounterVec(`stock_query`, `股票查询`, `code`)
	stockFetch     = newCounterVec(`stock_fetch`, `股票查询 实际网络请求`, `code`)
	stockFetchFail = newCounterVec(`stock_fetch_fail`, `股票查询 失败`, `code`)
)

func StockQuery(code string) {
	stockQuery.WithLabelValues(code).Inc()
}

func StockFetchFail(code string) {
	stockFetchFail.WithLabelValues(code).Inc()
}

func StockFetch(code string) {
	stockFetch.WithLabelValues(code).Inc()
}

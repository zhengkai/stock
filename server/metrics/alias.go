package metrics

import "github.com/prometheus/client_golang/prometheus"

var baseName = `stock_`

func newCounter(name, help string) prometheus.Counter {
	re := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: baseName + name,
			Help: help,
		},
	)
	pool = append(pool, re)
	return re
}
func newSummary(name, help string) prometheus.Summary {
	re := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       baseName + name,
		Help:       help,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	pool = append(pool, re)
	return re
}

func newGauge(name, help string) prometheus.Gauge {
	re := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: baseName + name,
		Help: help,
	})
	pool = append(pool, re)
	return re
}

func newCounterVec(name, help, field string) *prometheus.CounterVec {
	re := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: baseName + name,
			Help: help,
		},
		[]string{field},
	)
	pool = append(pool, re)
	return re
}

func newGaugeVec(name, help, field string) *prometheus.GaugeVec {
	re := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: baseName + name,
			Help: help,
		},
		[]string{field},
	)
	pool = append(pool, re)
	return re
}

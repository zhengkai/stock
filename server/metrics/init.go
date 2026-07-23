// Package metrics prometheus 监控
package metrics

import "github.com/prometheus/client_golang/prometheus"

var pool []prometheus.Collector

func init() {
	prometheus.MustRegister(pool...)
}

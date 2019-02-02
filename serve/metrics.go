package serve

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	accesses    prometheus.Counter
	errors      prometheus.Counter
	cacheMisses prometheus.Counter
	cacheSize   prometheus.Gauge
	calcTime    prometheus.Histogram
}

func register(name string) *metrics {
	errors := promauto.NewCounter(prometheus.CounterOpts{
		Name: name + "_errors",
	})
	accesses := promauto.NewCounter(prometheus.CounterOpts{
		Name: name + "_accesses",
	})
	cacheMisses := promauto.NewCounter(prometheus.CounterOpts{
		Name: name + "_cache_misses",
	})
	cacheSize := promauto.NewGauge(prometheus.GaugeOpts{
		Name: name + "_cache_size",
	})
	calcTime := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    name + "_calculation_time",
		Buckets: prometheus.ExponentialBuckets(0.5e-3, 2, 14),
	})

	return &metrics{
		accesses:    accesses,
		errors:      errors,
		cacheMisses: cacheMisses,
		cacheSize:   cacheSize,
		calcTime:    calcTime,
	}
}

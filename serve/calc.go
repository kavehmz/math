package serve

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var errTimeout = errors.New(`Operation timedout`)

func (m *MathHandler) calculate(f formula, b []byte) ([]byte, error) {
	f.metrics.accesses.Inc()
	timer := prometheus.NewTimer(f.metrics.calcTime)
	defer timer.ObserveDuration()

	var res []byte
	var err error
	if f.opt.Immutable() {
		res, err = m.calculateImmutable(f, b)
	} else {
		res, err = m.calculatemutable(f, b)

	}
	if err != nil {
		f.metrics.errors.Inc()
	}
	return res, err
}

func (m *MathHandler) calculateImmutable(f formula, b []byte) ([]byte, error) {
	defer f.metrics.cacheSize.Set(float64(f.cache.ItemCount()))

	var res []byte
	if res, found := f.cache.Get(string(b)); found {
		return res.([]byte), nil
	}
	f.metrics.cacheMisses.Inc()
	timer := prometheus.NewTimer(f.metrics.calcTime)
	defer timer.ObserveDuration()

	res, err := withTimeout(f.opt.Calc, b, m.timeout)
	if err != nil {
		return nil, err
	}
	f.cache.SetDefault(string(b), res)

	return res, nil
}

func (m *MathHandler) calculatemutable(f formula, b []byte) ([]byte, error) {
	return withTimeout(f.opt.Calc, b, m.timeout)
}

type result struct {
	b   []byte
	err error
}

func withTimeout(f func([]byte, chan bool) ([]byte, error), b []byte, timeout time.Duration) ([]byte, error) {
	end := make(chan bool)
	res := make(chan result, 1)
	go func() {
		r, err := f(b, end)
		res <- result{r, err}
	}()

	for {
		select {
		case <-time.After(timeout):
			return nil, errTimeout
		case r := <-res:
			return r.b, r.err
		}
	}
}

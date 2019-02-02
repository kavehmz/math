package serve

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpAccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_access",
	})
	httpErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_error",
	})
	httpTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_time",
		Buckets: prometheus.ExponentialBuckets(0.5e-3, 2, 14),
	})
)

func (m *MathHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpAccess.Inc()
	timer := prometheus.NewTimer(httpTime)
	defer timer.ObserveDuration()

	if r.Method != http.MethodPost {
		httpError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		httpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if f, ok := m.formulas[path.Clean(r.URL.Path)]; ok {
		f.metrics.accesses.Inc()

		result, err := m.calculate(f, b)
		if err == errTimeout {
			f.metrics.errors.Inc()
			httpError(w, err.Error(), http.StatusRequestTimeout)
			return
		}
		if err != nil {
			f.metrics.errors.Inc()
			httpError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(result)
		if err != nil {
			httpError(w, err.Error(), http.StatusBadGateway)
		}
		return
	}
	httpError(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func httpError(w http.ResponseWriter, error string, code int) {
	httpErrors.Inc()
	http.Error(w, error, code)
}

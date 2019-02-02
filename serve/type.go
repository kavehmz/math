package serve

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

type operation interface {
	Immutable() bool
	Calc(param []byte, end chan bool) ([]byte, error)
}

type formula struct {
	opt     operation
	metrics *metrics
	cache   *cache.Cache
}

// MathHandler will handle the math operations
type MathHandler struct {
	formulas map[string]formula
	timeout  time.Duration
}

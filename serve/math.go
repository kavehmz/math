package serve

import (
	"time"

	"github.com/kavehmz/math/formula/add"
	"github.com/kavehmz/math/formula/factorial"

	cache "github.com/patrickmn/go-cache"
)

// Math returns a math handler
func Math(timeout, ttl time.Duration) *MathHandler {
	m := initMath(timeout, ttl)
	m.formulas["/math/add"] = formulate("add", &add.Add{}, ttl)
	m.formulas["/math/factorial"] = formulate("factorial", &factorial.Factorial{}, ttl)
	return m
}

func initMath(timeout, ttl time.Duration) *MathHandler {
	return &MathHandler{
		timeout:  timeout,
		formulas: make(map[string]formula),
	}
}

func formulate(name string, opt operation, ttl time.Duration) formula {
	f := formula{cache: cache.New(ttl, time.Minute)}
	f.metrics = register(name)
	f.opt = opt
	return f
}

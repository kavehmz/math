package serve

import (
	"errors"
	"testing"
	"time"
)

type testFormula struct {
	im bool
}

func (a *testFormula) Immutable() bool {
	return a.im
}

func (a *testFormula) Calc(b []byte, end chan bool) ([]byte, error) {
	if string(b) == "err" {
		return nil, errors.New("test error")
	}
	time.Sleep(time.Millisecond)
	return []byte("2"), nil
}

func Test_math_calculateImmutable(t *testing.T) {
	m := MathHandler{
		timeout: time.Second,
	}
	m.formulas = make(map[string]formula)
	m.formulas["/math/test"] = formulate("test", &testFormula{im: true}, time.Minute)

	_, err := m.calculate(m.formulas["/math/test"], []byte(`err`))
	if err == nil {
		t.Error("expecting error")
	}

	m.formulas["/math/test"].cache.SetDefault("cachedTest", []byte("cachedValue"))
	val, err := m.calculate(m.formulas["/math/test"], []byte("cachedTest"))
	if err != nil || string(val) != "cachedValue" {
		t.Error("cache value wrong")
	}

	val, err = m.calculate(m.formulas["/math/test"], []byte(`{"x":1, "y":1}`))
	if err != nil || string(val) != "2" {
		t.Fatal("non cache value wrong")
	}

	cVal, found := m.formulas["/math/test"].cache.Get(`{"x":1, "y":1}`)
	if !found || string(cVal.([]byte)) != "2" {
		t.Error("cache failed")
	}
}

func Test_math_calculatemutable(t *testing.T) {
	m := MathHandler{timeout: time.Second}
	m.formulas = make(map[string]formula)
	m.formulas["/math/test"] = formulate("test2", &testFormula{}, time.Minute)

	_, err := m.calculate(m.formulas["/math/test"], []byte(`err`))
	if err == nil {
		t.Error("expecting error")
	}

	val, err := m.calculate(m.formulas["/math/test"], []byte(`{"x":1, "y":1}`))
	if err != nil || string(val) != "2" {
		t.Fatal("value wrong")
	}
}

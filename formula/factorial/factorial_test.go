package factorial

import (
	"testing"
)

func Test_extract(t *testing.T) {
	_, err := extract([]byte(`{"wrong":1}`))
	if err != ErrInvalidParameters {
		t.Error("Expected to error if x is not present")
	}

	_, err = extract([]byte(`{"wrong,}`))
	if err != ErrInvalidParameters {
		t.Error("Expected to error if x is not present")
	}

	_, err = extract([]byte(`{"x":-1}`))
	if err != ErrInvalidParameters {
		t.Error("Expected to error if x is negative")
	}

	p, err := extract([]byte(`{"x":1}`))
	if err != nil {
		t.Fatal("Expected no error")
	}

	if p.X == nil || *p.X != 1 {
		t.Fatal("Wrong x")
	}

}

func Test_factorial(t *testing.T) {
	end := make(chan bool, 1)
	v, e := factorial(0, end)
	if string(v) != "1" || e != nil {
		t.Error("Wrong number")
	}

	v, e = factorial(1, end)
	if string(v) != "1" || e != nil {
		t.Error("Wrong number")
	}

	v, e = factorial(5, end)
	if string(v) != "120" || e != nil {
		t.Error("Wrong number")
	}

	end <- true
	_, e = factorial(3, end)
	if e != ErrTimeout {
		t.Error("Expected error")
	}
}

func TestFactorial_Immutable(t *testing.T) {
	m := Factorial{}
	if !m.Immutable() {
		t.Error("Expected to be true")
	}
}

func TestFactorial_Calc(t *testing.T) {
	m := Factorial{}
	end := make(chan bool, 1)
	if _, err := m.Calc([]byte("bad input"), end); err != ErrInvalidParameters {
		t.Error("Was expecting ErrInvalidParameters", err)
	}
	end <- true
	if _, err := m.Calc([]byte(`{"x":5}`), end); err != ErrTimeout {
		t.Error("Was expecting ErrTimeout", err)
	}
}

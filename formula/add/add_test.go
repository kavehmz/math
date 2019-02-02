package add

import (
	"testing"
)

func Test_extract(t *testing.T) {
	_, err := extract([]byte(`{"x":1,"wrong":1}`))
	if err != ErrInvalidParameters {
		t.Error("Expected to error if y is not present")
	}

	_, err = extract([]byte(`{"wrong":1,"y":1}`))
	if err != ErrInvalidParameters {
		t.Error("Expected to error if x is not present")
	}
	_, err = extract([]byte(`{"wrong,"y":1}`))
	if err != ErrInvalidParameters {
		t.Error("Expected to error if x is not present")
	}

	p, err := extract([]byte(`{"x":1,"y":1}`))
	if err != nil {
		t.Fatal("Expected no error")
	}

	if p.X == nil || *p.X != 1 {
		t.Fatal("Wrong x")
	}

	if p.Y == nil || *p.Y != 1 {
		t.Fatal("Wrong y")
	}
}

func Test_add(t *testing.T) {
	if string(add(1, 1)) != "2" {
		t.Error("Wrong sum", string(add(1, 1)))
	}
}
func TestAdd_Immutable(t *testing.T) {
	m := Add{}
	if m.Immutable() {
		t.Error("Expected to be false")
	}
}

func TestAdd_Calc(t *testing.T) {
	m := Add{}
	if _, err := m.Calc([]byte("bad input"), make(chan bool)); err != ErrInvalidParameters {
		t.Error("Was expecting ErrInvalidParameters", err)
	}

	if _, err := m.Calc([]byte(`{"x":1,"y":1}`), make(chan bool)); err != nil {
		t.Error("Was expecting no errors", err)
	}
}

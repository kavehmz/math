package serve

import (
	"testing"
	"time"
)

func TestMath(t *testing.T) {
	m := Math(time.Second, time.Minute)
	if m.formulas == nil || len(m.formulas) == 0 {
		t.Error("route init issue")
	}
}

func Test_initMath(t *testing.T) {
	m := initMath(time.Second, time.Minute)
	if m.formulas == nil || m.timeout != time.Second {
		t.Error("Init issue")
	}
}

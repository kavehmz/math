package add

import (
	"encoding/json"
	"errors"
	"math/big"
)

// ErrInvalidParameters indicates input error
var ErrInvalidParameters = errors.New(`Add needs: {"x": int, "y":int}`)

// Add operation
type Add struct {
}

type params struct {
	X *int64 `json:"x"`
	Y *int64 `json:"y"`
}

// Immutable if it is but for add we say no!
func (a *Add) Immutable() bool {
	return false
}

// Calc the value and ignore the end signal
func (a *Add) Calc(b []byte, end chan bool) ([]byte, error) {
	p, err := extract(b)
	if err != nil {
		return nil, err
	}

	return add(*p.X, *p.Y), nil
}

func extract(b []byte) (*params, error) {
	p := params{}
	err := json.Unmarshal(b, &p)
	if err != nil || p.X == nil || p.Y == nil {
		return nil, ErrInvalidParameters
	}

	return &p, nil
}

func add(x, y int64) []byte {
	var xb, yb big.Int
	xb.SetInt64(x)
	yb.SetInt64(y)

	return []byte(xb.Add(&xb, &yb).String())
}

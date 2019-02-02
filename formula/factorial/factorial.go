package factorial

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

// ErrInvalidParameters indicates input error
var ErrInvalidParameters = errors.New(`Factorial needs: {"x": int} and x>=0`)

// ErrTimeout indicated timeout error based on caller request
var ErrTimeout = errors.New(`Operation timedout`)

// Factorial operation
type Factorial struct {
}

type params struct {
	X *int64 `json:"x"`
}

// Immutable if it iss
func (a *Factorial) Immutable() bool {
	return true
}

// Calc the value
func (a *Factorial) Calc(b []byte, end chan bool) ([]byte, error) {
	p, err := extract(b)
	if err != nil {
		return nil, err
	}

	return factorial(*p.X, end)
}

func extract(b []byte) (*params, error) {
	p := params{}
	err := json.Unmarshal(b, &p)
	if err != nil || p.X == nil || *p.X < 0 {
		return nil, ErrInvalidParameters
	}

	return &p, nil
}

func factorial(n int64, end chan bool) ([]byte, error) {
	xb := big.Int{}
	xb.SetInt64(1)

	for i := int64(2); i <= n; i++ {
		select {
		case x := <-end:
			fmt.Println(x)
			return nil, ErrTimeout
		default:
		}

		yb := big.Int{}
		yb.SetInt64(i)
		xb.Mul(&xb, &yb)
	}

	return []byte(xb.String()), nil
}

package serve

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMathHandler_ServeHTTP_app_err(t *testing.T) {
	m := MathHandler{timeout: time.Second}
	m.formulas = make(map[string]formula)
	m.formulas["/math/test"] = formulate("test_serve", &testFormula{}, time.Minute)

	var jsonStr = []byte(`err`)
	req := httptest.NewRequest("POST", "http://127.0.0.1/math/test", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected error", resp.StatusCode)
	}
}
func TestMathHandler_ServeHTTP_timeout(t *testing.T) {
	m := MathHandler{timeout: time.Nanosecond}
	m.formulas = make(map[string]formula)
	m.formulas["/math/test"] = formulate("test_serve2", &testFormula{}, time.Minute)

	var jsonStr = []byte(`param`)
	req := httptest.NewRequest("POST", "http://127.0.0.1/math/test", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusRequestTimeout {
		t.Error("Expected timeout", resp.StatusCode)
	}
}

type testReadWrite struct {
	// We will override the Rad/Write only. For other calls this will be used
	w *httptest.ResponseRecorder
}

func (f testReadWrite) Read(p []byte) (n int, err error) { return 0, errors.New("read error") }
func (f testReadWrite) Header() http.Header              { return f.w.Header() }
func (f testReadWrite) Write([]byte) (int, error)        { return 0, errors.New("write error") }
func (f testReadWrite) WriteHeader(statusCode int)       { f.w.WriteHeader(statusCode) }

func TestMathHandler_ServeHTTP_http_err(t *testing.T) {
	m := MathHandler{timeout: time.Second}
	m.formulas = make(map[string]formula)
	m.formulas["/math/test"] = formulate("test_serve3", &testFormula{}, time.Minute)

	var jsonStr = []byte(`param`)
	req := httptest.NewRequest("POST", "http://127.0.0.1/math/notexists", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected error", resp.StatusCode)
	}

	req = httptest.NewRequest("GET", "http://127.0.0.1/math/test", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	m.ServeHTTP(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Expected error", resp.StatusCode)
	}

	req = httptest.NewRequest("POST", "http://127.0.0.1/math/test", testReadWrite{})
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	m.ServeHTTP(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Expected error", resp.StatusCode)
	}

	req = httptest.NewRequest("POST", "http://127.0.0.1/math/test", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	m.ServeHTTP(testReadWrite{w: w}, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusBadGateway {
		t.Error("Expected error", resp.StatusCode)
	}
}
